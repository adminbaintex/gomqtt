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
	"fmt"
	"io"

	"github.com/gomqtt/packet"
)

// EncodeToWriter encodes a writes a packet to a bufio.Writer.
// It returns the bytes written to the writer and eventual errors.
func EncodeToWriter(w *bufio.Writer, pkt packet.Packet) (int, error) {
	buf := make([]byte, pkt.Len())

	// encode packet
	n, err := pkt.Encode(buf)
	if err != nil {
		return 0, err
	}

	// write to writer
	n, err = w.Write(buf)
	if err != nil {
		return n, err
	}

	// flush writer
	err = w.Flush()
	if err != nil {
		return n, err
	}

	return n, nil
}

// DecodeFromReader will read and decode the next packet from a bufio.Reader
// and return it as well as the bytes read and eventual errors.
func DecodeFromReader(r *bufio.Reader) (packet.Packet, int, error) {
	// initial detection length
	l := 2

	for {
		// check length
		if l > 5 {
			return nil, 0, fmt.Errorf("DecodeFrom: Error while detecting next packet.")
		}

		// try read detection bytes
		d, err := r.Peek(l)
		if err != nil {
			return nil, 0, err
		}

		// detect packet
		ml, t := packet.DetectPacket(d)

		// on zero packet length: increment detection length and try again
		if ml <= 0 {
			l++
			continue
		}

		// allocate buffer
		m := make([]byte, ml)

		// read whole packet
		n, err := io.ReadFull(r, m)
		if err != nil {
			return nil, n, err
		}

		// create packet
		pkt, err := t.New()
		if err != nil {
			return nil, n, err
		}

		// decode buffer
		_, err = pkt.Decode(m)
		if err != nil {
			return nil, n, err
		}

		// return packet
		return pkt, n, nil
	}
}
