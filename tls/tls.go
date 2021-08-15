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

package tls

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/iTrellis/common/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	errKeyMissing  = errors.New("certificate given but no key configured")
	errCertMissing = errors.New("key given but no certificate configured")
)

// ClientConfig is the config for client TLS.
type ClientConfig struct {
	CertPath           string `yaml:"tls_cert_path"`
	KeyPath            string `yaml:"tls_key_path"`
	CAPath             string `yaml:"tls_ca_path"`
	ServerName         string `yaml:"tls_server_name"`
	InsecureSkipVerify bool   `yaml:"tls_insecure_skip_verify"`
}

// GetTLSConfig initialises tls.Config from config options
func (p *ClientConfig) GetTLSConfig() (*tls.Config, error) {
	config := &tls.Config{
		InsecureSkipVerify: p.InsecureSkipVerify,
		ServerName:         p.ServerName,
	}

	// read ca certificates
	if p.CAPath != "" {
		var caCertPool *x509.CertPool
		caCert, err := ioutil.ReadFile(p.CAPath)
		if err != nil {
			// return nil, errors.Wrapf(err, "error loading ca cert: %s", p.CAPath)
			return nil, errors.NewErrors(errors.Newf("error loading ca cert: %s", p.CAPath), err)
		}
		caCertPool = x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		config.RootCAs = caCertPool
	}

	// read client certificate
	if p.CertPath != "" || p.KeyPath != "" {
		if p.CertPath == "" {
			return nil, errCertMissing
		}
		if p.KeyPath == "" {
			return nil, errKeyMissing
		}
		clientCert, err := tls.LoadX509KeyPair(p.CertPath, p.KeyPath)
		if err != nil {
			return nil, errors.NewErrors(
				errors.Newf("failed to load TLS certificate %s,%s", p.CertPath, p.KeyPath), err)
		}
		config.Certificates = []tls.Certificate{clientCert}
	}

	return config, nil
}

// GetGRPCDialOptions creates GRPC DialOptions for TLS
func (cfg *ClientConfig) GetGRPCDialOptions(enabled bool) ([]grpc.DialOption, error) {
	if !enabled {
		return []grpc.DialOption{grpc.WithInsecure()}, nil
	}

	tlsConfig, err := cfg.GetTLSConfig()
	if err != nil {
		// return nil, errors.Wrap(err, "error creating grpc dial options")
		// errors.NewErrors(errors.Newf("error creating grpc dial options"), err)
		return nil, errors.NewErrors(errors.New("error creating grpc dial options"), err)
	}

	return []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))}, nil
}
