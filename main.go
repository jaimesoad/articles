package main

import (
	"article/pkg/methods"
	"article/pkg/models"
	"bytes"
	"flag"
	"fmt"
	"os"
	"text/template"
)

func main() {
	var article models.Article

	flag.StringVar(&article.Author, "author", "Jaime Acosta", "")
	flag.StringVar(&article.Title, "title", "", "")
	flag.StringVar(&article.Filename, "file", "", "")
	flag.Parse()

	markdown := methods.NewMarkdown()

	converted, err := methods.LoadFile(article.Filename, markdown)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	article.HTML = converted
	article.Date = methods.GetCurrentDate()
	article.ReadTime = methods.GetReadTime(article.Filename)

	var buf bytes.Buffer

	templ := template.Must(template.ParseFiles("templates/index.html"))

	templ.Execute(&buf, article)

	outdir := "articles/" + article.Title

	_, err = os.Stat(outdir)

	if err != nil {
		err = os.Mkdir(outdir, os.ModePerm)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	f, err := os.Create(outdir + "/index.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	f.WriteString(buf.String())
}
