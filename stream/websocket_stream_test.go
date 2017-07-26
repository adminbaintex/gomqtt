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
	"testing"

	"git.baintex.com/sentio/gomqtt/packet"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestWebSocketStream(t *testing.T) {
	done := make(chan int, 1)

	s := startTestWSServer(func(conn *websocket.Conn) {
		s1 := NewWebSocketStream(conn)

		m := <-s1.Incoming()
		require.Equal(t, m.Type(), packet.CONNECT)

		s1.Send(packet.NewConnackPacket())

		_, ok := <-s1.Incoming()
		require.False(t, ok)
		require.True(t, s1.Closed())

		done <- 1
	})

	s2 := NewWebSocketStream(newTestWSConnection())
	s2.Send(packet.NewConnectPacket())

	m := <-s2.Incoming()
	require.Equal(t, m.Type(), packet.CONNACK)

	s2.Close()

	<-done
	s.Close()

	require.True(t, s2.Closed())
}

func TestWebSocketStreamClose(t *testing.T) {
	done := make(chan int, 1)

	s := startTestWSServer(func(conn *websocket.Conn) {
		s1 := NewWebSocketStream(conn)
		s1.Close()
		require.True(t, s1.Closed())

		done <- 1
	})

	s2 := NewWebSocketStream(newTestWSConnection())

	_, ok := <-s2.Incoming()
	require.False(t, ok)

	<-done
	s.Close()

	require.True(t, s2.Closed())
}

func TestWebSocketStreamEncodeError(t *testing.T) {
	done := make(chan int, 1)

	s := startTestWSServer(func(conn *websocket.Conn) {
		s1 := NewWebSocketStream(conn)

		pkt := packet.NewConnackPacket()
		pkt.ReturnCode = 11 // invalid return code
		s1.Send(pkt)

		_, ok := <-s1.Incoming()
		require.False(t, ok)
		require.Error(t, s1.Error())
		require.True(t, s1.Closed())

		done <- 1
	})

	s2 := NewWebSocketStream(newTestWSConnection())

	_, ok := <-s2.Incoming()
	require.False(t, ok)

	<-done
	s.Close()

	require.True(t, s2.Closed())
}

func TestWebSocketStreamDecodeError1(t *testing.T) {
	done := make(chan int, 1)

	s := startTestWSServer(func(conn *websocket.Conn) {
		s1 := NewWebSocketStream(conn)

		// invalid packet type
		s1.conn.WriteMessage(websocket.BinaryMessage, []byte{0x00, 0x00})

		_, ok := <-s1.Incoming()
		require.False(t, ok)
		require.NoError(t, s1.Error())
		require.True(t, s1.Closed())

		done <- 1
	})

	s2 := NewWebSocketStream(newTestWSConnection())

	_, ok := <-s2.Incoming()
	require.False(t, ok)
	require.Error(t, s2.Error())

	<-done
	s.Close()

	require.True(t, s2.Closed())
}

func TestWebSocketStreamDecodeError2(t *testing.T) {
	done := make(chan int, 1)

	s := startTestWSServer(func(conn *websocket.Conn) {
		s1 := NewWebSocketStream(conn)

		// invalid packet
		s1.conn.WriteMessage(websocket.BinaryMessage, []byte{0x20, 0x02, 0x00, 0x06})

		_, ok := <-s1.Incoming()
		require.False(t, ok)
		require.NoError(t, s1.Error())
		require.True(t, s1.Closed())

		done <- 1
	})

	s2 := NewWebSocketStream(newTestWSConnection())

	_, ok := <-s2.Incoming()
	require.False(t, ok)
	require.Error(t, s2.Error())

	<-done
	s.Close()

	require.True(t, s2.Closed())
}
