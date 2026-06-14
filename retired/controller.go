package Main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Controller struct {
	conn      net.Conn
	connected bool
	ID        string
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
			ID:        genRandomID(),
		}

		fmt.Println("established connection:", c.ID)
		go c.HandleConnection()

	}
}

func (c *Controller) HandleConnection() {
	defer func() {
		c.conn.Close()
		fmt.Println("connection terminated:", c.ID)
		c.connected = false
	}()
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {

		message := scanner.Text()

		if message == "ping" {
			c.connected = true

			continue
		}

		fmt.Println(len(message), message)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
	}
}
