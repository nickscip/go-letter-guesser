package main

import (
        "log"
        "os"
        "net/http"

        "github.com/gorilla/websocket"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/ws", socketHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

var upgrader = websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)

        if err != nil {
                log.Println(err)
                return
        }

        for {
                messageType, p, err := conn.ReadMessage()
                if err != nil {
                        log.Println(err)
                        return
                }
                if err := conn.WriteMessage(messageType, p); err != nil {
                        log.Println(err)
                        return
                }
        }
}
