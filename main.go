package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type tmplData struct {
	BaseURL string
	VueSRC  string
}

var prodTmplData = tmplData{
	BaseURL: "https://j4rv.github.io/genshin-calc",
	VueSRC:  "https://unpkg.com/vue@3.0.4/dist/vue.global.prod.js",
}

var devTmplData = tmplData{
	BaseURL: "http://localhost",
	VueSRC:  "https://unpkg.com/vue@3.0.4/dist/vue.global.js",
}

func main() {
	parseDocs()

	fs := http.FileServer(http.Dir("local"))
	http.Handle("/", fs)

	log.Println("Listening on :80...")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}

var templates = []string{"template.html", "header.html", "footer.html", "head.html"}

func parseDocs() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Panic(err)
	}

	os.Mkdir("docs", os.ModePerm)
	os.Mkdir("local", os.ModePerm)

	for _, file := range files {
		name := file.Name()
		switch {
		case contains(name, templates):
			continue
		case strings.HasSuffix(name, ".html"), strings.HasSuffix(name, ".css"):
			renderTmplToFile(name, "./local/"+name, devTmplData)
			renderTmplToFile(name, "./docs/"+name, prodTmplData)
		}
	}
}

func contains(s string, ss []string) bool {
	for i := range ss {
		if ss[i] == s {
			return true
		}
	}
	return false
}

func renderTmplToFile(tmplName, fileName string, data tmplData) {
	filesToParse := []string{tmplName}
	filesToParse = append(filesToParse, templates...)
	log.Println("parsing", filesToParse)
	withCustomDelims := template.New("").Delims("[[", "]]")
	tmpl, err := withCustomDelims.ParseFiles(filesToParse...)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.Println("rendering", fileName)
	err = tmpl.ExecuteTemplate(f, tmplName, data)
	if err != nil {
		log.Fatal(err)
	}
}
