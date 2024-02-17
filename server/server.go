package main

import (
	hdl "article/server/handlers"
	"html/template"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	g := e.Group("/articles")

	g.Static("/static", "docs/static")

	e.RouteNotFound("*", func(c echo.Context) error {
		templ := template.Must(template.ParseFiles("templates/404.html"))

		return templ.Execute(c.Response(), nil)
	})

	e.GET("/", hdl.GoToHomePage)
	g.GET("", hdl.HomePage)
	e.GET("/editor", hdl.LoadEditor)
	e.POST("/convert", hdl.ConvertToMD)
	g.GET("/raw", hdl.ListMarkdownFiles)
	g.GET("/raw/:filename", hdl.RawMarkdown)
	g.GET("/md", hdl.ListMarkdownFiles)
	g.GET("/md/:filename", hdl.MarkdownByFileName)
	g.GET("/:filename", hdl.LoadArticle)

	log.Fatal(e.Start(":3000"))
}
