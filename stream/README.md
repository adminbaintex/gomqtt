# gomqtt/stream

[![Circle CI](https://img.shields.io/circleci/project/gomqtt/stream.svg)](https://circleci.com/gh/gomqtt/stream)
[![Coverage Status](https://coveralls.io/repos/gomqtt/stream/badge.svg?branch=master&service=github)](https://coveralls.io/github/gomqtt/stream?branch=master)
[![GoDoc](https://godoc.org/github.com/gomqtt/stream?status.svg)](http://godoc.org/github.com/gomqtt/stream)
[![Release](https://img.shields.io/github/release/gomqtt/stream.svg)](https://github.com/gomqtt/stream/releases)
[![Go Report Card](http://goreportcard.com/badge/gomqtt/stream)](http://goreportcard.com/report/gomqtt/stream)

**Package stream implements functionality for handling [MQTT 3.1.1](http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/) packet streams.**

## Features

- Support for basic `io.Reader` and `io.Writer` based streams.
- Support for `net.Conn` based streams.
- Support for `websocket.Conn` based streams through <https://github.com/gorilla/websocket>.

## Installation

Get it using go's standard toolset:

```bash
$ go get github.com/gomqtt/stream
```

## Usage

```go
done := make(chan int, 1)

s := startTestTCPServer(func(conn net.Conn) {
    s1 := NewNetStream(conn)

    for {
        select{
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
```
