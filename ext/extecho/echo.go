package extecho

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.guoyk.net/ext/extos"
)

func NewEcho(limit int64, rs ...HealthResource) *echo.Echo {
	webroot := "public"
	extos.EnvStr(&webroot, "WEBROOT", "WEB_ROOT", "WWWROOT", "WWW_ROOT", "PUBLIC")

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Static("/", webroot)
	e.Use(NewColimit(limit))
	e.Use(NewHealth(rs...))
	e.Use(middleware.Logger())
	return e
}
