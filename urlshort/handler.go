package urlshort

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

type pathUrl struct {
	Path string
	URL  string
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.EscapedPath()
		if url, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
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
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yml []byte) ([]pathUrl, error) {
	pathUrls := []pathUrl{}
	err := yaml.Unmarshal([]byte(yml), &pathUrls)
	if err != nil {
		return nil, errors.New("could not parse yaml")
	}
	return pathUrls, nil
}

func buildMap(pathUrls []pathUrl) map[string]string {
	pathsToUrls := map[string]string{}
	for _, pathUrl := range pathUrls {
		pathsToUrls[pathUrl.Path] = pathUrl.URL
	}
	return pathsToUrls
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
// 		[
//			{
// 				"path": "/some-path",
//      		"url": "https://www.some-url.com/demo"
// 			}
// 		]
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(data)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJSON)
	return MapHandler(pathMap, fallback), nil
}

func parseJSON(byt []byte) ([]pathUrl, error) {
	pathUrls := []pathUrl{}
	err := json.Unmarshal(byt, &pathUrls)
	if err != nil {
		return nil, errors.New("could not parse json")
	}
	return pathUrls, nil
}

// BoltHandler will load the bucket urlshort and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not found in the bucket, then the
// fallback http.Handler will be called instead.
func BoltHandler(b *bolt.Bucket, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.EscapedPath()
		url := b.Get([]byte(path))
		if url != nil {
			http.Redirect(w, r, string(url), http.StatusTemporaryRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})
}
