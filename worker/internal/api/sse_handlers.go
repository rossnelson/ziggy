package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"ziggy/internal/workflow"
)

type SSEEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var lastStateJSON string
	var lastChatJSON string

	// Send initial state immediately
	s.sendStateUpdate(w, flusher, &lastStateJSON)
	s.sendChatUpdate(w, flusher, &lastChatJSON)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.sendStateUpdate(w, flusher, &lastStateJSON)
			s.sendChatUpdate(w, flusher, &lastChatJSON)
		}
	}
}

func (s *Server) sendStateUpdate(w http.ResponseWriter, flusher http.Flusher, lastJSON *string) {
	state, err := s.queryState(context.Background())
	if err != nil {
		log.Printf("[SSE] Error querying state: %v", err)
		return
	}

	data, err := json.Marshal(state)
	if err != nil {
		return
	}

	currentJSON := string(data)
	if currentJSON == *lastJSON {
		return
	}
	*lastJSON = currentJSON

	event := SSEEvent{Type: "state", Data: state}
	eventData, _ := json.Marshal(event)
	fmt.Fprintf(w, "data: %s\n\n", eventData)
	flusher.Flush()
}

func (s *Server) sendChatUpdate(w http.ResponseWriter, flusher http.Flusher, lastJSON *string) {
	if s.chatWorkflowID == "" {
		return
	}

	result, err := s.registry.QueryWorkflow(context.Background(), s.chatWorkflowID, workflow.QueryChatHistory)
	if err != nil {
		return
	}

	data, err := json.Marshal(result)
	if err != nil {
		return
	}

	currentJSON := string(data)
	if currentJSON == *lastJSON {
		return
	}
	*lastJSON = currentJSON

	event := SSEEvent{Type: "chat", Data: result}
	eventData, _ := json.Marshal(event)
	fmt.Fprintf(w, "data: %s\n\n", eventData)
	flusher.Flush()
}
