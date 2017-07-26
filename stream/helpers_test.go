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
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/gomqtt/packet"
	"github.com/stretchr/testify/require"
	"io/ioutil"
)

func TestEncodeToWriter(t *testing.T) {
	b := new(bytes.Buffer)
	w := bufio.NewWriter(b)
	m := packet.NewConnackPacket()

	n, err := EncodeToWriter(w, m)

	require.Equal(t, m.Len(), b.Len())
	require.Equal(t, m.Len(), n)
	require.NoError(t, err)
}

func TestEncodeToWriterError1(t *testing.T) {
	b := new(bytes.Buffer)
	w := bufio.NewWriter(b)
	m := packet.NewConnackPacket()
	m.ReturnCode = 11 // invalid return code

	n, err := EncodeToWriter(w, m)

	require.Equal(t, 0, b.Len())
	require.Equal(t, 0, n)
	require.Error(t, err)
}

func TestEncodeToWriterError2(t *testing.T) {
	b := new(errorWriter) // returns error
	w := bufio.NewWriter(b)
	m := packet.NewConnackPacket()

	n, err := EncodeToWriter(w, m)

	require.Equal(t, 4, n)
	require.Error(t, err)
}

func TestEncodeToWriterError3(t *testing.T) {
	b := new(errorWriter)          // returns error
	w := bufio.NewWriterSize(b, 1) // force direct write
	m := packet.NewConnackPacket()

	n, err := EncodeToWriter(w, m)

	require.Equal(t, 0, n)
	require.Error(t, err)
}

func TestDecodeFromReader(t *testing.T) {
	b := new(bytes.Buffer)
	w := bufio.NewWriter(b)
	r := bufio.NewReader(b)
	m := packet.NewConnackPacket()

	EncodeToWriter(w, m)

	m2, n, err := DecodeFromReader(r)

	require.Equal(t, m.String(), m2.String())
	require.Equal(t, m.Len(), n)
	require.NoError(t, err)
}

func TestDecodeFromReaderError1(t *testing.T) {
	b := new(bytes.Buffer)
	w := bufio.NewWriter(b)
	r := bufio.NewReader(b)
	m := packet.NewConnackPacket()

	EncodeToWriter(w, m)

	b.Truncate(2) // to short

	_, n, err := DecodeFromReader(r)

	require.Equal(t, 2, n)
	require.Error(t, err)
}

func TestDecodeFromReaderError2(t *testing.T) {
	buf := make([]byte, 6)
	buf[1] = 0x00
	binary.PutVarint(buf[1:], 1234567890) // not detectable length

	b := bytes.NewBuffer(buf)
	r := bufio.NewReader(b)

	_, n, err := DecodeFromReader(r)

	require.Equal(t, 0, n)
	require.Error(t, err)
}

func TestDecodeFromReaderError3(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x00, 0x00}) // invalid packet type
	r := bufio.NewReader(b)

	_, n, err := DecodeFromReader(r)

	require.Equal(t, 2, n)
	require.Error(t, err)
}

func TestDecodeFromReaderError4(t *testing.T) {
	b := bytes.NewBuffer([]byte{0x20, 0x02, 0x00, 0x06}) // invalid packet
	r := bufio.NewReader(b)

	_, n, err := DecodeFromReader(r)

	require.Equal(t, 4, n)
	require.Error(t, err)
}

func BenchmarkEncodeToWriter(b *testing.B) {
	w := bufio.NewWriter(ioutil.Discard)
	m := packet.NewConnectPacket()

	for i := 0; i < b.N; i++ {
		EncodeToWriter(w, m)
	}
}

func BenchmarkDecodeFromReader(b *testing.B) {
	m := packet.NewConnectPacket()
	r := bufio.NewReader(newMessageReader(m))

	for i := 0; i < b.N; i++ {
		DecodeFromReader(r)
	}
}
