package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func main() {
	parseDocs()

	fs := http.FileServer(http.Dir("local"))
	http.Handle("/", fs)

	log.Println("Listening on :80...")
	err := http.ListenAndServe(":80", execRequest(http.DefaultServeMux))
	if err != nil {
		log.Panic(err)
	}
}

func execRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := path.Base(r.URL.String())
		log.Printf("from: %s %s\n", r.RemoteAddr, url)
		if url == "/" {
			url = "index.html"
		}
		if strings.HasSuffix(url, ".html") || strings.HasSuffix(url, ".css") {
			parseDoc(url)
		}
		handler.ServeHTTP(w, r)
	})
}

func parseDocs() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Panic(err)
	}

	os.Mkdir("docs", os.ModePerm)
	os.Mkdir("local", os.ModePerm)
	for _, file := range files {
		name := path.Base(file.Name())
		if strings.HasSuffix(name, ".html") || strings.HasSuffix(name, ".css") {
			if !contains(name, partialTemplates) {
				parseDoc(name)
			}
		}
	}
}

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

func parseDoc(name string) {
	renderTmplToFile(name, path.Join("local", name), devTmplData)
	renderTmplToFile(name, path.Join("docs", name), prodTmplData)
}

var partialTemplates = []string{
	"template.html",
	"header.html",
	"footer.html",
	"head.html",
}

func renderTmplToFile(tmplName, fileName string, data tmplData) {
	log.Println("rendering", fileName)
	filesToParse := []string{tmplName}
	filesToParse = append(filesToParse, partialTemplates...)
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

	err = tmpl.ExecuteTemplate(f, tmplName, data)
	if err != nil {
		log.Fatal(err)
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
