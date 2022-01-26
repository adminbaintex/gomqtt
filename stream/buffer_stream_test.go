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
	"io/ioutil"
	"testing"

	"github.com/adminbaintex/gomqtt/packet"
	"github.com/stretchr/testify/require"
)

func TestBufferStream(t *testing.T) {
	r := new(bytes.Buffer)
	w := new(bytes.Buffer)

	m := packet.NewConnackPacket()
	r.Write(encodePacket(m))

	s := NewBufferStream(r, w)
	ret1 := s.Send(m)
	m2 := <-s.Incoming()
	s.Close()
	ret2 := s.Send(m)

	require.Equal(t, m.Len(), w.Len())
	require.Equal(t, m.Len(), m2.Len())
	require.True(t, s.Closed())
	require.True(t, ret1)
	require.False(t, ret2)
}

func TestBufferStreamEOF(t *testing.T) {
	r := new(bytes.Buffer)
	w := new(bytes.Buffer)
	s := NewBufferStream(r, w)

	_, in := <-s.Incoming()
	require.False(t, in)
	require.True(t, s.Closed())
}

func TestBufferStreamError(t *testing.T) {
	r := new(infiniteReader)
	w := new(bytes.Buffer)
	s := NewBufferStream(r, w)
	m := packet.NewConnackPacket()
	m.ReturnCode = 11 // invalid return code

	s.Outgoing() <- m
	_, ok := <-s.Incoming()

	require.False(t, ok)
	require.Equal(t, 0, w.Len())
	require.Error(t, s.Error())
	require.True(t, s.Closed())
}

func BenchmarkBufferStream(b *testing.B) {
	m := packet.NewConnectPacket()
	r := bufio.NewReader(newMessageReader(m))
	w := bufio.NewWriter(ioutil.Discard)

	s := NewBufferStream(r, w)

	for i := 0; i < b.N; i++ {
		s.Send(<-s.Incoming())
	}

	s.Close()
}
