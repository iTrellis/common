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

package discovery

import (
	"context"

	"github.com/iTrellis/common/codec"
	"github.com/iTrellis/common/discovery/etcd"
	"github.com/iTrellis/common/errors"
)

// // Labels method returns Prometheus labels relevant to itself.
// func (r *role) Labels() prometheus.Labels {
// 	return prometheus.Labels{"role": string(*r)}
// }

// StoreConfig is a configuration used for building single store client, either
// Consul, Etcd, Memberlist or MultiClient. It was extracted from Config to keep
// single-client config separate from final client-config (with all the wrappers)
type Config struct {
	// // Consul consul.Config `yaml:"consul"`
	Etcd etcd.Config `yaml:"etcd"`
	// // Multi  MultiConfig   `yaml:"multi"`

	// Function that returns memberlist.KV store to use. By using a function, we can delay
	// initialization of memberlist.KV until it is actually required.
	// MemberlistKV func() (*memberlist.KV, error) `yaml:"-"`
}

// Client is a high-level client for key-value stores (such as Etcd and
// Consul) that exposes operations such as CAS and Watch which take callbacks.
// It also deals with serialisation by using a Codec and having a instance of
// the the desired type passed in to methods ala json.Unmarshal.
type Client interface {
	// List returns a list of keys under the given prefix. Returned keys will
	// include the prefix.
	List(ctx context.Context, prefix string) ([]string, error)

	// Get a specific key.  Will use a codec to deserialise key to appropriate type.
	// If the key does not exist, Get will return nil and no error.
	Get(ctx context.Context, key string) (interface{}, error)

	// Delete a specific key. Deletions are best-effort and no error will
	// be returned if the key does not exist.
	Delete(ctx context.Context, key string) error

	// CAS stands for Compare-And-Swap.  Will call provided callback f with the
	// current value of the key and allow callback to return a different value.
	// Will then attempt to atomically swap the current value for the new value.
	// If that doesn't succeed will try again - callback will be called again
	// with new value etc.  Guarantees that only a single concurrent CAS
	// succeeds.  Callback can return nil to indicate it is happy with existing
	// value.
	CAS(ctx context.Context, key string, f func(in interface{}) (out interface{}, retry bool, err error)) error

	// WatchKey calls f whenever the value stored under key changes.
	WatchKey(ctx context.Context, key string, f func(interface{}) bool)

	// WatchPrefix calls f whenever any value stored under prefix changes.
	WatchPrefix(ctx context.Context, prefix string, f func(string, interface{}) bool)
}

// , reg prometheus.Registerer
func createClient(backend string, prefix string, cfg Config, codec codec.Codec) (Client, error) {
	var client Client
	var err error

	switch backend {
	// case "consul":
	// 	client, err = consul.NewClient(cfg.Consul, codec)

	case "etcd":
		client, err = etcd.New(cfg.Etcd, codec)

	// This case is for testing. The mock KV client does not do anything internally.
	case "mock":
		client, err = buildMockClient()

	default:
		return nil, errors.Newf("invalid discovery type: %s", backend)
	}

	if err != nil {
		return nil, err
	}

	if prefix != "" {
		client = PrefixClient(client, prefix)
	}

	// If no Registerer is provided return the raw client.
	// if reg == nil {
	// return client, nil
	// }
	// return newMetricsClient(backend, client, prometheus.WrapRegistererWith(role.Labels(), reg)), nil

	return client, nil

}
