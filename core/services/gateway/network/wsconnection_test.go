package network_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	gw_net "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

var upgrader = websocket.Upgrader{}

type serverSideLogic struct {
	readCh        chan gw_net.ReadItem
	writeCh       chan gw_net.WriteItem
	wsConnWrapper *gw_net.WSConnectionWrapper
}

func (ssl *serverSideLogic) wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// one wsConnWrapper per client
	ssl.wsConnWrapper.Restart(c)
}

func TestWSConnectionWrapper_ClientReconnect(t *testing.T) {
	// server
	ssl := &serverSideLogic{
		readCh:  make(chan gw_net.ReadItem),
		writeCh: make(chan gw_net.WriteItem),
	}
	ssl.wsConnWrapper = gw_net.NewWSConnectionWrapper(ssl.readCh, ssl.writeCh)
	s := httptest.NewServer(http.HandlerFunc(ssl.wsHandler))
	serverURL := "ws" + strings.TrimPrefix(s.URL, "http")
	defer s.Close()

	// client
	clientReadCh := make(chan gw_net.ReadItem)
	clientWriteCh := make(chan gw_net.WriteItem)
	clientConnWrapper := gw_net.NewWSConnectionWrapper(clientReadCh, clientWriteCh)
	devNullCh := make(chan error, 100)

	// connect, write a message, disconnect
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	require.NoError(t, err)
	clientConnWrapper.Restart(conn)
	clientWriteCh <- gw_net.WriteItem{websocket.TextMessage, []byte("hello"), devNullCh}
	<-ssl.readCh // consumed by server
	conn.Close()

	// re-connect, write another message, disconnect
	conn, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
	require.NoError(t, err)
	clientConnWrapper.Restart(conn)
	clientWriteCh <- gw_net.WriteItem{websocket.TextMessage, []byte("hello again"), devNullCh}
	<-ssl.readCh // consumed by server
	conn.Close()

	ssl.wsConnWrapper.Shutdown()
	clientConnWrapper.Shutdown()
}
