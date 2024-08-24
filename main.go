package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v4/mem"
)

//go:embed index.html
var index embed.FS

func main() {
	indexFS, _ := fs.Sub(index, ".")
	http.Handle("/", http.FileServer(http.FS(indexFS)))
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
			time.Sleep(1 * time.Second)
			w.(http.Flusher).Flush()
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
