package main

import (
	"article/pkg/methods"
	"article/pkg/models"
	"article/pkg/tpls"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Code struct {
	Text string `json:"text"`
}

const MDDirectory = "markdown"

func main() {

	markdown := methods.NewMarkdown()

	e := echo.New()

	g := e.Group("/articles")

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/articles")
	})

	g.Static("/static", "docs/static")

	g.GET("", func(c echo.Context) error {
		articles := tpls.LoadArticles(methods.GetDirectoryMD(MDDirectory))

		var buf bytes.Buffer

		articles.Render(context.Background(), &buf)

		templ := template.Must(template.ParseFiles("templates/root.html"))

		return templ.Execute(c.Response(), buf.String())
	})

	g.GET("/editor", func(c echo.Context) error {
		templ := template.Must(template.ParseFiles("templates/editor.html"))
		return templ.Execute(c.Response(), nil)
	})

	g.POST("/convert", func(c echo.Context) error {
		data := c.FormValue("data")

		file, err := methods.FromString(data, markdown)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		return c.HTML(http.StatusOK, file)
	})

	g.GET("/raw", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/md")
	})

	g.GET("/raw/:filename", func(c echo.Context) error {
		filename := "markdown/" + c.Param("filename") + ".md"

		file, err := methods.LoadFile(filename, markdown)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}

		return c.HTML(http.StatusOK, file)
	})

	g.GET("/md", func(c echo.Context) error {
		files := methods.GetDirectoryMD(MDDirectory)

		return c.JSON(http.StatusOK, models.JSON{"files": files})
	})

	g.GET("/md/:filename", func(c echo.Context) error {
		filename := "markdown/" + c.Param("filename") + ".md"

		file, err := os.Open(filename)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}

		blob, err := io.ReadAll(file)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, string(blob))
	})

	g.GET("/:filename", func(c echo.Context) error {
		filename := c.Param("filename")

		data := models.Article{
			Title:    filename,
			Author:   "Jaime Acosta",
			Date:     methods.GetCurrentDate(),
			Filename: "markdown/" + filename,
		}

		file, err := methods.LoadFile(data.Filename+".md", markdown)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}

		data.HTML = file
		data.ReadTime = methods.GetReadTime(data.Filename + ".md")

		templ := template.Must(template.ParseFiles("templates/index.html"))

		return templ.Execute(c.Response(), data)
	})

	fmt.Println("Server running on http://localhost:3000")

	log.Fatal(e.Start(":3000"))
}
