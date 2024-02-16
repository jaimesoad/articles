package main

import (
	"article/pkg/methods"
	"article/pkg/models"
	"article/pkg/tpls"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const author = "Jaime Acosta"
const dirname = "markdown"

func main() {
	articles := methods.GetDirectoryMD(dirname)

	templ := template.Must(template.ParseFiles("templates/root.html"))

	links := tpls.LoadArticles(articles)

	root, err := os.Create("docs/index.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	var linkBlob bytes.Buffer

	links.Render(context.Background(), &linkBlob)

	templ.Execute(root, linkBlob.String())
	root.Close()

	markdown := methods.NewMarkdown()

	for _, filename := range articles {
		outdir := "docs/" + strings.TrimSuffix(filename, filepath.Ext(filename))

		article, err := methods.LoadMetadata(outdir)
		if err != nil {
			article = models.Article{
				Author:   author,
				Filename: dirname + "/" + filename,
				Title:    strings.TrimSuffix(filename, filepath.Ext(filename)),
				Date:     methods.GetCurrentDate(),
			}
		}

		converted, err := methods.LoadFile(article.Filename, markdown)
		if err != nil {
			log.Panic(err)
			return
		}

		article.HTML = converted
		article.ReadTime = methods.GetReadTime(article.Filename)

		templ = template.Must(template.ParseFiles("templates/index.html"))

		_, err = os.Stat(outdir)

		if err != nil {
			err = os.Mkdir(outdir, os.ModePerm)

			if err != nil {
				log.Panic(err)
				return
			}
		}

		f, err := os.Create(outdir + "/index.html")
		if err != nil {
			log.Panic(err)
			return
		}

		templ.Execute(f, article)
		f.Close()

		f, err = os.Create(outdir + "/metadata.json")
		if err != nil {
			log.Panic(err)
			return
		}

		blob, err := json.MarshalIndent(article, "", "\t")
		if err != nil {
			log.Panic(err)
			return
		}

		f.Write(blob)
		f.Close()
	}
}
