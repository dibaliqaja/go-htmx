package main

import (
	"fmt"
	"html/template"
	"time"

	// "io"
	"log"
	"net/http"
)

type Film struct {
	Title 		string
	Director 	string
}

func main() {
	fmt.Println("Hello Wooooo")

	handlerHome := func (w http.ResponseWriter, r *http.Request) {
		// io.WriteString(w, "Hellowowwo\n")
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

		fmt.Println(title)
		fmt.Println(director)

		// htmlString := fmt.Sprintf("<li class='list-group-item bg-primary text-white m-1'>%s - %s</li>", title, director)
		// tmpl, _    := template.New("t").Parse(htmlString)
		// tmpl.Execute(w, nil)

		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "film-list-element", Film{ Title: title, Director: director })
	}

	http.HandleFunc("/", handlerHome)
	http.HandleFunc("/add-film/", handlerAddFilm)

	log.Fatal(http.ListenAndServe(":8000", nil))
}