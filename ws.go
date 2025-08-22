package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/coder/websocket"
	"github.com/siiimooon/go-polar/pkg/h10"
)

type WSServer struct {
	HR chan h10.HeartRateMeasurement
}

func (wss *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)

	if err != nil {
		return
	}

	defer c.CloseNow()

	b10 := make([]byte, 0, 8)
	ctx := c.CloseRead(context.Background())

	for {
		select {
		case msg, ok := <-wss.HR:
			if !ok {
				return
			}

			b10 = b10[:0]
			b10 = strconv.AppendInt(b10, int64(msg.GetHeartRate()), 10)
			err = c.Write(ctx, websocket.MessageText, b10)

			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
