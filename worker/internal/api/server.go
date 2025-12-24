package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"ziggy/internal/temporal"
	"ziggy/internal/workflow"
)

type Server struct {
	registry   *temporal.Registry
	workflowID string
	port       int
}

func NewServer(registry *temporal.Registry, workflowID string, port int) *Server {
	return &Server{
		registry:   registry,
		workflowID: workflowID,
		port:       port,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("GET /api/state", s.handleGetState)
	mux.HandleFunc("POST /api/signal/feed", s.handleFeed)
	mux.HandleFunc("POST /api/signal/play", s.handlePlay)
	mux.HandleFunc("POST /api/signal/pet", s.handlePet)
	mux.HandleFunc("POST /api/signal/wake", s.handleWake)
	mux.HandleFunc("GET /api/health", s.handleHealth)

	// CORS middleware
	handler := corsMiddleware(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: handler,
	}

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	fmt.Printf("API server listening on http://localhost:%d\n", s.port)
	return server.ListenAndServe()
}

func (s *Server) handleGetState(w http.ResponseWriter, r *http.Request) {
	state, err := s.queryState(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    state,
	})
}

func (s *Server) handleFeed(w http.ResponseWriter, r *http.Request) {
	err := s.registry.SignalWorkflow(r.Context(), s.workflowID, workflow.SignalFeed, struct{}{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	state, _ := s.queryState(r.Context())
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    state,
	})
}

func (s *Server) handlePlay(w http.ResponseWriter, r *http.Request) {
	err := s.registry.SignalWorkflow(r.Context(), s.workflowID, workflow.SignalPlay, struct{}{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	state, _ := s.queryState(r.Context())
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    state,
	})
}

func (s *Server) handlePet(w http.ResponseWriter, r *http.Request) {
	err := s.registry.SignalWorkflow(r.Context(), s.workflowID, workflow.SignalPet, struct{}{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	state, _ := s.queryState(r.Context())
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    state,
	})
}

func (s *Server) handleWake(w http.ResponseWriter, r *http.Request) {
	err := s.registry.SignalWorkflow(r.Context(), s.workflowID, workflow.SignalWake, struct{}{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	state, _ := s.queryState(r.Context())
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    state,
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

func (s *Server) queryState(ctx context.Context) (*workflow.ZiggyStateResponse, error) {
	result, err := s.registry.QueryWorkflow(ctx, s.workflowID, workflow.QueryState)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	var state workflow.ZiggyState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	// Calculate current state with real time (decay applied)
	now := time.Now()
	current := state.CalculateCurrentState(now)
	// Keep the message from workflow state (only changes on actions/mood transitions)
	response := current.ToResponse(now)
	return &response, nil
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
