package network

import (
	"errors"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

// WSConnectionWrapper is a websocket connection abstraction that supports re-connects.
// I/O is separated from connection management:
//   - component doing I/O provides read and write channels that never change
//   - component managing connections can listen to connection-closed events and call Restart()
//     to swap the underlying connection object
//
// The Wrapper can be used by a server expecting long-lived connections from a given client,
// as well as a client maintaining such long-lived connection with a given server.
// This fits the Gateway very well as servers accept connections only from a fixed set of nodes
// and conversely, nodes only connect to a fixed set of servers (Gateways).
//
// The concept of "pumps" is borrowed from https://github.com/smartcontractkit/wsrpc
type WSConnectionWrapper struct {
	conn atomic.Pointer[websocket.Conn]

	writeCh    <-chan WriteItem
	readCh     chan<- ReadItem
	shutdownCh chan struct{}
}

type ReadItem struct {
	MsgType int
	Data    []byte
}

type WriteItem struct {
	MsgType int
	Data    []byte
	ErrCh   chan error
}

func NewWSConnectionWrapper(readCh chan<- ReadItem, writeCh <-chan WriteItem) *WSConnectionWrapper {
	cw := &WSConnectionWrapper{
		writeCh:    writeCh,
		readCh:     readCh,
		shutdownCh: make(chan struct{}),
	}
	// write pump runs until Shutdown() is called
	go cw.writePump()
	return cw
}

func (c *WSConnectionWrapper) Shutdown() {
	c.shutdownCh <- struct{}{}
	c.Restart(nil)
}

// Restart
//  1. replaces the underlying connection and shuts the old one down
//  2. starts a new read goroutine that pushes received messages to readCh
//  3. returns channel that closes when connection closes
func (c *WSConnectionWrapper) Restart(newConn *websocket.Conn) <-chan struct{} {
	oldConn := c.conn.Swap(newConn)

	if oldConn != nil {
		oldConn.Close()
	}
	if newConn == nil {
		return nil
	}
	closeCh := make(chan struct{})
	// readPump goroutine is tied to the lifecycle of the underlying conn object
	go readPump(newConn, c.readCh, closeCh)
	return closeCh
}

func (c *WSConnectionWrapper) writePump() {
	for {
		select {
		case wsMsg := <-c.writeCh:
			// synchronization is a tradeoff for the ability to use a single write channel
			conn := c.conn.Load()
			if conn == nil {
				wsMsg.ErrCh <- errors.New("no active connection")
				break
			}
			wsMsg.ErrCh <- conn.WriteMessage(wsMsg.MsgType, wsMsg.Data)
		case <-c.shutdownCh:
			return
		}
	}
}

func readPump(conn *websocket.Conn, readCh chan<- ReadItem, closeCh chan<- struct{}) {
	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			close(closeCh)
			return
		}
		readCh <- ReadItem{msgType, data}
	}
}
