package main

import (
	"article/pkg/methods"
	"article/pkg/models"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type Code struct {
	Text string `json:"text"`
}

const file = "markdown/mips-16bit.md"

func main() {

	markdown := methods.NewMarkdown()

	static := http.FileServer(http.Dir("./docs"))
	http.Handle("/articles/", http.StripPrefix("/articles/", static))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html, err := methods.LoadFile(file, markdown)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
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

		templ.Execute(w, article)
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
		file, err := methods.LoadFile(file, markdown)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(file))
	})

	fmt.Println("Server running on http://localhost:3000")

	log.Fatal(http.ListenAndServe(":3000", nil))
}
