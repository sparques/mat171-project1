package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func main() {
	listenPort := flag.Int("server", 0, "start a server")
	flag.Parse()

	if *listenPort != 0 {
		http.HandleFunc("/", servePage)
		http.ListenAndServe(fmt.Sprintf(":%d", *listenPort), nil)
	}

	os.Stdout.Write(renderFile(flag.Arg(0)))
}

func servePage(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Clean(strings.Trim(r.URL.Path, "/"))
	log.Println("request for", filePath)
	if filePath != "." {
		fh, err := os.Open(filePath)
		if err != nil {
			// set 404
			log.Println("error opening", filePath, ":", err)
			return
		}
		defer fh.Close()
		switch filepath.Ext(filePath) {
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		}
		_, err = io.Copy(w, fh)
		if err != nil {
			log.Println("err copying to socket:", err)
		}
		return
	}
	w.Write(renderFile(flag.Arg(0)))
}

func renderFile(fname string) []byte {
	fh, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	content, err := io.ReadAll(fh)
	if err != nil {
		panic(err)
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Footnotes | parser.Attributes | parser.Mmark | parser.Includes | parser.MathJax | parser.Tables
	pars := parser.NewWithExtensions(extensions)
	opts := html.RendererOptions{Flags: html.CommonFlags | html.CompletePage}
	renderer := html.NewRenderer(opts)

	astDoc := markdown.Parse(content, pars)

	return markdown.Render(astDoc, renderer)
}
