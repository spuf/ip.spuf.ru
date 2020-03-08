package websocket_ping

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type messageData float64

func NewWebsocketPing() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			Subprotocols: []string{"timestamp-ping"},
		})
		if err != nil {
			panic(err)
		}
		defer c.Close(websocket.StatusInternalError, "unexpected error")

		ctx, cancel := context.WithTimeout(r.Context(), time.Hour)
		defer cancel()

		m := make(chan messageData, 1)
		for {
			select {
			case <-ctx.Done():
				c.Close(websocket.StatusNormalClosure, "")
				return
			case v := <-m:
				if err := wsjson.Write(ctx, c, v); err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						return
					}
					panic(err)
				}
			default:
				var v messageData
				if err := wsjson.Read(ctx, c, &v); err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						return
					}
					if errors.Is(err, io.EOF) {
						return
					}
					panic(err)
				}
				m <- v
			}
		}
	})
}
