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

package server

import (
	"fmt"
	"net"

	"git.baintex.com/sentio/gomqtt/stream"
)

func ExampleServer() {
	done := make(chan struct{})

	server := NewServer(func(s stream.Stream){
		fmt.Println("new connection")

		<-s.Error()
		close(done)
	})

	server.LaunchTCPConfiguration("localhost:1337")
	server.LaunchWSConfiguration("localhost:1338")

	c, _ := net.Dial("tcp", "localhost:1337")
	c.Close()

	<-done
	server.Stop()

	// Output:
	// new connection
}
