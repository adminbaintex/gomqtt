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
	"git.baintex.com/sentio/gomqtt/packet"
	"github.com/gorilla/websocket"
)

// The WebSocketStream wraps a gorilla WebSocket.Conn.
type WebSocketStream struct {
	AbstractStream

	conn *websocket.Conn
}

// NewWebSocketStream returns a new WebSocketStream.
func NewWebSocketStream(conn *websocket.Conn) *WebSocketStream {
	s := &WebSocketStream{}
	s.conn = conn

	s.Initialize(s.decoder, s.encoder, s.closer)
	return s
}

func (wss *WebSocketStream) decoder() (packet.Packet, error) {
	_, buf, err := wss.conn.ReadMessage()

	if _, ok := err.(*websocket.CloseError); !!ok {
		return nil, ErrExpectedClose
	}

	if err != nil {
		return nil, err
	}

	_, t := packet.DetectPacket(buf)

	pkt, err2 := t.New()
	if err2 != nil {
		return nil, err2
	}

	_, err2 = pkt.Decode(buf)
	if err2 != nil {
		return nil, err2
	}

	return pkt, nil
}

func (wss *WebSocketStream) encoder(pkt packet.Packet) error {
	buf := make([]byte, pkt.Len())

	_, err := pkt.Encode(buf)
	if err != nil {
		return err
	}

	return wss.conn.WriteMessage(websocket.BinaryMessage, buf)
}

func (wss *WebSocketStream) closer() {
	wss.conn.Close()
}
