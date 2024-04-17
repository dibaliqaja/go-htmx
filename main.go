package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"time"
	"io"
	"log"
	"net/http"
)

type Film struct {
	Title 		string
	Director 	string
}

type Todo struct {
	UserID		int		`json:"userId"`
	ID			int		`json:"id"`
	Title		string	`json:"title"`
	Completed 	bool	`json:"completed"`
}

func main() {
	fmt.Println("Listening serve on http://localhost:8000")

	handlerHome := func (w http.ResponseWriter, r *http.Request) {
		// io.WriteString(w, "Hello Go - HTMX\n")
		// io.WriteString(w, r.Method)

		tmpl := template.Must(template.ParseFiles("index.html"))

		films := map[string][]Film {
			"Films": {
				{ Title: "The Food", Director: "Aberama Hanyu" },
				{ Title: "The Dasiid", Director: "asdad Hanyu" },
				{ Title: "The Fhasi", Director: "klo Hanyu" },
			},
		}

		tmpl.Execute(w, films)
	}

	handlerAddFilm := func (w http.ResponseWriter, r *http.Request) {
		// log.Print("HTMX request received")
		// log.Print(r.Header.Get("HX-Request"))

		time.Sleep(1 * time.Second)

		title := r.PostFormValue("title")
		director := r.PostFormValue("director")

		// fmt.Println(title)
		// fmt.Println(director)

		// htmlString := fmt.Sprintf("<li class='list-group-item bg-primary text-white m-1'>%s - %s</li>", title, director)
		// tmpl, _    := template.New("t").Parse(htmlString)
		// tmpl.Execute(w, nil)

		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "film-list-element", Film{ Title: title, Director: director })
	}

	handlerJSON := func (w http.ResponseWriter, r *http.Request) {
		url := "https://jsonplaceholder.typicode.com/todos/1"

		response, err := http.Get(url)

		if err != nil {
			log.Fatal(err)
		}

		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			// bodyBytesWeb, err := io.ReadAll(response.Body)

			// if err != nil {
			// 	log.Fatal(err)
			// }

			// data := string(bodyBytesWeb)
			// fmt.Println(data)
			// io.WriteString(w, data)

			todoItem := Todo{}
			// todoItem := &Todo{1, 1, "This Title Todo", true}
			// json.Unmarshal(bodyBytes, &todoItem)

			decoder := json.NewDecoder(response.Body)
			decoder.DisallowUnknownFields()

			if err := decoder.Decode(&todoItem); err != nil {
				log.Fatal("Decode error: ", err)
			}

			// fmt.Printf("Data from API: %+v", todoItem)

			// Convert back to JSON
			todo, err := json.MarshalIndent(todoItem, "", "\t")
			if err != nil {
				log.Fatal(err)
			}

			// fmt.Println(string(todo))
			io.WriteString(w, string(todo))
		}
	}

	http.HandleFunc("/", handlerHome)
	http.HandleFunc("/add-film/", handlerAddFilm)
	http.HandleFunc("/json", handlerJSON)

	log.Fatal(http.ListenAndServe(":8000", nil))
}