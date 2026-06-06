package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Controller struct {
	conn      net.Conn
	connected bool
}

func setUpListener() {
	listener, err := net.Listen("tcp", "0.0.0.0:8090")
	if err != nil {
		log.Fatal("error setting up tcp server")
	}
	defer listener.Close()

	fmt.Println("tcp server up")

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("error establishing connection", err)
			continue
		}

		c := Controller{
			conn:      conn,
			connected: false,
		}

		go c.HandleConnection()

	}
}

func (c *Controller) HandleConnection() {
	defer c.conn.Close()
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {

		message := scanner.Text()

		if message == "ping" {
			c.connected = true
			continue
		}

		fmt.Println("message recieved:", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
	}
	fmt.Println("Connection terminated")
	c.connected = false
}
