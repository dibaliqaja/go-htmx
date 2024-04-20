package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
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

type Product struct {
	Name		string
	Price		float64
	Available	bool
}

type Todois struct {
	Id			int		`json:"id"`
	Name		string	`json:"name"`
	IsCompleted bool	`json:"isCompleted"`
}

var todos = []Todois{
	{ Id: 1, Name: "Adding service list all products", IsCompleted: true },
	{ Id: 2, Name: "Adding service create product", IsCompleted: false },
	{ Id: 3, Name: "Adding service update product", IsCompleted: false },
}

var templates map[string]*template.Template
var lastID int64 = 3

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates["todois.html"] = template.Must(template.ParseFiles("todois.html"))
	templates["todois_list.html"] = template.Must(template.ParseFiles("todois_list.html"))
}

func todoisHandler(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(todos)

	if err != nil {
		log.Fatal(err)
	}

	tmpl := templates["todois.html"]
	tmpl.ExecuteTemplate(w, "todois.html", map[string]template.JS{"Todois": template.JS(json)})
}

func createTodoisHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

	name := r.PostFormValue("name")
	completed := r.PostFormValue("completed") == "true"
	newID := int(atomic.AddInt64(&lastID, 1))

	todo := Todois{ Id: newID, Name: name, IsCompleted: completed }

	todos = append(todos, todo)

	tmpl := templates["todois_list.html"]
	tmpl.ExecuteTemplate(w, "todois_list.html", todo)
}

func deleteTodoisHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

	id, err := strconv.Atoi(r.PostFormValue("id"))
	if err != nil {
		log.Fatal(err)
	}

	for i, todo := range todos {
		if todo.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			break
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", todoisHandler)
	http.HandleFunc("/create-todo/", createTodoisHandler)
	http.HandleFunc("/delete-todo/", deleteTodoisHandler)

	fmt.Println("Server is running at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func connectDB() {
	connStr := "postgres://postgres:p@localhost:5432/gopgtest?sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}
	
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	
	defer db.Close()

	// createProductTable(db)
	// getProductWithInsert(db)
	// getAllProduct(db)
}

func createProductTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS product (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price NUMERIC(6,2) NOT NULL,
		available BOOLEAN,
		created_at timestamp DEFAULT NOW()
	)`

	_, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
	}
}

func insertProduct(db *sql.DB, product Product) int {
	query := `INSERT INTO product (name, price, available) VALUES ($1, $2, $3) RETURNING id`

	var primaryKey int

	err := db.QueryRow(query, product.Name, product.Price, product.Available).Scan(&primaryKey)

	if err != nil {
		log.Fatal(err)
	}

	return primaryKey
}

func getProductWithInsert(db *sql.DB) {
	product := Product{"Book", 15.55, true}
	primaryKey := insertProduct(db, product)

	fmt.Printf("ID = %d\n", primaryKey)

	var name string
	var price float64
	var available bool

	query := "SELECT name, price, available FROM product WHERE id = $1"

	err := db.QueryRow(query, primaryKey).Scan(&name, &price, &available)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("No rows found with ID %d", 111)
		}
		log.Fatal(err)
	}

	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Price: %f\n", price)
	fmt.Printf("Available: %t\n", available)   // %t is prints true or false
}

func getAllProduct(db *sql.DB) {
	data := []Product{}
	rows, err := db.Query("SELECT name, price, available FROM product")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var name string
	var price float64
	var available bool

	for rows.Next() {
		err := rows.Scan(&name, &price, &available)

		if err != nil {
			log.Fatal(err)
		}

		data = append(data, Product{name, price, available})
	}

	fmt.Println(data)
}

func routeHandler() {	
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