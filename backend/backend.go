package backend

import (
	"collaborart/frontend"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"log"
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
		return c.Redirect(http.StatusSeeOther, "/branch?bid=main")
	})

	e.GET("/branch", func(c echo.Context) error {
		branchId := c.QueryParam("branchId")
		log.Println("branchID", branchId)
		if branchId == "" {
			branchId = "main"
		}
		// TODO: Get list of branches
		branches := []string{"main"}

		params := map[string]interface{}{
			"BranchId": branchId,
			"Branches": branches,
		}

		return c.Render(http.StatusOK, "index", params)
	})

	e.GET("/branch/new_branch_settings", func(c echo.Context) error {
		return c.Render(http.StatusOK, "new_branch_settings", nil)
	})

	e.POST("/branch/create_new_branch", func(c echo.Context) error {
		branchName := c.FormValue("branch_name")

		if branchName == "" {
			return c.String(http.StatusBadRequest, "")
		}

		//TODO: Create new branch
		log.Println("Create new branch!", branchName)
		return c.String(http.StatusCreated, "")
	})

	e.GET("/branch/preview", func(c echo.Context) error {
		branchId := c.QueryParam("branchId")

		if branchId == "" {
			return c.String(http.StatusInternalServerError, "Missing branch id")
		}

		// TODO: Get the image to the client somehow - can be either as a byte array or straight from a file
		return c.String(http.StatusNotImplemented, "Not implemented yet : )")
	})

	e.POST("/branch/push", func(c echo.Context) error {
		c.QueryParam("ID")
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
