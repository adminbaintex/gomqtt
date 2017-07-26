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

	"github.com/gomqtt/packet"
)

// The BufferStream wraps an io.reader and an io.writer.
type BufferStream struct {
	AbstractStream

	reader *bufio.Reader
	writer *bufio.Writer
}

// NewBufferStream returns a new BufferStream.
func NewBufferStream(r io.Reader, w io.Writer) *BufferStream {
	s := &BufferStream{}
	s.reader = bufio.NewReader(r)
	s.writer = bufio.NewWriter(w)

	s.Initialize(s.decoder, s.encoder, nil)
	return s
}

func (bs *BufferStream) decoder() (packet.Packet, error) {
	m, _, e := DecodeFromReader(bs.reader)
	if e == io.EOF {
		e = ErrExpectedClose
	}
	return m, e
}

func (bs *BufferStream) encoder(pkt packet.Packet) error {
	_, e := EncodeToWriter(bs.writer, pkt)
	return e
}
