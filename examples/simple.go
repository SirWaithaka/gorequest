package main

import (
	"log"
	"net/http"

	"github.com/SirWaithaka/gorequest"
	"github.com/SirWaithaka/gorequest/corehooks"
)

type Post struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func simpleGetRequest() {
	url := "https://jsonplaceholder.typicode.com/posts/10"
	response := &Post{}
	// create an instance of request
	request := gorequest.New(gorequest.Config{Endpoint: url}, gorequest.Operation{Method: http.MethodGet}, corehooks.Default(), nil, nil, response)
	// make request
	if err := request.Send(); err != nil {
		log.Println(err)
		return
	}
	log.Println(response)

	return
}

func main() {
	simpleGetRequest()
}
