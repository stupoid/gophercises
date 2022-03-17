package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/stupoid/gophercises/urlshort"
)

func main() {
	yamlFlag := flag.String("yaml", "", "load path URls from YAML")
	jsonFlag := flag.String("json", "", "load path URls from JSON")
	boltFlag := flag.String("bolt", "", "load path URls from BoltDB")
	seedFlag := flag.Bool("seed", true, "seed data in BoltDB if bucket is not found")
	portFlag := flag.String("port", "8080", "server port")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort":       "https://github.com/gophercises/urlshort",
		"/urlshort-final": "https://github.com/gophercises/urlshort/tree/solution",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	if *yamlFlag != "" {
		yaml, err := os.ReadFile(*yamlFlag)
		if err != nil {
			log.Fatal(err)
		}
		yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Starting the server on :%s\n", *portFlag)
		http.ListenAndServe(fmt.Sprintf(":%s", *portFlag), yamlHandler)
	}

	if *jsonFlag != "" {
		data, err := os.ReadFile(*jsonFlag)
		if err != nil {
			log.Fatal(err)
		}
		jsonHandler, err := urlshort.JSONHandler(data, mapHandler)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Starting the server on :%s\n", *portFlag)
		http.ListenAndServe(fmt.Sprintf(":%s", *portFlag), jsonHandler)
	}

	if *boltFlag != "" {
		db, err := bolt.Open(*boltFlag, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("urlshort"))
			if bucket == nil {
				if !*seedFlag {
					log.Fatal("BoltDB bucket 'urlshort' not found")
				} else {
					// Seed data
					bucket, err = tx.CreateBucket([]byte("urlshort"))
					if err != nil {
						log.Fatal(err)
					}
					for path, url := range map[string]string{
						"/urlshort":       "https://github.com/gophercises/urlshort",
						"/urlshort-final": "https://github.com/gophercises/urlshort/tree/solution",
						"/bolt":           "https://github.com/boltdb/bolt",
					} {
						err := bucket.Put([]byte(path), []byte(url))
						if err != nil {
							log.Fatal(err)
						}
					}
				}
			}

			boltHandler := urlshort.BoltHandler(bucket, mapHandler)
			fmt.Printf("Starting the server on :%s\n", *portFlag)
			http.ListenAndServe(fmt.Sprintf(":%s", *portFlag), boltHandler)
			return nil
		})

	}

	fmt.Printf("Starting the server on :%s\n", *portFlag)
	http.ListenAndServe(fmt.Sprintf(":%s", *portFlag), mapHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
