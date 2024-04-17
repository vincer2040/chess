package routes

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/vincer2040/chess/internal/game"
	"github.com/vincer2040/chess/internal/protocol"
	"github.com/vincer2040/chess/internal/types"
)

var (
	upgrader = websocket.Upgrader{}
)

func GameGet(c echo.Context) error {
    g := game.New("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
            if err.Error() == "websocket: close 1001 (going away)" {
                fmt.Println("closing")
                break
            }
			c.Logger().Error(err)
            break
		}
        parser := protocol.NewParser(msg)
        data := parser.Parse()
        fmt.Printf("received: %+v\n", data)

        buf := handleData(&data, &g)

        fmt.Println("writing:", string(buf))
        err = ws.WriteMessage(websocket.TextMessage, buf)
        if err != nil {
            c.Logger().Error(err)
            break
        }
	}
    return nil
}

func handleData(data *types.Data, game *game.Game) protocol.Builder {
    b := protocol.NewBuilder()
    switch data.Type {
    case types.IllegalType:
        b = b.AddError("invalid message")
        break
    case types.CommandType:
        cmd := data.Data.(types.Command)
        fmt.Println("command:", cmd)
        b = b.AddCommand("OK")
        break
    case types.MoveType:
        move := data.Data.(types.Move)
        fmt.Printf("move: %+v\n", move)
        b = b.AddCommand("OK")
        break
    case types.PositionType:
        pos := data.Data.(types.Position)
        fmt.Println("position:", pos)
        b = b.AddCommand("OK")
        break
    }
    return b
}
