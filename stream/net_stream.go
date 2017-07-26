// Copyright (c) 2014 The gomqtt Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stream

import (
	"bufio"
	"io"
	"net"

	"git.baintex.com/sentio/gomqtt/packet"
)

// The NetStream wraps an io.reader and an io.writer.
type NetStream struct {
	AbstractStream

	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

// NewNetStream returns a new NetStream.
func NewNetStream(conn net.Conn) *NetStream {
	s := &NetStream{}
	s.conn = conn
	s.reader = bufio.NewReader(conn)
	s.writer = bufio.NewWriter(conn)

	s.Initialize(s.decoder, s.encoder, s.closer)
	return s
}

func (ns *NetStream) decoder() (packet.Packet, error) {
	m, _, e := DecodeFromReader(ns.reader)
	if e == io.EOF {
		e = ErrExpectedClose
	}
	return m, e
}

func (ns *NetStream) encoder(pkt packet.Packet) error {
	_, e := EncodeToWriter(ns.writer, pkt)
	return e
}

func (ns *NetStream) closer() {
	ns.conn.Close()
}
