package websocket

import (
    "encoding/json"
    "github.com/gin-gonic/gin"
    "github.com/gofrs/uuid"
    "github.com/gorilla/websocket"
    "github.com/sqzxcv/glog"
    "net/http"
    "time"
)

// ClientManager is a websocket manager
type ClientManager struct {
    Clients    map[*Client]bool
    Broadcast  chan []byte
    Register   chan *Client
    Unregister chan *Client
}

// Client is a websocket client
type Client struct {
    ID     string
    Socket *websocket.Conn
    Send   chan []byte
}

// Message is an object for websocket message which is mapped to json type
type Message struct {
    Sender    string `json:"sender,omitempty"`
    Recipient string `json:"recipient,omitempty"`
    Content   string `json:"content,omitempty"`
}

// Manager define a ws server manager
var Manager = ClientManager{
    Broadcast:  make(chan []byte),
    Register:   make(chan *Client),
    Unregister: make(chan *Client),
    Clients:    make(map[*Client]bool),
}

// Start is to start a ws server
func (manager *ClientManager) Start() {
    for {
        select {
        case conn := <-manager.Register:
            // 这地方做特殊设置, 同一时间只允许一个客户端连接当前ws, 所以有新的连接建立的时候, 会关闭之前的链接
            for diedConn, _ := range manager.Clients {
                close(diedConn.Send)
            }
            manager.Clients = make(map[*Client]bool)
            manager.Clients[conn] = true

        case conn := <-manager.Unregister:
            if _, ok := manager.Clients[conn]; ok {
                close(conn.Send)
                delete(manager.Clients, conn)
            }
        case message := <-manager.Broadcast:
            for conn := range manager.Clients {
                select {
                case conn.Send <- message:
                default:
                    close(conn.Send)
                    delete(manager.Clients, conn)
                }
            }
        }
    }
}

// Send is to send ws message to ws client
func (manager *ClientManager) Send(message []byte, ignore *Client) {
    for conn := range manager.Clients {
        if conn != ignore {
            conn.Send <- message
        }
    }
}

func (c *Client) Read() {
    defer func() {
        Manager.Unregister <- c
        c.Socket.Close()
    }()

    for {
        _, message, err := c.Socket.ReadMessage()
        if err != nil {
            Manager.Unregister <- c
            c.Socket.Close()
            break
        }
        jsonMessage, _ := json.Marshal(&Message{Sender: c.ID, Content: string(message)})
        Manager.Broadcast <- jsonMessage
    }
}

func (c *Client) Write() {
    defer func() {
        c.Socket.Close()
    }()

    for {
        select {
        case message, ok := <-c.Send:
            if !ok {
                c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            c.Socket.WriteMessage(websocket.TextMessage, message)
        }
    }
}

// WsPage is a websocket handler
// 将普通http请求升级为websocket协议
func WsPage(c *gin.Context) {
    // change the reqest to websocket model
    conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
    if error != nil {
        http.NotFound(c.Writer, c.Request)
        return
    }
    u4, err := uuid.NewV4()
    // websocket connect
    if err != nil {
        glog.Error("创建uuid失败, 原因:", err)
        http.NotFound(c.Writer, c.Request)
        return
    }
    client := &Client{ID: u4.String(), Socket: conn, Send: make(chan []byte)}

    Manager.Register <- client

    go client.Read()
    go client.Write()
}

//WSBroadcast 向客户端广播json消息.
func WSBroadcast(cmd string, message string) {
    hasClient := false
    for _, b := range Manager.Clients {
        if b {
            hasClient = b
            break
        }
    }
    if hasClient {
        msg := make(map[string]string)
        msg["cmd"] = cmd
        msg["msg"] = message
        content, _ := json.Marshal(msg)

        Manager.Broadcast <- content
    }
}

// 心跳机制: 服务器每隔一段时间广播一条心跳包. 客户端检测心跳包,如果长时间没收到心跳包, 客户端则主动端口链接重连
func heartBeat() {

    WSBroadcast("beat", "beat")
    time.AfterFunc(time.Second*20, func() {
        heartBeat()
    })
}
