package main

import (
	handlers "bomberman/Server/Handlers"
	"bomberman/Server/websocket"
	"fmt"
	"log"
	"net/http"

	monke "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmsgprefix)
}

func main() {

	headersOk := monke.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token"})
	originsOk := monke.AllowedOrigins([]string{"http://localhost:3000"})
	methodsOk := monke.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	credentialsOK := monke.AllowCredentials()

	r := mux.NewRouter()
	r.HandleFunc("/generate", handlers.GenerateMap)
	http.Handle("/", r)

	r.PathPrefix("/Public/").Handler(http.StripPrefix("/Public/", http.FileServer(http.Dir("Public"))))

	manager := websocket.NewManager()
	r.HandleFunc("/ws", manager.ServeWS)

	fmt.Printf("Starting server at port 8080\n")
	fmt.Printf("http://localhost:8080/\n")
	log.Fatal(http.ListenAndServe(":8080", monke.CORS(originsOk, headersOk, methodsOk, credentialsOK)(r)))
}
