package wsconn

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sivaosorg/govm/logger"
	"github.com/sivaosorg/govm/utils"
	"github.com/sivaosorg/govm/wsconnx"
)

var (
	_logger    = logger.NewLogger()
	conf       = wsconnx.GetWsConnOptionConfigSample()
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func NewWebsocket() *Websocket {
	ws := &Websocket{}
	ws.SetBroadcast(make(map[string]chan wsconnx.WsConnMessagePayload))
	ws.SetSubscribers(make(map[*websocket.Conn]wsconnx.WsConnSubscription))
	ws.SetConfig(*conf)
	ws.SetUpgrader(wsUpgrader)
	ws.SetAllowCloseConn(false)
	return ws
}

func (ws *Websocket) SetBroadcast(value map[string]chan wsconnx.WsConnMessagePayload) *Websocket {
	ws.Broadcast = value
	return ws
}

func (ws *Websocket) SetSubscribers(value map[*websocket.Conn]wsconnx.WsConnSubscription) *Websocket {
	ws.Subscribers = value
	return ws
}

func (ws *Websocket) SetConfig(value wsconnx.WsConnOptionConfig) *Websocket {
	ws.Config = value
	return ws
}

func (ws *Websocket) SetUpgrader(value websocket.Upgrader) *Websocket {
	ws.Upgrader = value
	return ws
}

func (ws *Websocket) SetAllowCloseConn(value bool) *Websocket {
	ws.AllowCloseConn = value
	return ws
}

func (ws *Websocket) Json() string {
	return utils.ToJson(ws)
}
