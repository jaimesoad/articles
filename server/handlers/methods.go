package hdl

import (
	"article/pkg/methods"
	"article/pkg/models"
	"article/pkg/tpls"
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/labstack/echo/v4"
)

var markdown = methods.NewMarkdown()

const MDDirectory = "markdown"

func GoToHomePage(c echo.Context) error {
	c.Response().WriteHeader(http.StatusNotFound)
	return c.Redirect(http.StatusPermanentRedirect, "/articles")
}

func HomePage(c echo.Context) error {
	articles := tpls.LoadArticles(methods.GetDirectoryMD(MDDirectory))

	var buf bytes.Buffer

	articles.Render(context.Background(), &buf)

	templ := template.Must(template.ParseFiles("templates/root.html"))

	return templ.Execute(c.Response(), buf.String())
}

func LoadEditor(c echo.Context) error {
	templ := template.Must(template.ParseFiles("templates/editor.html"))
	return templ.Execute(c.Response(), nil)
}

func ConvertToMD(c echo.Context) error {
	data := c.FormValue("data")

	file, err := methods.FromString(data, markdown)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.HTML(http.StatusOK, file)
}

func RawMarkdown(c echo.Context) error {
	filename := "markdown/" + c.Param("filename") + ".md"

	file, err := methods.LoadFile(filename, markdown)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.HTML(http.StatusOK, file)
}

func ListMarkdownFiles(c echo.Context) error {
	files := methods.GetDirectoryMD(MDDirectory)

	return c.JSON(http.StatusOK, models.JSON{"files": files})
}

func MarkdownByFileName(c echo.Context) error {
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
}

func LoadArticle(c echo.Context) error {
	filename := c.Param("filename")

	data, err := methods.LoadMetadata("docs/" + filename)
	if err != nil {
		data = models.Article{
			Title:    filename,
			Author:   "Jaime Acosta",
			Date:     methods.GetCurrentDate(),
			Filename: "markdown/" + filename + ".md",
		}
	}

	file, err := methods.LoadFile(data.Filename, markdown)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	data.HTML = file
	data.ReadTime = methods.GetReadTime(data.Filename)

	templ := template.Must(template.ParseFiles("templates/index.html"))

	return templ.Execute(c.Response(), data)
}
