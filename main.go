package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

type Option struct {
	Text string
	Arc string
}

type Story struct {
	Title string
	Story []string
	Options []Option
}

const(
	tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<h2>{{.Title}}</h2>
		{{range .Story}}<div>{{ . }}</div>{{end}}

		{{range .Options}}<p><a href="http://localhost:8080/{{.Arc}}">{{.Text}}</a></p>{{end}}
	</body>
</html>`
)

func main() {
	terminal := flag.Bool("terminal", false, "Launch in terminal")
	flag.Parse()

	jsonData, err := os.ReadFile("gopher.json")
	if err != nil {
		panic(err)
	}
	stories := make(map[string]Story)
	err = json.Unmarshal(jsonData, &stories)
	if err != nil {
		panic(err)
	}

	if *terminal {
		story := stories["intro"]
		for {
			fmt.Println(story.Title + "\n")
			for _, text := range story.Story {
				fmt.Print(text + " ")
			}
			fmt.Println("\n")
			for idx, val := range story.Options {
				fmt.Printf("%d. %s\n", idx + 1, val.Text)
			}
			fmt.Print("\nEnter: \n")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			num, _ := strconv.Atoi(input[:1])
			fmt.Println(num)
			arc := story.Options[num - 1].Arc
			story = stories[arc]
		}
	} else {
		mux := http.NewServeMux()
		mux.HandleFunc("/", MapHandler(stories))
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", mux)
	}
}

func MapHandler(stories map[string]Story) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.New("webpage").Parse(tpl)
		if err != nil {
			panic(err)
		}
		path := r.URL.Path
		story, ok := stories[path[1:]]
		if !ok {
			story = stories["intro"]
		}
		temp.Execute(w, story)
	})
}
