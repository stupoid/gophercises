package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

type StoryArc struct {
	Title   string
	Story   []string
	Options []struct {
		Text string
		Arc  string
	}
}

// Starts the prompt for the story arc and returns
// the arc name for the next option chosen
// if there are no options fo the story arc, "" is returned
func (s StoryArc) Prompt() string {
	fmt.Println(s.Title)
	for _, line := range s.Story {
		fmt.Println(line)
	}
	if len(s.Options) > 0 {
		var items []string
		for _, option := range s.Options {
			items = append(items, option.Text)
		}

		prompt := promptui.Select{
			Label: "Select Option",
			Items: items,
		}

		i, _, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}
		return s.Options[i].Arc
	}
	return ""
}

func main() {
	cliFlag := flag.Bool("cli", false, "cli mode")
	portFlag := flag.String("port", "8080", "port for web server")
	flag.Parse()

	data, err := os.ReadFile("gopher.json")
	if err != nil {
		log.Fatal(err)
	}
	storyArcs := map[string]StoryArc{}
	err = json.Unmarshal(data, &storyArcs)
	if err != nil {
		log.Fatal(err)
	}

	arc := "intro" // initial arc

	if *cliFlag {
		fmt.Println("Starting cli mode")
		for arc != "" {
			arc = storyArcs[arc].Prompt()
		}
		fmt.Println("The End.")
	} else {
		tmpl := template.Must(template.ParseFiles("templates/main.tmpl"))
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimLeft(r.URL.EscapedPath(), "/")
			if storyArc, ok := storyArcs[path]; ok {
				tmpl.Execute(w, storyArc)
			} else {
				http.Redirect(w, r, fmt.Sprintf("/%s", arc), http.StatusPermanentRedirect)
			}
		})
		fmt.Printf("Starting the server on :%s\n", *portFlag)
		http.ListenAndServe(fmt.Sprintf(":%s", *portFlag), nil)
	}
}
