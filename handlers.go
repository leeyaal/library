package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

//функция для получения одной конкретной книги...
func getbook(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("need post"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var br BookRequest
	err := decoder.Decode(&br)
	if err != nil {
		fmt.Fprint(w, "request can't be decoded: ", err)
		return
	}

  //подключение к БД...
	connStr := "user=postgres password=qwerty7 dbname=library sslmode=disable port=8080"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("1", err)
	}
	defer db.Close()

	booksid := []BookRequest{}
	booksath := []BookRequest{}
	booksnm := []BookRequest{}

	// поиск книги по id если в json пришел id
	if br.Author == "" && br.Name == "" {
		rows1, err := db.Query("select * from books where id = $1", br.Id)
		if err != nil {
			fmt.Fprint(w, "unexpected error")
		}
		if err == sql.ErrNoRows {
			fmt.Fprint(w, "gg")
		}

		defer rows1.Close()

		for rows1.Next() {
			err := rows1.Scan(&br.Id, &br.Author, &br.Name)
			if err != nil {
				fmt.Fprint(w, "smth happens")
				continue
			}
			booksid = append(booksid, br)
		}
		for _, br := range booksid {
			sntnc, _ := fmt.Printf("%v %s : %s \n", br.Id, br.Author, br.Name)
			fmt.Fprint(w, sntnc)
		}
	}
	// поик книги по author если в json пришел автор 
	if br.Id == 0 && br.Name == "" {
		rows2, err := db.Query("select * from books where author = $1", br.Author)
		if err != nil {
			panic(err)
		}
		defer rows2.Close()

		for rows2.Next() {
			err := rows2.Scan(&br.Id, &br.Author, &br.Name)
			if err != nil {
				fmt.Println(err)
				continue
			}
			booksath = append(booksath, br)
		}
		for _, br := range booksath {
			sntnc, _ := fmt.Printf("%v %s : %s \n", br.Id, br.Author, br.Name)
			fmt.Fprint(w, sntnc)
		}
	}
	//поиск книги по name если в json пришло название книги
	if br.Id == 0 && br.Author == "" {
		rows3, err := db.Query("select * from books where name = $1", br.Name)
		if err != nil {
			panic(err)
		}
		defer rows3.Close()

		for rows3.Next() {
			err := rows3.Scan(&br.Id, &br.Author, &br.Name)
			if err != nil {
				fmt.Println(err)
				continue
			}
			booksnm = append(booksnm, br)
		}
		for _, br := range booksnm {
			sntnc, _ := fmt.Printf("%v %s : %s \n", br.Id, br.Author, br.Name)
			fmt.Fprint(w, sntnc)
		}
	}
}

//функция для получения всего списка книг
func getbooks(w http.ResponseWriter, r *http.Request) {

	br := BookRequest{}
	books := []BookRequest{}
	connStr := "user=postgres password=qwerty7 dbname=library sslmode=disable port=8080"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("1", err)
	}
	defer db.Close()

	rows, err := db.Query("select * from books")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&br.Id, &br.Author, &br.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		books = append(books, br)
	}
	for _, br := range books {
		sntnc, _ := fmt.Printf("%v %s : %s \n", br.Id, br.Author, br.Name)
		fmt.Fprint(w, sntnc)
	}
}

//функция для добавления книги в БД
func createbook(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("need post"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var br BookRequest
	err := decoder.Decode(&br)
	if err != nil {
		fmt.Fprint(w, "request can't be decoded: ", err)
		return
	}

	connStr := "user=postgres password=qwerty7 dbname=library sslmode=disable port=8080"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("1", err)
	}
	defer db.Close()

	row1 := db.QueryRow("select * from books where author = $1 and name = $2", br.Author, br.Name)
	err = row1.Scan(&br.Author, &br.Name)
	if err != sql.ErrNoRows {
		fmt.Fprint(w, "book is already exists, try another one")
	} else {

		result1, err := db.Exec("insert into books (author, name) values ($1, $2)", br.Author, br.Name)
		if err != nil {
			fmt.Fprint(w, "unexpected error")
		} else {
			fmt.Println(result1)
		}
	}

}

//функция для изменения данных по конкретной книге (проводится по id)
func refreshbook(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("need post"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var br BookRequest
	err := decoder.Decode(&br)
	if err != nil {
		fmt.Fprint(w, "request can't be decoded: ", err)
		return
	}

	connStr := "user=postgres password=qwerty7 dbname=library sslmode=disable port=8080"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("1", err)
	}
	defer db.Close()

	row1 := db.QueryRow("select * from books where id = $1", br.Id)
	err = row1.Scan(&br.Id)
	if err == sql.ErrNoRows {
		errors.New("book does not exist")
	}

	result1, err := db.Exec("update books set author = $1, name = $2 where id = $3", br.Author, br.Name, br.Id)
	if err != nil {
		fmt.Fprint(w, "unexpected error")
	} else {
		fmt.Println(result1)
	}
}

//функция для удаления книги из БД
func deletebook(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("need post"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var br BookRequest
	err := decoder.Decode(&br)
	if err != nil {
		fmt.Fprint(w, "request can't be decoded: ", err)
		return
	}

	connStr := "user=postgres password=qwerty7 dbname=library sslmode=disable port=8080"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("1", err)
	}
	defer db.Close()

	result1, err := db.Exec("delete from books where id = $1", br.Id)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(result1.RowsAffected())
	}

}
