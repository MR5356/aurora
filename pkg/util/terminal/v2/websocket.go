package terminal

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	EventResize = "resize"
)

type Event struct {
	Operate string `json:"operate"`
	Cols    int    `json:"cols"`
	Rows    int    `json:"rows"`
}

func HandleTerminal(ws *websocket.Conn, term Terminal) {
	defer term.Close()

	go func() {
		for {
			_, data, err := ws.ReadMessage()
			if err != nil {
				return
			}
			var msg Event
			if err := json.Unmarshal(data, &msg); err == nil {
				switch msg.Operate {
				case EventResize:
					logrus.Infof("terminal resize: %d x %d", msg.Cols, msg.Rows)
					if err := term.Resize(uint32(msg.Cols), uint32(msg.Rows)); err != nil {
						logrus.Errorf("failed to resize terminal: %v", err)
						continue
					}
				default:
					logrus.Warnf("unknown operate: %s", msg.Operate)
				}
			} else {
				if _, err := term.Write(data); err != nil {
					logrus.Errorf("failed to write terminal: %v", err)
					continue
				}
			}
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := term.Read(buf)
		if err != nil {
			return
		}
		if err := ws.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
			return
		}
	}
}
