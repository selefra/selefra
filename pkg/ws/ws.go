package ws

import (
	"github.com/gorilla/websocket"
	"github.com/selefra/selefra/global"
	"sync"
	"time"
)

type connClient struct {
	conn   *websocket.Conn
	l      sync.Mutex
	ID     string
	Token  string
	TaskId string
	Remote string
}

func (c *connClient) WriteJsonLock(v any) error {
	c.l.Lock()
	defer c.l.Unlock()
	err := c.conn.WriteJSON(v)
	return err
}

type BaseConnectionInfo struct {
	ID         string
	Token      string
	TaskId     string
	RemoteType string
}

type connectMsg struct {
	ActionName string             `json:"action_name"`
	Data       interface{}        `json:"data"`
	Msg        string             `json:"msg"`
	BaseInfo   BaseConnectionInfo `json:"base_info"`
}

const LogStream = "logStream"
const Register = "register"
const Ping = "ping"
const Reconnect = "reconnect"
const TaskCompleted = "task_completed"

var Client connClient
var registerSuccess bool

func (c *connClient) Close() error {
	registerSuccess = false
	return c.conn.Close()
}

func Init() {
	di := websocket.Dialer{}
	conn, _, err := di.Dial("wss://"+global.SERVER+"/cli/ws/log_stream", nil)
	if err != nil {
		return
	}
	Client.conn = conn
	go onMessage()
}

func SendLog(msg string) error {
	if registerSuccess {
		msg := connectMsg{
			ActionName: LogStream,
			Data:       msg,
			Msg:        "",
			BaseInfo: BaseConnectionInfo{
				ID:         "",
				Token:      Client.Token,
				TaskId:     Client.TaskId,
				RemoteType: "cli",
			},
		}
		err := Client.WriteJsonLock(msg)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func Regis(token, taskId string) error {
	msg := connectMsg{
		ActionName: Register,
		Data:       nil,
		Msg:        "",
		BaseInfo: BaseConnectionInfo{
			ID:         "",
			Token:      token,
			TaskId:     taskId,
			RemoteType: "cli",
		},
	}
	Client.Token = token
	Client.TaskId = taskId
	registerSuccess = true
	err := Client.WriteJsonLock(msg)
	if err != nil {
		return err
	}
	go PingF()
	return nil
}

func PingF() {
	for {
		time.Sleep(5 * time.Second)
		msg := connectMsg{
			ActionName: Ping,
			Data:       nil,
			Msg:        "",
			BaseInfo: BaseConnectionInfo{
				ID:         "",
				Token:      Client.Token,
				TaskId:     Client.TaskId,
				RemoteType: "cli",
			},
		}
		err := Client.WriteJsonLock(msg)
		if err != nil {
			return
		}
	}
}

func Completed() error {
	msg := connectMsg{
		ActionName: TaskCompleted,
		Data:       nil,
		Msg:        "",
		BaseInfo: BaseConnectionInfo{
			ID:         "",
			Token:      Client.Token,
			TaskId:     Client.TaskId,
			RemoteType: "cli",
		},
	}
	err := Client.WriteJsonLock(msg)
	if err != nil {
		return err
	}
	return nil
}

func onMessage() {
	for {
		_, _, err := Client.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
