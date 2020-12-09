package main

import (
	"github.com/betNevS/Go-000/Week02/controller"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/user", controller.GetUser)
	log.Println(http.ListenAndServe(":8080", nil))
}
