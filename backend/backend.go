package backend

import (
	"bytes"
	"collaborart/backend/composedImage"
	"collaborart/backend/vcs"
	"collaborart/frontend"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"image"
	"image/jpeg"
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
		return c.Redirect(http.StatusSeeOther, "/branch")
	})

	e.GET("/branch/checkout_new_branch", func(c echo.Context) error {
		bid := c.QueryParam("bid")
		log.Println("Checkout new branch", bid)
		c.Response().Header().Add("HX-Redirect", "/beep?bid="+bid)
		return c.Redirect(http.StatusSeeOther, "/beep?bid="+bid)
	})

	e.POST("/branch/checkout_branch", func(c echo.Context) error {
		bid := c.FormValue("bid")
		fmt.Println("Checking out: ", bid)
		c.Response().Header().Add("HX-Push-Url", "/branch?bid="+bid)
		return c.Redirect(http.StatusSeeOther, "/branch?bid="+bid)
	})

	e.GET("/branch", func(c echo.Context) error {
		branchName := c.QueryParam("bid")
		log.Println("branchID", branchName)

		//var branches []string
		branches := GetBranchNames()
		log.Println("branches found", branches)
		branchDetails, branchFound := vcs.GetBranch(branchName)

		var commits []uuid.UUID
		if branchFound == nil {
			commits = branchDetails.Commits
		}

		params := map[string]interface{}{
			"BranchId": branchName,
			"Branches": branches,
			"Commits":  commits,
		}

		return c.Render(http.StatusOK, "index", params)
	})

	e.GET("/branch/change_branch_settings", func(c echo.Context) error {
		params := map[string]interface{}{
			"Branches": GetBranchNames(),
		}
		return c.Render(http.StatusOK, "change_branch_settings", params)
	})

	e.POST("/branch/upload", func(c echo.Context) error {
		branchName := c.FormValue("branchId")
		fmt.Println("Upload file to branch", branchName)

		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		img, str, err := image.Decode(src)
		if err != nil {
			return err
		}

		//fmt.Println(img, str, err)

		log.Println("Info!!!", img.Bounds(), str, err)
		PushToBranch(branchName, &img)
		c.Response().Header().Add("HX-Refresh", "true")
		return c.String(http.StatusOK, "")
	})

	e.GET("/branch/new_branch_settings", func(c echo.Context) error {
		branchName := c.QueryParam("bid")
		log.Println("Branch name:", branchName)
		params := map[string]interface{}{
			"BranchName": branchName,
		}
		return c.Render(http.StatusOK, "new_branch_settings", params)
	})

	e.GET("/branch/upload_image_settings", func(c echo.Context) error {
		branchName := c.QueryParam("bid")
		log.Println("Branch name:", branchName)
		params := map[string]interface{}{
			"BranchName": branchName,
		}
		return c.Render(http.StatusOK, "upload_image_settings", params)
	})

	e.GET("/branch/merge_branch_settings", func(c echo.Context) error {
		branches := GetBranchNames()
		params := map[string]interface{}{
			"Branches": branches,
		}
		return c.Render(http.StatusOK, "merge_branch_settings", params)
	})

	e.GET("/branch/merge_preview", func(c echo.Context) error {
		//TODO: Acc do this
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
		log.Println("Merging from:", mergingFrom, "Merging Into:", mergingTo)
		Merge(mergingFrom, mergingTo)
		c.Response().Header().Add("HX-Refresh", "true")
		return c.String(http.StatusOK, "")
	})

	e.POST("/branch/create_new_branch", func(c echo.Context) error {
		oldBranchName := c.FormValue("old_branch")
		newBranchName := c.FormValue("new_branch")
		log.Println("old branch", oldBranchName, "new branch", newBranchName)
		if newBranchName == "" {
			return c.String(http.StatusBadRequest, "")
		}

		CreateNewBranch(newBranchName, oldBranchName)
		newbr, _ := vcs.GetBranch(newBranchName)
		log.Println("Created new branch:", newBranchName, "From branch:", oldBranchName, "New Branch num commits:", len(newbr.Commits))
		c.Response().Header().Add("HX-Refresh", "true")
		return c.String(http.StatusCreated, "")
	})

	e.GET("/branch/preview", func(c echo.Context) error {
		branchId := c.QueryParam("bid")
		log.Println("Branch id", branchId)
		if branchId == "" {
			return c.String(http.StatusInternalServerError, "Missing branch id")
		}

		branchDetails, branchFound := vcs.GetBranch(branchId)

		if branchFound != nil {
			return c.String(http.StatusBadRequest, "")
		}

		target := composedImage.New(branchDetails)

		//fmt.Println("target Img", target.Img)

		//c.Response().Header().Set("Content-Type", "image/jpeg")
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, &target.Img, &jpeg.Options{Quality: 100}); err != nil {
			log.Printf("failed to encode: %v", err)
		}
		image := buf.Bytes()
		imgBase64Str := base64.StdEncoding.EncodeToString(image)
		//fmt.Println("base64 encoding", imgBase64Str)

		params := map[string]interface{}{
			"Encoding": imgBase64Str,
		}
		return c.Render(http.StatusOK, "preview", params)
	})

	e.Static("/public", "./frontend/public")

	frontend.NewTemplateRenderer(e, "./frontend/templates/*.html")

	e.Logger.Fatal(e.Start(":8000"))
}
