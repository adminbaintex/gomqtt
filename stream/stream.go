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

// Package stream implements functionality for handling MQTT 3.1.1
// (http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/) packet streams.
package stream

import "github.com/gomqtt/packet"

// Stream is an abstract interface for all transport streams.
type Stream interface {
	// Incoming returns the channel used for reading incoming packets.
	// The channel gets automatically closed when the stream gets closed.
	Incoming() chan packet.Packet

	// Outgoing return the channel used for sending outgoing packets.
	// The channel gets automatically closed when the stream gets closed.
	Outgoing() chan packet.Packet

	// Send will safely write the packet to the outgoing channel and return true
	// if the packet has been written without panicking on a closed channel.
	Send(packet.Packet) bool

	// Error returns the last occurred error. The return value can be consulted
	// when the stream gets closed unexpectedly because of an potential error.
	Error() error

	// Close will close the stream and cleanup open channels and running
	// go routines.
	Close()

	// Closed will return a boolean indicating if the stream has been
	// already closed by Close(), EOF or an error.
	Closed() bool
}
