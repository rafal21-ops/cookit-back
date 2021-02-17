package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Daniorocket/cookit-back/handlers"
	"github.com/Daniorocket/cookit-back/routing"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World")
}
func main() {
	port := os.Getenv("PORT")
	router, err := routing.NewRouter()
	if err != nil {
		log.Print("Failed to init Router", err)
		return
	}
	router.Use(handlers.Authenticate)

	fmt.Println("Server started!")
	log.Print("Listening on :" + port)
	if err = http.ListenAndServe(":"+port, router); err != nil {
		log.Println("Failed to close server: ", err)
		return
	}

}
