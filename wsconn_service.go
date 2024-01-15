package wsconn

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sivaosorg/govm/entity"
	"github.com/sivaosorg/govm/wsconnx"
)

type WebsocketService interface {
	Run(topic string)
	WriteMessage(conn *websocket.Conn, message wsconnx.WsConnMessagePayload) error
	CloseSubscriber(conn *websocket.Conn)
	AddSubscriber(conn *websocket.Conn, subscription wsconnx.WsConnSubscription)
	SubscribeMessage(c *gin.Context)
	BroadcastMessage(message wsconnx.WsConnMessagePayload)
	RegisterTopic(c *gin.Context)
}

type websocketServiceImpl struct {
	wsConf *Websocket
}

func NewWebsocketService(wsConf *Websocket) WebsocketService {
	s := &websocketServiceImpl{
		wsConf: wsConf,
	}
	return s
}

func (ws *websocketServiceImpl) Run(topic string) {
	channel, ok := ws.wsConf.broadcast[topic]
	if !ok {
		_logger.Warn("Topic not found: %v", topic)
		return
	}
	for {
		message := <-channel
		ws.wsConf.mutex.Lock()
		for subscriber, subscription := range ws.wsConf.subscribers {
			if subscription.Topic == topic {
				err := ws.WriteMessage(subscriber, message)
				if err != nil {
					_logger.Error("An error occurred while writing message", err)
					ws.CloseSubscriber(subscriber)
				}
			}
		}
		ws.wsConf.mutex.Unlock()
	}
}

func (ws *websocketServiceImpl) WriteMessage(conn *websocket.Conn, message wsconnx.WsConnMessagePayload) error {
	message.SetGenesisTimestamp(time.Now())
	conn.SetWriteDeadline(time.Now().Add(ws.wsConf.Option.WriteWait))
	return conn.WriteJSON(message)
}

func (ws *websocketServiceImpl) CloseSubscriber(conn *websocket.Conn) {
	conn.Close()
	delete(ws.wsConf.subscribers, conn)
}

func (ws *websocketServiceImpl) AddSubscriber(conn *websocket.Conn, subscription wsconnx.WsConnSubscription) {
	ws.wsConf.mutex.Lock()
	defer ws.wsConf.mutex.Unlock()
	if conn != nil {
		ws.wsConf.subscribers[conn] = subscription
		if _, ok := ws.wsConf.broadcast[subscription.Topic]; !ok {
			ws.wsConf.broadcast[subscription.Topic] = make(chan wsconnx.WsConnMessagePayload)
			go ws.Run(subscription.Topic)
		}
	}
}

// Parse user ID and desired topic from the WebSocket message
// Read incoming messages, but ignore them as we handle sending only
func (ws *websocketServiceImpl) SubscribeMessage(c *gin.Context) {
	conn, err := ws.wsConf.upgrader.Upgrade(c.Writer, c.Request, nil)
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
	if ws.wsConf.IsEnabledClosure {
		conn.SetReadLimit(int64(ws.wsConf.Option.MaxMessageSize))
		conn.SetReadDeadline(time.Now().Add(ws.wsConf.Option.PongWait))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(ws.wsConf.Option.PongWait))
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
	ws.wsConf.mutex.Lock()
	defer ws.wsConf.mutex.Unlock()
	if channel, ok := ws.wsConf.broadcast[message.Topic]; ok {
		channel <- message
	}
}

func (ws *websocketServiceImpl) RegisterTopic(c *gin.Context) {
	ws.wsConf.mutex.Lock()
	defer ws.wsConf.mutex.Unlock()
	response := entity.NewResponseEntity()
	var subscription wsconnx.WsConnSubscription
	if err := c.ShouldBindJSON(&subscription); err != nil {
		response.SetStatusCode(http.StatusBadRequest).SetError(err).SetMessage(err.Error())
		c.JSON(response.StatusCode, response)
		return
	}
	if _, ok := ws.wsConf.Topics[subscription.Topic]; ok {
		response.SetStatusCode(http.StatusOK).SetMessage(fmt.Sprintf("Topic %s already registered", subscription.Topic)).SetData(subscription)
		c.JSON(response.StatusCode, response)
		return
	}
	ws.wsConf.Topics[subscription.Topic] = true
	ws.wsConf.broadcast[subscription.Topic] = make(chan wsconnx.WsConnMessagePayload)
	go ws.Run(subscription.Topic)
	response.SetStatusCode(http.StatusOK).SetMessage(fmt.Sprintf("Topic %s registered successfully", subscription.Topic)).SetData(subscription)
	c.JSON(response.StatusCode, response)
	return
}
