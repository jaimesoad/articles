package main

import (
	"article/pkg/methods"
	"article/pkg/models"
	"article/pkg/tpls"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

type Code struct {
	Text string `json:"text"`
}

const MDDirectory = "markdown"

func main() {

	markdown := methods.NewMarkdown()

	static := http.FileServer(http.Dir("./docs"))
	http.Handle("/articles/", http.StripPrefix("/articles/", static))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		articles := tpls.LoadArticles(methods.GetDirectoryMD(MDDirectory))

		var buf bytes.Buffer

		articles.Render(context.Background(), &buf)

		templ := template.Must(template.ParseFiles("templates/root.html"))

		templ.Execute(w, buf.String())
	})

	http.HandleFunc("/editor", func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFiles("templates/editor.html"))
		templ.Execute(w, nil)
	})

	http.HandleFunc("POST /convert", func(w http.ResponseWriter, r *http.Request) {
		data := r.FormValue("data")

		file, err := methods.FromString(data, markdown)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(file))
	})

	http.HandleFunc("/raw", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/md", http.StatusPermanentRedirect)
	})

	http.HandleFunc("/raw/{filename}", func(w http.ResponseWriter, r *http.Request) {
		filename := "markdown/" + r.PathValue("filename") + ".md"

		file, err := methods.LoadFile(filename, markdown)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(file))
	})

	http.HandleFunc("/md", func(w http.ResponseWriter, r *http.Request) {
		files := methods.GetDirectoryMD(MDDirectory)

		blob, err := json.Marshal(models.JSON{"files": files})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(blob)
	})

	http.HandleFunc("/md/{filename}", func(w http.ResponseWriter, r *http.Request) {
		filename := "markdown/" + r.PathValue("filename") + ".md"

		file, err := os.Open(filename)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		blob, err := io.ReadAll(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(blob)
	})

	http.HandleFunc("/{filename}", func(w http.ResponseWriter, r *http.Request) {
		filename := r.PathValue("filename")

		data := models.Article{
			Title:    filename,
			Author:   "Jaime Acosta",
			Date:     methods.GetCurrentDate(),
			Filename: "markdown/" + filename,
		}

		file, err := methods.LoadFile(data.Filename+".md", markdown)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		data.HTML = file
		data.ReadTime = methods.GetReadTime(data.Filename + ".md")

		templ := template.Must(template.ParseFiles("templates/index.html"))

		templ.Execute(w, data)
	})

	fmt.Println("Server running on http://localhost:3000")

	log.Fatal(http.ListenAndServe(":3000", nil))
}
