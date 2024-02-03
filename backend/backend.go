package backend

import (
	"collaborart/frontend"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"net/http"
)

func StartServer() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(20),
	)))

	e.GET("/", func(c echo.Context) error {
		params := map[string]interface{}{
			"Name": "Claire",
		}
		return c.Render(http.StatusOK, "index", params)
	})

	e.POST("/branch/push", func(c echo.Context) error {
		return c.String(http.StatusOK, "branch-push")
	})

	e.POST("/checkout/commit", func(c echo.Context) error {
		return c.String(http.StatusOK, "checkout-commit")
	})

	e.POST("/branch/new", func(c echo.Context) error {
		return c.String(http.StatusOK, "branch-new")
	})

	e.POST("/merge", func(c echo.Context) error {
		return c.String(http.StatusOK, "merge")
	})

	e.Static("/public", "./frontend/public")

	frontend.NewTemplateRenderer(e, "./frontend/templates/*.html")

	e.Logger.Fatal(e.Start(":8000"))
}
