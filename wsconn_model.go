package wsconn

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sivaosorg/govm/wsconnx"
)

type Websocket struct {
	mutex            sync.Mutex
	upgrader         websocket.Upgrader
	broadcast        map[string]chan wsconnx.WsConnMessagePayload
	subscribers      map[*websocket.Conn]wsconnx.WsConnSubscription
	Option           wsconnx.WsConnOptionConfig `json:"option"`
	IsEnabledClosure bool                       `json:"enabled_closure"`
	Topics           map[string]bool            `json:"topics"`
}
