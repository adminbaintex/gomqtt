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

// Package server implements basic functionality for launching an MQTT 3.1.1
// (http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/) server that use multiple
// listeners and protocols.
package server

import (
	"net"
	"net/http"

	"github.com/gomqtt/stream"
	"github.com/gorilla/websocket"
)

// Handler will receive new connections as streams.
type Handler func(stream stream.Stream)

// A Configuration describes a single server accepting connections using one
// protocol and address.
type Configuration struct {
	Protocol string
	Address string

	listener net.Listener
}

// The Server manages multiple Configurations and yields new connection as
// streams to the Handler.
type Server struct {
	// The Handler that receives new Streams.
	Handler Handler

	// The currently running configurations.
	Configurations []Configuration

	// The channel used to send internal errors.
	Error chan error
}

// NewServer returns a new Server.
func NewServer(handler Handler) *Server {
	return &Server{
		Handler: handler,
		Configurations: make([]Configuration, 0),
	}
}

// LaunchTCPConfiguration will run a simple TCP server.
func (s *Server) LaunchTCPConfiguration(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.Configurations = append(s.Configurations, Configuration{
		Protocol: "tcp",
		Address: address,
		listener: l,
	})

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				s.Error <- err
				return
			}

			go s.Handler(stream.NewNetStream(conn))
		}
	}()

	return nil
}

// LaunchWSConfiguration will run a simple WS (HTTP) server.
func (s *Server) LaunchWSConfiguration(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.Error <- err
		}

		go s.Handler(stream.NewWebSocketStream(conn))
	})

	h := &http.Server{
		Handler: mux,
	}

	s.Configurations = append(s.Configurations, Configuration{
		Protocol: "ws",
		Address: address,
		listener: l,
	})

	go h.Serve(l)

	return nil
}

// Stop will stop all started configurations.
// Eventual errors are sent via the Error channel.
func (s *Server) Stop() {
	for _, c := range s.Configurations {
		err := c.listener.Close()
		if err != nil {
			s.Error <- err
		}
	}
}
