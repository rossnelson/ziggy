package api

import (
	"encoding/json"
	"net/http"

	"ziggy/internal/workflow/chat"
)

func (s *Server) handleGetChatHistory(w http.ResponseWriter, r *http.Request) {
	if s.chatWorkflowID == "" {
		writeError(w, http.StatusNotFound, "chat not initialized")
		return
	}

	result, err := s.reg.QueryWorkflow(r.Context(), s.chatWorkflowID, chat.QueryChatHistory)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

func (s *Server) handleSendMessage(w http.ResponseWriter, r *http.Request) {
	if s.chatWorkflowID == "" {
		writeError(w, http.StatusNotFound, "chat not initialized")
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Content == "" {
		writeError(w, http.StatusBadRequest, "content is required")
		return
	}

	signal := chat.SendMessageSignal{Content: req.Content}
	err := s.reg.SignalWorkflow(r.Context(), s.chatWorkflowID, chat.SignalSendMessage, signal)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Give workflow time to process and generate response
	// In production, would use a more robust approach like polling or websockets

	result, err := s.reg.QueryWorkflow(r.Context(), s.chatWorkflowID, chat.QueryChatHistory)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

func (s *Server) handleGetMysteryStatus(w http.ResponseWriter, r *http.Request) {
	if s.chatWorkflowID == "" {
		writeError(w, http.StatusNotFound, "chat not initialized")
		return
	}

	result, err := s.reg.QueryWorkflow(r.Context(), s.chatWorkflowID, chat.QueryMysteryStatus)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

func (s *Server) handleStartMystery(w http.ResponseWriter, r *http.Request) {
	if s.chatWorkflowID == "" {
		writeError(w, http.StatusNotFound, "chat not initialized")
		return
	}

	var req struct {
		MysteryID string `json:"mysteryId"`
		Track     string `json:"track"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	signal := chat.StartMysterySignal{
		MysteryID: req.MysteryID,
		Track:     req.Track,
	}
	err := s.reg.SignalWorkflow(r.Context(), s.chatWorkflowID, chat.SignalStartMystery, signal)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func (s *Server) handleGetMysteries(w http.ResponseWriter, r *http.Request) {
	track := r.URL.Query().Get("track")
	if track == "" {
		track = "fun"
	}

	// Get solved mysteries from chat workflow if available
	var solved []string
	if s.chatWorkflowID != "" {
		result, err := s.reg.QueryWorkflow(r.Context(), s.chatWorkflowID, chat.QueryMysteryStatus)
		if err == nil {
			if status, ok := result.(chat.MysteryStatus); ok {
				// Would need to track solved mysteries in state
				_ = status
			}
		}
	}

	mysteries := chat.GetAvailableMysteries(track, solved)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    mysteries,
	})
}
