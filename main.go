package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/code-grey/digital-notice-board/db"
    "github.com/code-grey/digital-notice-board/handlers"
    "github.com/code-grey/digital-notice-board/ws"
)

func main() {
    // Connect to MySQL and auto-migrate the Announcement model.
    db.Connect()
    log.Println("Connected to MySQL database.")

    // Start the WebSocket broadcast listener.
    go ws.StartBroadcast()

    // Create a new router.
    router := mux.NewRouter()

    // RESTful endpoints.
    router.HandleFunc("/announcements", handlers.CreateAnnouncement).Methods("POST")
    router.HandleFunc("/announcements", handlers.ListAnnouncements).Methods("GET")
    router.HandleFunc("/announcements/{id}", handlers.GetAnnouncement).Methods("GET")
    router.HandleFunc("/announcements/{id}", handlers.UpdateAnnouncement).Methods("PUT")
    router.HandleFunc("/announcements/{id}", handlers.DeleteAnnouncement).Methods("DELETE")

    // WebSocket endpoint.
    router.HandleFunc("/ws", ws.HandleWebSocket)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

    // Start the server.
    log.Println("Server starting on :8080...")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatal("ListenAndServe error:", err)
    }
}

