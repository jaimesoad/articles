package methods

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func openFile(name string) ([]byte, error) {
	file, err := os.Open(name)
	if err != nil {
		return []byte{}, err
	}

	return io.ReadAll(file)
}

func LoadFile(name string, markdown goldmark.Markdown) (string, error) {
	file, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	source, err := openFile(name)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	err = markdown.Convert(source, &buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func FromString(data string, markdown goldmark.Markdown) (string, error) {
	var buf bytes.Buffer

	err := markdown.Convert([]byte(data), &buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func NewMarkdown() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Footnote,
			extension.Typographer,
			extension.Strikethrough,
			extension.GFM,
			extension.DefinitionList,
			extension.TaskList,
			extension.Linkify,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(html.WithXHTML()),
	)
}

func GetCurrentDate() string {
	t := time.Now()

	return fmt.Sprintf("%s %d, %d", t.Month().String()[0:3], t.Day(), t.Year())
}

func GetReadTime(name string) int {
	article, err := openFile(name)
	if err != nil {
		return 0
	}

	post := string(article)

	WORDS_PER_MINUTE := 200

	regex := regexp.MustCompile(`\w+`)
	wordCount := len(regex.FindAllString(post, -1))

	return wordCount / WORDS_PER_MINUTE
}
