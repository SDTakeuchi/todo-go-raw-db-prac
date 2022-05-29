package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

const (
	DB_NAME     = "postgres"
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_HOST     = "db"
	DB_PORT     = "5432"
)

type Todo struct {
	ID        uint64    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Title     string    `json:"title" db:"title"`
	Memo      string    `json:"memo" db:"memo"`
	IsDone    bool      `json:"is_done" db:"is_done"`
	DueDate   time.Time `json:"due_date" db:"due_date"`
}

var db *sql.DB

func init() {
	var err error
	config := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST,
		DB_PORT,
		DB_USER,
		DB_PASSWORD,
		DB_NAME,
	)
	db, err = sql.Open("postgres", config)
	if err != nil {
		panic(err)
	}
}

func main() {
	r := mux.NewRouter()
	rTodo := r.PathPrefix("/api/todo").Subrouter()
	rTodo.HandleFunc("/", Get).Methods(http.MethodGet)
	// rTodo.HandleFunc("/{id}", FindById).Methods(http.MethodGet)
	// rTodo.HandleFunc("/", Create).Methods(http.MethodPost)
	// rTodo.HandleFunc("/", Update).Methods(http.MethodPut)
	// rTodo.HandleFunc("/", Delete).Methods(http.MethodDelete)

	http.ListenAndServe(":8080", r)
}

func Get(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, memo, is_done, due_date FROM todo")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "internal error has occurred, db"}`))
		return
	}
	defer rows.Close()
	var todos []Todo
	for rows.Next() {
		t := Todo{}
		if err := rows.Scan(&t.ID, &t.Title, &t.Memo, &t.IsDone, &t.DueDate); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "internal error has occurred, loop"}`))
			return
		}
		todos = append(todos, t)
	}
	res, err := json.Marshal(todos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "internal error has occurred, json"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
