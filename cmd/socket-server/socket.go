//go:build !js && !wasm

package main

import (
	"context"
	"exnes/Core"
	"fmt"
	"math/rand"
	"net/http"
	"slices"
	"time"

	"github.com/coder/websocket"
)

var connectedClients []*client

type client struct {
	ID     string
	Output chan Core.ScreenInfo
	conn   *websocket.Conn
}

func acceptScreenConn(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if err != nil {
		http.Error(w, "bad connection try", http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	defer conn.CloseNow()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	ctx = conn.CloseRead(ctx)

	c := &client{
		ID:     genRandomID(),
		conn:   conn,
		Output: make(chan Core.ScreenInfo, 60),
	}

	connectedClients = append(connectedClients, c)
	defer removeClient(c)

	if debugConsole.Console != nil {
		c.runReciever(ctx)
	}

}

func (c *client) runReciever(ctx context.Context) {

	// debugConsole.Console.RunDisplayUpdates()
	// c.conn.Write(ctx, websocket.MessageBinary, debugConsole.Console.Ppu.BackBuffer[:])
	fmt.Println("opening conn:", c.ID)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("closing connection", c.ID)
			return
		case info := <-c.Output:
			c.conn.Write(ctx, websocket.MessageBinary, (*info.Buffer)[:])
		}
	}

}

func HandleScreenUpdates() {
	for s := range debugConsole.Console.ScreenChannel {
		for _, client := range connectedClients {
			client.Output <- s
		}
	}
}

func removeClient(c *client) {
	for id, conn := range connectedClients {
		if c == conn {
			connectedClients = slices.Delete(connectedClients, id, id+1)
		}
	}
}

func genRandomID() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 5)
	for i := range b {
		b[i] = chars[r.Intn(len(chars))]
	}
	return string(b)
}
