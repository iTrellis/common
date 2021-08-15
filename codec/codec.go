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

package codec

import (
	"github.com/golang/snappy"
	"github.com/iTrellis/common/json"
	"google.golang.org/protobuf/proto"
)

var (
	_ Codec = (*Proto)(nil)
	_ Codec = (*String)(nil)
	_ Codec = (*JSON)(nil)
)

// NewCodec Takes in a connection/buffer and returns a new Codec
type NewCodec func() Codec

// Codec is a simple encoding interface used for the broker/transport
// where headers are not supported by the underlying implementation.
type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(bs []byte) (interface{}, error)
	String() string
}

// Proto is a Codec for proto/snappy
type Proto struct {
	id      string
	factory func() proto.Message
}

func NewProtoCodec(id string, factory func() proto.Message) Proto {
	return Proto{id: id, factory: factory}
}

func (p Proto) String() string {
	return p.id
}

// Marshal implements Codec
func (p Proto) Marshal(msg interface{}) ([]byte, error) {
	bytes, err := proto.Marshal(msg.(proto.Message))
	if err != nil {
		return nil, err
	}
	return snappy.Encode(nil, bytes), nil
}

// Unmarshal implements Codec
func (p Proto) Unmarshal(bytes []byte) (interface{}, error) {
	out := p.factory()
	bytes, err := snappy.Decode(nil, bytes)
	if err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(bytes, out); err != nil {
		return nil, err
	}
	return out, nil
}

// String is a code for strings.
type String struct{}

func (String) String() string {
	return "string"
}

// Unmarshal implements Codec.
func (String) Unmarshal(bytes []byte) (interface{}, error) {
	return string(bytes), nil
}

// Marshal implements Codec.
func (String) Marshal(msg interface{}) ([]byte, error) {
	return []byte(msg.(string)), nil
}

type JSON struct {
	id      string
	factory func() interface{}
}

func NewJSONCodec(id string, factory func() interface{}) *JSON {
	return &JSON{id: id, factory: factory}
}

func (*JSON) String() string {
	return "json"
}

// Unmarshal implements Codec.
func (p *JSON) Unmarshal(msg []byte) (interface{}, error) {
	out := p.factory()
	if err := json.Unmarshal(msg, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Marshal implements Codec.
func (p *JSON) Marshal(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}
