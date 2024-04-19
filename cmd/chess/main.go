package chess

import (
	// "fmt"
	// "github.com/vincer2040/chess/internal/game"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	"github.com/vincer2040/chess/internal/render"
	"github.com/vincer2040/chess/internal/routes"
)

func Main() error {
	// game := game.New("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	// legalMoves := game.GetLegalMoves()
	// fmt.Printf("legalMoves: %v\n", legalMoves)
	// game.PrintBoard()
	e := echo.New()

	e.Renderer = render.New()

	// e.Use(middleware.Logger())
	e.Static("pieces", "public/pieces")

	e.GET("/", routes.RootGet)
	e.GET("/game", routes.GameGet)

	e.Logger.Fatal(e.Start(":8080"))
	return nil
}
