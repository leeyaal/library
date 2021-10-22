package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/getbook", getbook)
	http.HandleFunc("/getbooks", getbooks)
	http.HandleFunc("/createbook", createbook)
	http.HandleFunc("/refreshbook", refreshbook)
	http.HandleFunc("/deletebook", deletebook)

	fmt.Println("listen now")
	http.ListenAndServe(":7000", nil)
}
