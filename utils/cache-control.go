package utils

import (
	"github.com/labstack/echo/v4"
)

// noCacheMiddleware sets Cache-Control to no-store in development mode.
func NoCacheMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// You could make this conditional on an environment variable if needed.
		c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		return next(c)
	}
}
