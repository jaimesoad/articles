package main

import (
	"article/pkg/methods"
	"article/pkg/models"
	"log"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Code struct {
	Text string `json:"text"`
}

const file = "markdown/hello.md"

func main() {

	markdown := methods.NewMarkdown()

	e := echo.New()

	e.Static("/css", "./docs/css")
	e.Static("/script", "./docs/script")

	e.GET("/", func(c echo.Context) error {
		html, err := methods.LoadFile(file, markdown)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		article := models.Article{
			Author:   "Jaime Acosta",
			Date:     methods.GetCurrentDate(),
			HTML:     html,
			ReadTime: methods.GetReadTime(file),
			Filename: file,
			Title:    "Learn Golang in one blog",
		}

		templ := template.Must(template.ParseFiles("templates/index.html"))

		return templ.Execute(c.Response(), article)
	})

	e.GET("/editor", func(c echo.Context) error {
		templ := template.Must(template.ParseFiles("templates/editor.html"))
		return templ.Execute(c.Response(), nil)
	})

	e.POST("/convert", func(c echo.Context) error {
		data := c.FormValue("data")

		file, err := methods.FromString(data, markdown)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		return c.HTML(http.StatusOK, file)
	})

	e.GET("/raw", func(c echo.Context) error {
		file, err := methods.LoadFile("hello.md", markdown)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.HTML(http.StatusOK, file)
	})

	log.Fatal(e.Start(":3000"))
}
