package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"

	"todo/model"
)

const (
	DB_NAME     = "postgres"
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_HOST     = "db"
	DB_PORT     = "5432"
)

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
	rTodo.HandleFunc("/{id:[0-9]+}", FindById).Methods(http.MethodGet)
	rTodo.HandleFunc("/", Create).Methods(http.MethodPost)
	rTodo.HandleFunc("/{id:[0-9]+}", Update).Methods(http.MethodPut)
	rTodo.HandleFunc("/{id:[0-9]+}", Delete).Methods(http.MethodDelete)

	http.ListenAndServe(":8080", r)
}

func HandleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"message": "internal error has occurred"}`))
	fmt.Println(err)
}

func Get(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, memo, is_done, due_date FROM todos")
	if err != nil {
		HandleError(w, err)
		return
	}
	defer rows.Close()
	var todos []model.Todo
	for rows.Next() {
		var t model.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Memo, &t.IsDone, &t.DueDate); err != nil {
			HandleError(w, err)
			return
		}
		todos = append(todos, t)
	}
	res, err := json.Marshal(todos)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func FindById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var t model.Todo
	err := db.QueryRow("SELECT id, title, memo, is_done, due_date FROM todos WHERE id = $1", id).Scan(&t.ID, &t.Title, &t.Memo, &t.IsDone, &t.DueDate)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("id: ", id)
		return
	}
	if err != nil {
		HandleError(w, err)
		return
	}
	res, err := json.Marshal(t)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		HandleError(w, err)
		return
	}
	var t model.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		HandleError(w, err)
		return
	}
	defer r.Body.Close()

	statement := "INSERT INTO todos (title, memo, is_done, due_date) VALUES ($1, $2, $3, $4) RETURNING id, title, memo, is_done, due_date"
	stmt, err := db.Prepare(statement)
	if err != nil {
		HandleError(w, err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(t.Title, t.Memo, t.IsDone, t.DueDate).Scan(&t.ID, &t.Title, &t.Memo, &t.IsDone, &t.DueDate)
	if err != nil {
		HandleError(w, err)
		return
	}

	res, err := json.Marshal(t)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var t model.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		HandleError(w, err)
		return
	}
	defer r.Body.Close()

	result, err := db.Exec(
		"UPDATE todos SET title = $2, memo = $3, is_done = $4, due_date = $5 WHERE id = $1",
		id,
		t.Title,
		t.Memo,
		t.IsDone,
		t.DueDate,
	)
	if err != nil {
		HandleError(w, err)
		return
	}
	count, err := result.RowsAffected()
	if count == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		HandleError(w, err)
		return
	}
	res, err := json.Marshal(t)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	row, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		HandleError(w, err)
		return
	}
	count, err := row.RowsAffected()
	if err != nil {
		HandleError(w, err)
		return
	}
	if count == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(`{"message": "delete successful"}`))
}
