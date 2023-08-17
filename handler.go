package urlshort

import (
	"net/http"
	yaml "gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.

// MapHandler вернет http.HandlerFunc (который также реализует http.Handler), 
// который попытается сопоставить любые пути (ключи в карте) с соответствующими им URL (значениями, 
// на которые указывает каждый ключ в карте, в строковом формате).
// Если путь не указан в карте, то вместо него будет вызван резервный http.Handler.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		redirectUrl, ok := pathsToUrls[path]
		if ok {
			http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var paredYamls []ShortUrl
	err := yaml.Unmarshal(yml, &paredYamls)
	if err != nil {
		return nil, err
	}
	pathMap := YAMLtoMap(paredYamls)
	return MapHandler(pathMap, fallback), nil
}

func YAMLtoMap(parsedYamls []ShortUrl) map[string]string {
	mapData := make(map[string]string)
	for _, value := range parsedYamls {
		mapData[value.Path] = value.Url
	}
	return mapData
}

type ShortUrl struct {
	Path string `yaml:"path"`
	Url string	`yaml:"url"`
}