package main

import (
	"article/pkg/methods"
	"article/pkg/models"
	"bytes"
	"fmt"
	"os"
	"path"
	"text/template"
)

var files = GetDirectoryMD("markdown")

var data = models.Article{
	Title:    "Title",
	Author:   "Jaime Acosta",
	Date:     methods.GetCurrentDate(),
	Filename: "markdown/" + files[0],
}

func main() {
	markdown := methods.NewMarkdown()

	file, err := methods.LoadFile(data.Filename, markdown)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data.HTML = file
	data.ReadTime = methods.GetReadTime(data.Filename)

	templ := template.Must(template.ParseFiles("templates/index.html"))

	var buf bytes.Buffer

	templ.Execute(&buf, data)

	fmt.Println(buf.String())
}

func GetDirectoryMD(directory string) []string {
	var out []string

	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println(err.Error())
		return out
	}

	for _, file := range files {

		if file.IsDir() || path.Ext(file.Name()) != ".md" {
			continue
		}

		out = append(out, file.Name())
	}

	return out
}
