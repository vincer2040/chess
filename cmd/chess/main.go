package chess

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vincer2040/chess/internal/render"
	"github.com/vincer2040/chess/internal/routes"
)

func Main() error {
    e := echo.New()

    e.Renderer = render.New()

    e.Use(middleware.Logger())

    e.GET("/", routes.RootGet)

    e.Logger.Fatal(e.Start(":8080"))
    return nil
}
