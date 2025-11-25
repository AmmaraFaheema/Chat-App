package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/FhmiSddq/ProyekJarkom/internal/infra/env"
	"github.com/gorilla/websocket"
)

type Registry struct {
	mu      sync.RWMutex
	conns   map[string]*websocket.Conn
	address map[string]string
}

func NewRegistry() *Registry {
	return &Registry{
		conns:   make(map[string]*websocket.Conn),
		address: make(map[string]string),
	}
}

func (r *Registry) Add(name string, conn *websocket.Conn, addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.conns[name] = conn
	r.address[name] = addr
}

func (r *Registry) RemoveByConn(conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for n, c := range r.conns {
		if c == conn {
			delete(r.conns, n)
			delete(r.address, n)
			return
		}
	}
}

func (r *Registry) GetConn(name string) *websocket.Conn {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.conns[name]
}

func (r *Registry) List() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]string)
	for n, a := range r.address {
		out[n] = a
	}
	return out
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	cfg := env.New()                        
	port := cfg.Port                         
	addr := fmt.Sprintf("0.0.0.0:%d", port) 

	reg := NewRegistry()

	staticDir := filepath.Join("internal", "app", "chat", "static")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})

	http.Handle("/style.css", http.FileServer(http.Dir(staticDir)))
	http.Handle("/app.js", http.FileServer(http.Dir(staticDir)))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(filepath.Join(staticDir, "css")))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(filepath.Join(staticDir, "assets")))))

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()
		nameBytes := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(nameBytes)
		name := string(nameBytes)
		if name == "" {
			http.Error(w, "empty name", http.StatusBadRequest)
			return
		}
		reg.mu.Lock()
		reg.address[name] = r.RemoteAddr
		reg.mu.Unlock()
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users := reg.List()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(users)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing name", http.StatusBadRequest)
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("ws upgrade:", err)
			return
		}
		defer c.Close()

		reg.Add(name, c, r.RemoteAddr)
		defer reg.RemoveByConn(c)

		log.Printf("ws connected: %s (%s)\n", name, r.RemoteAddr)

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("ws read error:", err)
				break
			}
			payload := string(msg)
			sep := -1
			for i := 0; i < len(payload); i++ {
				if payload[i] == ':' {
					sep = i
					break
				}
			}
			if sep == -1 {
				continue
			}
			to := payload[:sep]
			message := payload[sep+1:]

			target := reg.GetConn(to)
			if target != nil {
				err := target.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s:%s", name, message)))
				if err != nil {
					log.Println("failed to forward ws msg:", err)
				}
			} else {
				_ = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("server:Pengguna %s tidak tersambung", to)))
			}
		}
	})

	log.Printf("Starting Web server on %s (static dir: %s)\n", addr, staticDir)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
