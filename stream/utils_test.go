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
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gomqtt/packet"
)

// the errorWriter
type errorWriter struct{}

// Write will return always an error
func (ew *errorWriter) Write(p []byte) (int, error) {
	return 0, errors.New("error")
}

// the infiniteReader
type infiniteReader struct{}

// Read will block for 10us and read zero bytes
func (ir *infiniteReader) Read(p []byte) (int, error) {
	time.Sleep(10 * time.Microsecond)

	return 0, nil
}

// the messageReader
type messageReader struct{
	data []byte
}

// creates a new messageReader
func newMessageReader(p packet.Packet) *messageReader {
	mr := &messageReader{
		data: make([]byte, p.Len()),
	}
	p.Encode(mr.data)
	return mr
}

// Read will return the stored message
func (mr *messageReader) Read(p []byte) (int, error) {
	copy(p, mr.data)
	return len(mr.data), nil
}

// will encode the packet and return the buffer
func encodePacket(p packet.Packet) []byte {
	b := make([]byte, p.Len())
	p.Encode(b)
	return b
}

// runs tcp server, yields connections and returns listener
func startTestTCPServer(handler func(net.Conn)) net.Listener {
	l, err := net.Listen("tcp", "localhost:3456")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}

			go handler(conn)
		}
	}()

	return l
}

// creates a tcp connection and returns it
func newTestTCPConnection() net.Conn {
	conn, err := net.Dial("tcp", "localhost:3456")
	if err != nil {
		panic(err)
	}

	return conn
}

// runs http server, yield connections and returns listener
func startTestWSServer(handler func(*websocket.Conn)) net.Listener {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}

		go handler(conn)
	})

	s := &http.Server{
		Handler: mux,
	}

	l, err := net.Listen("tcp", "localhost:2345")
	if err != nil {
		panic(err)
	}

	go s.Serve(l)
	return l
}

// creates a tcp connection and returns it
func newTestWSConnection() *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:2345/", nil)
	if err != nil {
		panic(err)
	}

	return conn
}
