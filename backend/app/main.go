package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type ToDo struct {
	ID       int    `json:"ID"`
	Typename string `json:"typename"`
	Todo     string `json:"todo"`
}

var db *sql.DB

func init() {
	var err error

	db, err = sql.Open("mysql", "root:admin@tcp(todo-db:3306)/todo")
	if err != nil {
		log.Fatal(err)
	}
}

func getToDos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM todo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todos := []ToDo{}
	for rows.Next() {
		var (
			id       int
			typename string
			todo     string
		)
		err := rows.Scan(&id, &typename, &todo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t := ToDo{id, typename, todo}
		todos = append(todos, t)
	}
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createToDo(w http.ResponseWriter, r *http.Request) {
	var todo ToDo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(todo.Typename)
	fmt.Println(todo.Todo)

	if _, err := db.Exec("insert into todo (typename, todo) values (?, ?)", todo.Typename, todo.Todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteToDo(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query()["id"]
	if len(ids) == 0 {
		http.Error(w, "cannot get id parameter", http.StatusInternalServerError)
		return
	}

	id := ids[0]

	if _, err := db.Exec("delete from todo where id = ?", id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE")

		switch r.Method {
		case http.MethodGet:
			getToDos(w, r)
		case http.MethodPost:
			createToDo(w, r)
		case http.MethodDelete:
			deleteToDo(w, r)
		}
	})

	defer db.Close()

	log.Println("start http server :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
