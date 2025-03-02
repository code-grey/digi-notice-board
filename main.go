package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "digi-notice-board/handlers"
    "digi-notice-board/ws"
    "digi-notice-board/db"
)

func main() {

    // Connect to DB
    db.Connect()

    // Start the WebSocket broadcast in a goroutine
    go ws.StartBroadcast()

    // Set up the router
    router := mux.NewRouter()

    // REST endpoints
    router.HandleFunc("/announcements", handlers.CreateAnnouncement).Methods("POST")
    router.HandleFunc("/announcements", handlers.ListAnnouncements).Methods("GET")
    router.HandleFunc("/announcements/{id}", handlers.GetAnnouncement).Methods("GET")
    router.HandleFunc("/announcements/{id}", handlers.UpdateAnnouncement).Methods("PUT")
    router.HandleFunc("/announcements/{id}", handlers.DeleteAnnouncement).Methods("DELETE")

    // WebSocket endpoint
    router.HandleFunc("/ws", ws.HandleWebsocket)
    
    //auth endpoints
    
    router.HandleFunc("/register", handlers.Register).Methods("POST")
    router.HandleFunc("/login", handlers.Login).Methods("POST")
    
	router.PathPrefix("/static/").Handler(
	    http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
	    )
	log.Println("Listening on :8080...")
    log.Fatal(http.ListenAndServe(":8080", router))



    // Start server
    log.Println("Listening on :8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}

