package backend

import (
	"collaborart/backend/composedImage"
	"collaborart/backend/vcs"
	"collaborart/frontend"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"image/jpeg"
	"log"
	"net/http"
	"os"
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
		// TODO: Get list of branches and commits
		branches := []string{"main"}
		commits := []string{"First Commit"}

		params := map[string]interface{}{
			"BranchId": branchId,
			"Branches": branches,
			"Commits":  commits,
		}

		return c.Render(http.StatusOK, "index", params)
	})

	e.POST("/branch/upload", func(c echo.Context) error {
		//TODO: Get branch main
		branchName := "main"

		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		//src, err := file.Open()
		//if err != nil {
		//	return err
		//}
		imgFile, err := os.Create(file.Filename)
		PushToBranch(branchName, imgFile)
		return c.String(http.StatusOK, "")
	})

	e.GET("/branch/new_branch_settings", func(c echo.Context) error {
		return c.Render(http.StatusOK, "new_branch_settings", nil)
	})

	e.GET("/branch/upload_image_settings", func(c echo.Context) error {
		return c.Render(http.StatusOK, "upload_image_settings", nil)
	})

	e.GET("/branch/merge_branch_settings", func(c echo.Context) error {
		//TODO: Get branches!!
		branches := GetBranchNames()
		params := map[string]interface{}{
			"Branches": branches,
		}
		return c.Render(http.StatusOK, "merge_branch_settings", params)
	})

	e.GET("/branch/merge_preview", func(c echo.Context) error {

		return c.String(http.StatusOK, "")
	})

	e.POST("/branch/merge_branches", func(c echo.Context) error {
		mergingFrom := c.FormValue("merging_from")
		if mergingFrom == "" {
			return c.String(http.StatusBadRequest, "")
		}
		mergingTo := c.FormValue("merging_to")
		if mergingTo == "" {
			return c.String(http.StatusBadRequest, "")
		}
		// TODO: Create Image Preview!
		log.Println("Merging from:", mergingFrom, "Merging Into:", mergingTo)
		return c.String(http.StatusOK, "")
	})

	e.POST("/branch/create_new_branch", func(c echo.Context) error {
		//TODO: Get current branch in request!
		currentBranchName := "main"
		branchName := c.FormValue("branch_name")

		if branchName == "" {
			return c.String(http.StatusBadRequest, "")
		}

		//TODO: Create new branch
		CreateNewBranch(branchName, currentBranchName)
		log.Println("Created new branch:", branchName, "From branch:", currentBranchName)
		return c.String(http.StatusCreated, "")
	})

	e.GET("/branch/preview", func(c echo.Context) error {
		branchId := c.QueryParam("branchId")
		log.Println("Branch id", branchId)
		if branchId == "" {
			return c.String(http.StatusInternalServerError, "Missing branch id")
		}

		branch := vcs.GetBranch(branchId)

		target := composedImage.New(branch)

		w := c.Response().Writer

		if err := jpeg.Encode(w, &target.Img, &jpeg.Options{Quality: 100}); err != nil {
			log.Printf("failed to encode: %v", err)
		}

		return c.String(http.StatusOK, "")
	})

	e.Static("/public", "./frontend/public")

	frontend.NewTemplateRenderer(e, "./frontend/templates/*.html")

	e.Logger.Fatal(e.Start(":8000"))
}
