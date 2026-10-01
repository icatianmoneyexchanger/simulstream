package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

type StreamTarget struct {
	ID        string `json:"id"`
	Platform  string `json:"platform"`
	StreamURL string `json:"stream_url"`
	StreamKey string `json:"stream_key"`
	Active    bool   `json:"active"`
}

type ServerState struct {
	mu      sync.RWMutex
	Targets map[string]StreamTarget
}

var state = &ServerState{
	Targets: make(map[string]StreamTarget),
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/api/targets", handleTargets)
	http.HandleFunc("/api/reload", handleReload)

	log.Printf("SimulStream Go Control Engine running on port %s", port)
	log.Printf("IPv6 / Dual-Stack HTTP API ready at http://localhost:%s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "online",
		"timestamp": time.Now().UTC(),
		"engine":    "SimulStream Go v1.0",
	})
}

func handleTargets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		state.mu.RLock()
		defer state.mu.RUnlock()
		targets := make([]StreamTarget, 0, len(state.Targets))
		for _, t := range state.Targets {
			targets = append(targets, t)
		}
		json.NewEncoder(w).Encode(targets)

	case http.MethodPost:
		var target StreamTarget
		if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if target.ID == "" {
			target.ID = fmt.Sprintf("target-%d", time.Now().UnixNano())
		}
		target.Active = true

		state.mu.Lock()
		state.Targets[target.ID] = target
		state.mu.Unlock()

		log.Printf("[SimulStream] Added new egress target: %s (%s)", target.Platform, target.ID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(target)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleReload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Hot reload NGINX container
	cmd := exec.Command("podman", "exec", "rtmp-relay", "nginx", "-s", "reload")
	output, err := cmd.CombinedOutput()

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Printf("[SimulStream] Reload failed: %v | Output: %s", err, string(output))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "error",
			"detail": string(output),
		})
		return
	}

	log.Println("[SimulStream] NGINX RTMP Relay hot-reloaded successfully.")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"detail": "RTMP Relay state updated without stream drop",
	})
}
