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
	"sync"

	"github.com/adminbaintex/gomqtt/packet"
)

// ErrExpectedClose is returned by the Encoder and Decoder functions if they
// encounter an expected end of the underlying stream, eg. EOF.
var ErrExpectedClose = errors.New("Expected close")

// A Decoder function accepts packets and writes them to the underlying stream.
type Decoder func() (packet.Packet, error)

// An Encoder function reads from the underlying stream and returns packets.
type Encoder func(packet.Packet) error

// A Closer function closes the underlying stream.
type Closer func()

// The AbstractStream forms the basis for all streams.
type AbstractStream struct {
	in     chan packet.Packet
	out    chan packet.Packet
	error  error
	closed bool

	decoder Decoder
	encoder Encoder
	closer  Closer

	mutex  sync.Mutex
	start  sync.WaitGroup
	finish sync.WaitGroup
}

// Initialize setups the stream and starts the internal go routines.
func (as *AbstractStream) Initialize(decoder Decoder, encoder Encoder, closer Closer) {
	as.in = make(chan packet.Packet)
	as.out = make(chan packet.Packet)
	as.closed = false
	as.decoder = decoder
	as.encoder = encoder
	as.closer = closer

	as.start.Add(2)
	as.finish.Add(2)

	go as.write()
	go as.read()

	as.start.Wait()
}

// Incoming returns the channel used for reading incoming packets.
// The channel gets automatically closed when the stream gets closed.
func (as *AbstractStream) Incoming() chan packet.Packet {
	return as.in
}

// Outgoing return the channel used for sending outgoing packets.
// The channel gets automatically closed when the stream gets closed.
func (as *AbstractStream) Outgoing() chan packet.Packet {
	return as.out
}

// Send will safely write the packet to the outgoing channel and return true
// if the packet has been written without panicking on a closed channel.
func (as *AbstractStream) Send(pkt packet.Packet) (ret bool) {
	defer func() {
		if recover() != nil {
			ret = false
		}
	}()
	as.out <- pkt
	return true
}

// Error returns the last occurred error. The return value can be consulted
// when the stream gets closed unexpectedly because of an potential error.
func (as *AbstractStream) Error() error {
	return as.error
}

// Close will close the stream and cleanup open channels and running
// go routines.
func (as *AbstractStream) Close() {
	as.exit(nil)
	as.finish.Wait()
}

// Closed will return a boolean indicating if the stream has been
// already closed by Close(), EOF or an error.
func (as *AbstractStream) Closed() bool {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	return as.closed
}

// will close down the stream and eventually the potential error
func (as *AbstractStream) exit(err error) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	if !as.closed {
		as.closed = true
		err = as.error

		if as.closer != nil {
			as.closer()
		}

		close(as.in)
		close(as.out)
	}
}

// write process
func (as *AbstractStream) write() {
	as.start.Done()
	defer as.finish.Done()

	for {
		if as.Closed() {
			return
		}

		pkt := <-as.out

		if pkt != nil {
			err := as.encoder(pkt)

			if err != nil {
				as.error = err
				as.exit(err)
			}
		}
	}
}

// read process
func (as *AbstractStream) read() {
	as.start.Done()
	defer as.finish.Done()

	for {
		if as.Closed() {
			return
		}

		pkt, err := as.decoder()

		if pkt != nil {
			as.send(pkt)
		}

		if err != nil {
			as.error = err
			if err == ErrExpectedClose {
				as.exit(nil)
			} else {
				as.exit(err)
			}
		}
	}
}

// safe write without panicking on a closed channel
func (as *AbstractStream) send(pkt packet.Packet) {
	defer func() { recover() }()
	as.in <- pkt
}
