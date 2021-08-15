/*
Copyright © 2021 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package etcd

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/mitchellh/hashstructure"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"

	"github.com/iTrellis/common/backoff"
	"github.com/iTrellis/common/codec"
	"github.com/iTrellis/common/errors"
	"github.com/iTrellis/common/formats"
	"github.com/iTrellis/common/logger"
	iTLS "github.com/iTrellis/common/tls"
)

// Config for a new etcd.Client.
type Config struct {
	Endpoints   []string          `yaml:"endpoints"`
	DialTimeout formats.Duration  `yaml:"dial_timeout"`
	MaxRetries  int               `yaml:"max_retries"`
	EnableTLS   bool              `yaml:"tls_enabled"`
	UserName    string            `yaml:"username"`
	Password    formats.Hide      `yaml:"password"`
	TLS         iTLS.ClientConfig `yaml:",inline"`
	TTL         formats.Duration  `yaml:"ttl"` // use with lease

	ZapLogger *zap.Logger `yaml:"-"`
}

// Client implements ring.KVClient for etcd.
type Client struct {
	cfg   Config
	codec codec.Codec
	cli   *clientv3.Client

	logger logger.Logger

	mu     sync.RWMutex
	leases map[string]clientv3.LeaseID
	hashs  map[string]uint64
}

// GetTLS sets the TLS config field with certs
func (p *Config) GetTLS() (*tls.Config, error) {
	if !p.EnableTLS {
		return nil, nil
	}
	tlsInfo := &transport.TLSInfo{
		CertFile:           p.TLS.CertPath,
		KeyFile:            p.TLS.KeyPath,
		TrustedCAFile:      p.TLS.CAPath,
		ServerName:         p.TLS.ServerName,
		InsecureSkipVerify: p.TLS.InsecureSkipVerify,
	}
	return tlsInfo.ClientConfig()
}

// New makes a new Client.
func New(cfg Config, codec codec.Codec) (*Client, error) {
	tlsConfig, err := cfg.GetTLS()
	if err != nil {
		return nil, errors.NewErrors(errors.New("unable to initialise TLS configuration for etcd"), err)
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: time.Duration(cfg.DialTimeout),
		// Configure the keepalive to make sure that the client reconnects
		// to the etcd service endpoint(s) in case the current connection is
		// dead (ie. the node where etcd is running is dead or a network
		// partition occurs).
		//
		// The settings:
		// - DialKeepAliveTime: time before the client pings the server to
		//   see if transport is alive (10s hardcoded)
		// - DialKeepAliveTimeout: time the client waits for a response for
		//   the keep-alive probe (set to 2x dial timeout, in order to avoid
		//   exposing another config option which is likely to be a factor of
		//   the dial timeout anyway)
		// - PermitWithoutStream: whether the client should send keepalive pings
		//   to server without any active streams (enabled)
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 2 * time.Duration(cfg.DialTimeout),
		PermitWithoutStream:  true,
		TLS:                  tlsConfig,
		Username:             cfg.UserName,
		Password:             string(cfg.Password),
		Logger:               cfg.ZapLogger,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		cfg:    cfg,
		codec:  codec,
		cli:    cli,
		logger: logger.NewWithZapLogger(cfg.ZapLogger),
		leases: make(map[string]clientv3.LeaseID),
		hashs:  make(map[string]uint64),
	}, nil
}

// CAS implements kv.Client.
// TODO Test with lease
func (p *Client) CAS(ctx context.Context, key string,
	f func(in interface{}) (out interface{}, retry bool, err error)) error {
	var revision int64
	var lastErr error

	for i := 0; i < p.cfg.MaxRetries; i++ {
		p.mu.Lock()
		leaseID, ok := p.leases[key]
		p.mu.Unlock()

		var intermediate interface{}
		if !ok {
			// missing lease, check if the key exists.
			resp, err := p.cli.Get(ctx, key, clientv3.WithSerializable())
			if err != nil {
				p.logger.Error("error getting key", "key", key, "err", err)
				lastErr = err
				continue
			}

			for _, kv := range resp.Kvs {
				if kv.Lease > 0 {
					leaseID = clientv3.LeaseID(kv.Lease)

					intermediate, err = p.codec.Unmarshal(kv.Value)
					if err != nil {
						lastErr = err
						continue
					}

					intermediate, _, err = f(intermediate)
					if err != nil {
						lastErr = err
						continue
					}

					h, err := hashstructure.Hash(intermediate, nil)
					if err != nil {
						lastErr = err
						continue
					}

					p.mu.Lock()
					p.leases[key] = leaseID
					p.hashs[key] = h
					p.mu.Unlock()

					revision = kv.Version

					break
				}
			}
		}

		var leaseNotFound bool

		if leaseID > 0 {
			p.logger.Debugf("Renewing existing lease for %s %d", key, leaseID)

			if _, err := p.cli.KeepAliveOnce(context.TODO(), leaseID); err != nil {
				if err != rpctypes.ErrLeaseNotFound {
					return err
				}

				p.logger.Debugf("Lease not found for %s %d", key, leaseID)

				// lease not found do register
				leaseNotFound = true
			}
		}

		// get existing hash for the service node
		p.mu.RLock()
		v, ok := p.hashs[key]
		p.mu.RUnlock()

		var err error
		if intermediate == nil {
			intermediate, _, err = f(nil)
			if err != nil {
				lastErr = err
				continue
			}
		}

		// create hash of service; uint64
		h, err := hashstructure.Hash(intermediate, nil)
		if err != nil {
			return err
		}

		// the service is unchanged, skip registering
		if ok && v == h && !leaseNotFound {
			p.logger.Infof(" %s unchanged skipping registration", key)
			return nil
		}

		buf, err := p.codec.Marshal(intermediate)
		if err != nil {
			lastErr = err
			continue
		}

		var lgr *clientv3.LeaseGrantResponse
		if p.cfg.TTL.Seconds() > 0 {
			// get a lease used to expire keys since we have a ttl
			lgr, err = p.cli.Grant(ctx, int64(p.cfg.TTL.Seconds()))
			if err != nil {
				return err
			}
		}

		// create an entry for the node
		var putOpts []clientv3.OpOption
		if lgr != nil {
			putOpts = append(putOpts, clientv3.WithLease(lgr.ID))

			p.logger.Debugf("Registering %s with lease %v and leaseID %v and ttl %v", key, lgr, lgr.ID, p.cfg.TTL)
		}

		result, err := p.cli.Txn(ctx).
			If(clientv3.Compare(clientv3.Version(key), "=", revision)).
			Then(clientv3.OpPut(key, string(buf), putOpts...)).
			Commit()
		if err != nil {
			// level.Error(util_log.Logger).Log("msg", "error CASing", "key", key, "err", err)
			lastErr = err
			continue
		}
		// result is not Succeeded if the the comparison was false, meaning if the modify indexes did not match.
		if !result.Succeeded {
			// level.Debug(util_log.Logger).Log("msg", "failed to CAS, revision and version did not match in etcd", "key", key, "revision", revision)
			continue
		}

		p.mu.Lock()
		p.leases[key] = leaseID
		p.hashs[key] = h
		p.mu.Unlock()

		return nil
	}

	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("failed to CAS %s", key)
}

// Delete implements kv.Client.
func (c *Client) Delete(ctx context.Context, key string) error {
	_, err := c.cli.Delete(ctx, key)
	return err
}

// Get implements kv.Client.
func (c *Client) Get(ctx context.Context, key string) (interface{}, error) {
	resp, err := c.cli.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	} else if len(resp.Kvs) != 1 {
		return nil, fmt.Errorf("got %d kvs, expected 1 or 0", len(resp.Kvs))
	}
	return c.codec.Unmarshal(resp.Kvs[0].Value)
}

// List implements kv.Client.
func (c *Client) List(ctx context.Context, prefix string) ([]string, error) {
	resp, err := c.cli.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		keys = append(keys, string(kv.Key))
	}
	return keys, nil
}

// WatchKey implements kv.Client.
func (c *Client) WatchKey(ctx context.Context, key string, f func(interface{}) bool) {
	backoff := backoff.New(ctx, backoff.Config{
		MinBackoff: 1 * time.Second,
		MaxBackoff: 1 * time.Minute,
	})

	// Ensure the context used by the Watch is always cancelled.
	watchCtx, cancel := context.WithCancel(ctx)
	defer cancel()

outer:
	for backoff.Ongoing() {
		for resp := range c.cli.Watch(watchCtx, key) {
			if err := resp.Err(); err != nil {
				// level.Error(util_log.Logger).Log("msg", "watch error", "key", key, "err", err)
				continue outer
			}

			backoff.Reset()

			for _, event := range resp.Events {
				out, err := c.codec.Unmarshal(event.Kv.Value)
				if err != nil {
					// level.Error(util_log.Logger).Log("msg", "error decoding key", "key", key, "err", err)
					continue
				}

				if !f(out) {
					return
				}
			}
		}
	}
}

// WatchPrefix implements kv.Client.
func (c *Client) WatchPrefix(ctx context.Context, key string, f func(string, interface{}) bool) {
	backoff := backoff.New(ctx, backoff.Config{
		MinBackoff: 1 * time.Second,
		MaxBackoff: 1 * time.Minute,
	})

	// Ensure the context used by the Watch is always cancelled.
	watchCtx, cancel := context.WithCancel(ctx)
	defer cancel()

outer:
	for backoff.Ongoing() {
		for resp := range c.cli.Watch(watchCtx, key, clientv3.WithPrefix()) {
			if err := resp.Err(); err != nil {
				// level.Error(util_log.Logger).Log("msg", "watch error", "key", key, "err", err)
				continue outer
			}

			backoff.Reset()

			for _, event := range resp.Events {
				if event.Kv.Version == 0 && event.Kv.Value == nil {
					// Delete notification. Since not all KV store clients (and Cortex codecs) support this, we ignore it.
					continue
				}

				out, err := c.codec.Unmarshal(event.Kv.Value)
				if err != nil {
					// level.Error(util_log.Logger).Log("msg", "error decoding key", "key", key, "err", err)
					continue
				}

				if !f(string(event.Kv.Key), out) {
					return
				}
			}
		}
	}
}
