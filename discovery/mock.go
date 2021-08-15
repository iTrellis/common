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

import "context"

// The mockClient does not anything.
// This is used for testing only.
type mockClient struct{}

func buildMockClient() (Client, error) {
	return mockClient{}, nil
}

func (m mockClient) List(ctx context.Context, prefix string) ([]string, error) {
	return []string{}, nil
}

func (m mockClient) Get(ctx context.Context, key string) (interface{}, error) {
	return "", nil
}

func (m mockClient) Delete(ctx context.Context, key string) error {
	return nil
}

func (m mockClient) CAS(ctx context.Context, key string,
	f func(in interface{}) (out interface{}, retry bool, err error)) error {
	return nil
}

func (m mockClient) WatchKey(ctx context.Context, key string, f func(interface{}) bool) {
}

func (m mockClient) WatchPrefix(ctx context.Context, prefix string, f func(string, interface{}) bool) {
}
