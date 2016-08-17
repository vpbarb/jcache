package transport

import (
	"bufio"
	"log"
	"net"
	"strings"
)

type Request struct {
	Name   string
	Params []string
}

type Response []byte

type handler interface {
	Handle(*Request) Response
}

func ListenTCP(h handler, addr string) {
	listener, _ := net.Listen("tcp", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error on connection accepting: %s\n", err)
			continue
		}

		log.Print("incoming tcp connect")

		go func(c net.Conn) {
			buf := bufio.NewReader(c)
			for {
				line, _, err := buf.ReadLine()

				if err != nil {
					break
				}
				if len(line) > 0 {
					log.Printf("read: %s", string(line))

					fields := strings.Fields(string(line))

					r := &Request{Name: fields[0]}
					if len(fields) > 0 {
						r.Params = fields[1:]
					}

					if response := h.Handle(r); len(response) > 0 {
						c.Write(response)
						c.Write([]byte("\n"))
					}
				}
			}
			log.Print("close tcp connect")
		}(conn)
	}
}
