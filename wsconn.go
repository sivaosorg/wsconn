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
	w := &Websocket{}
	w.SetBroadcast(make(map[string]chan wsconnx.WsConnMessagePayload))
	w.SetSubscribers(make(map[*websocket.Conn]wsconnx.WsConnSubscription))
	w.SetOption(*conf)
	w.SetUpgrader(wsUpgrader)
	w.SetEnabledClosure(false)
	w.SetTopics(make(map[string]bool))
	return w
}

func (w *Websocket) SetBroadcast(value map[string]chan wsconnx.WsConnMessagePayload) *Websocket {
	w.broadcast = value
	return w
}

func (w *Websocket) SetSubscribers(value map[*websocket.Conn]wsconnx.WsConnSubscription) *Websocket {
	w.subscribers = value
	return w
}

func (w *Websocket) SetOption(value wsconnx.WsConnOptionConfig) *Websocket {
	w.Option = value
	return w
}

func (w *Websocket) SetUpgrader(value websocket.Upgrader) *Websocket {
	w.upgrader = value
	return w
}

func (w *Websocket) Upgrader() websocket.Upgrader {
	return w.upgrader
}

func (w *Websocket) SetEnabledClosure(value bool) *Websocket {
	w.IsEnabledClosure = value
	return w
}

func (w *Websocket) SetTopics(value map[string]bool) *Websocket {
	w.Topics = value
	return w
}

func (w *Websocket) Json() string {
	return utils.ToJson(w)
}
