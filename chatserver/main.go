package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const welcome = `
   ____ _   _    _  _____  __
  / ___| | | |  / \|_   _|/ ___ _ ____   __
 | |   | |_| | / _ \ | | / / __| '__\ \ / /
 | |___|  _  |/ ___ \| |/ /\__ | |   \ V /
  \____|_| |_/_/   \_|_/_/ |___|_|    \_/
  		connected
..

type '/name username' to name yourself
`

type Connection struct {
	conn *net.Conn
	name string
}

var cons = []*Connection{}

func main() {
	fmt.Println("Server starting...")

	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	i := 1
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("err accepting = %+v\n", err)
		}

		i = i + 1
		id := fmt.Sprintf("anonymous-%d", i)
		c := &Connection{conn: &conn, name: id}
		cons = append(cons, c)

		go func(conn net.Conn) {
			sendConns()

			_, err := conn.Write([]byte(welcome + "\r"))
			if err != nil {
				log.Printf("error writing: %v", err)
			}

			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				if strings.HasPrefix(scanner.Text(), "/name") {
					name := strings.TrimPrefix(scanner.Text(), "/name ")
					c.name = name
					conn.Write([]byte("name set to " + c.name + "\n\r"))
					sendConns()
					continue
				}

				writeToAll(c.name + ": " + scanner.Text() + "\n")
			}

			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "error reading input:", err)
			}
		}(conn)
	}
}

func sendConns() {
	names := []string{}
	for _, c := range cons {
		names = append(names, c.name)
	}
	writeToAll("u:\n" + strings.Join(names, "\n") + "\n")
}

func writeToAll(body string) {
	for i, c := range cons {
		_, err := (*c.conn).Write([]byte(body + "\r"))
		if err != nil {
			cons = append(cons[:i], cons[i+1:]...)
		}
	}
}
