package routes

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

func GameGet(c echo.Context) error {
    pos := ""
    websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			// Read
			msg := ""
            err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
                break
			}
            // Write
            err = websocket.Message.Send(ws, "Hello, Client!")
            if err != nil {
                c.Logger().Error(err)
                break
            }

            fmt.Println("pos", pos)
		}
	}).ServeHTTP(c.Response(), c.Request())
    return nil
}
