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
	"fmt"
	"net"

	"github.com/adminbaintex/gomqtt/packet"
)

func ExampleStream() {
	done := make(chan int, 1)

	s := startTestTCPServer(func(conn net.Conn) {
		s1 := NewNetStream(conn)

		for {
			select {
			case p, ok := <-s1.Incoming():
				if p != nil {
					switch p.Type() {
					case packet.CONNECT:
						c, _ := p.(*packet.ConnectPacket)
						fmt.Println(c)

						s1.Send(packet.NewConnackPacket())
					}
				}

				if !ok {
					done <- 1
				}
			}
		}
	})

	s2 := NewNetStream(newTestTCPConnection())
	s2.Send(packet.NewConnectPacket())

	m := <-s2.Incoming()

	switch m.Type() {
	case packet.CONNACK:
		c, _ := m.(*packet.ConnackPacket)
		fmt.Println(c)

		s2.Close()
	}

	<-done
	s.Close()

	// Output:
	// CONNECT: ClientID="" KeepAlive=0 Username="" Password="" CleanSession=true WillTopic="" WillPayload="" WillQOS=0 WillRetain=false
	// CONNACK: SessionPresent=false ReturnCode="Connection accepted"
}
