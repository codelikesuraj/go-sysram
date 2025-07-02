package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/shirou/gopsutil/v4/mem"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT_GO_SSE")
	if port == "" {
		log.Fatal("PORT_GO_SSE must be set")
	}
	addr := fmt.Sprintf("localhost:%s", port)

	staticFS, _ := fs.Sub(staticFiles, "static")
	http.Handle("/", http.FileServer(http.FS(staticFS)))
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Cache-Control", "no-cache")

		for {
			m, _ := mem.VirtualMemory()
			data, _ := json.Marshal(map[string]any{
				"total":     m.Total,
				"in_use":    m.Used,
				"available": m.Available,
			})
			fmt.Fprintf(w, "data: %s\n\n", string(data))
			time.Sleep(1500 * time.Millisecond)
			w.(http.Flusher).Flush()
		}
	})

	log.Println("Listening at http://" + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
