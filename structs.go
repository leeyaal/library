package main

type (
	BookRequest struct {
		Author string `json: "author"`
		Name   string `json: "name"`
		Id     int    `json: "id"`
	}
)
