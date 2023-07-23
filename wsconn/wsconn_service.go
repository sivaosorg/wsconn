package wsconn

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sivaosorg/govm/wsconnx"
)

type WebsocketService interface {
	Run(topic string)
	WriteMessage(conn *websocket.Conn, message wsconnx.WsConnMessagePayload) error
	CloseSubscriber(conn *websocket.Conn)
	AddSubscriber(conn *websocket.Conn, subscription wsconnx.WsConnSubscription)
	SubscribeMessage(c *gin.Context)
	BroadcastMessage(message wsconnx.WsConnMessagePayload)
}

type websocketServiceImpl struct {
	wsConf *Websocket `json:"-"`
}

func NewWebsocketService(wsConf *Websocket) WebsocketService {
	s := &websocketServiceImpl{
		wsConf: wsConf,
	}
	return s
}

func (ws *websocketServiceImpl) Run(topic string) {
	channel, ok := ws.wsConf.Broadcast[topic]
	if !ok {
		_logger.Warn("Topic not found: %v", topic)
		return
	}
	for {
		message := <-channel
		ws.wsConf.Mutex.Lock()
		for subscriber, subscription := range ws.wsConf.Subscribers {
			if subscription.Topic == topic {
				err := ws.WriteMessage(subscriber, message)
				if err != nil {
					_logger.Error("An error occurred while writing message", err)
					ws.CloseSubscriber(subscriber)
				}
			}
		}
		ws.wsConf.Mutex.Unlock()
	}
}

func (ws *websocketServiceImpl) WriteMessage(conn *websocket.Conn, message wsconnx.WsConnMessagePayload) error {
	message.SetGenesisTimestamp(time.Now())
	conn.SetWriteDeadline(time.Now().Add(ws.wsConf.Config.WriteWait))
	return conn.WriteJSON(message)
}

func (ws *websocketServiceImpl) CloseSubscriber(conn *websocket.Conn) {
	conn.Close()
	delete(ws.wsConf.Subscribers, conn)
}

func (ws *websocketServiceImpl) AddSubscriber(conn *websocket.Conn, subscription wsconnx.WsConnSubscription) {
	ws.wsConf.Mutex.Lock()
	defer ws.wsConf.Mutex.Unlock()
	if conn != nil {
		ws.wsConf.Subscribers[conn] = subscription
		if _, ok := ws.wsConf.Broadcast[subscription.Topic]; !ok {
			ws.wsConf.Broadcast[subscription.Topic] = make(chan wsconnx.WsConnMessagePayload)
			go ws.Run(subscription.Topic)
		}
	}
}

// Parse user ID and desired topic from the WebSocket message
// Read incoming messages, but ignore them as we handle sending only
func (ws *websocketServiceImpl) SubscribeMessage(c *gin.Context) {
	conn, err := ws.wsConf.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		_logger.Error("An error occurred while upgrading connection", err)
		return
	}
	defer ws.CloseSubscriber(conn)
	var subscription wsconnx.WsConnSubscription
	if err := conn.ReadJSON(&subscription); err != nil {
		_logger.Error("An error occurred while reading subscription", err)
		return
	}
	ws.AddSubscriber(conn, subscription)
	if ws.wsConf.AllowCloseConn {
		conn.SetReadLimit(int64(ws.wsConf.Config.MaxMessageSize))
		conn.SetReadDeadline(time.Now().Add(ws.wsConf.Config.PongWait))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(ws.wsConf.Config.PongWait))
			return nil
		})
	}
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				_logger.Error("An error occurred while reading message", err)
			}
			break
		}
	}
}

// Find the channel for the specific topic and send the message to it
func (ws *websocketServiceImpl) BroadcastMessage(message wsconnx.WsConnMessagePayload) {
	ws.wsConf.Mutex.Lock()
	defer ws.wsConf.Mutex.Unlock()
	if channel, ok := ws.wsConf.Broadcast[message.Topic]; ok {
		channel <- message
	}
}
