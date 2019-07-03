package extecho

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	HealthPath = "/_health"
)

// HealthResource something that need be health checked
type HealthResource interface {
	HealthCheck() error // returns error if health check no passed
}

// NewHealth create a health check middleware
func NewHealth(rs ...HealthResource) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path == HealthPath {
				for _, r := range rs {
					if err := r.HealthCheck(); err != nil {
						return c.String(http.StatusInternalServerError, err.Error())
					}
				}
				return c.String(http.StatusOK, "OK")
			}

			return next(c)
		}
	}
}
