package wsconn

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sivaosorg/govm/wsconnx"
)

type Websocket struct {
	Config           wsconnx.WsConnOptionConfig                     `json:"conf"`
	AllowCloseConn   bool                                           `json:"allow_close_conn"`
	RegisteredTopics map[string]bool                                `json:"registered_topics"`
	Upgrader         websocket.Upgrader                             `json:"-"`
	Broadcast        map[string]chan wsconnx.WsConnMessagePayload   `json:"-"`
	Subscribers      map[*websocket.Conn]wsconnx.WsConnSubscription `json:"-"`
	Mutex            sync.Mutex                                     `json:"-"`
}
