package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"urlshort"
)

func main() {
	readFile := flag.String("yamlFile", "", "Set yaml file for file config")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	var yaml []byte
	if *readFile == "" {
		// Build the YAMLHandler using the mapHandler as the
		// fallback
		yamlData := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
		yaml = []byte(yamlData)
	} else {
		var err error
		yaml, err = os.ReadFile(fmt.Sprintf("../%s", *readFile))
		if err != nil {
			panic(err)
		}
	}

	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
