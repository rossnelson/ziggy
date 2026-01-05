package workflow

import (
	"fmt"
	"time"
)

type ChatMessage struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"` // "user" | "ziggy"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type ChatState struct {
	Owner           string        `json:"owner"`
	Messages        []ChatMessage `json:"messages"`
	ActiveMystery   *Mystery      `json:"activeMystery,omitempty"`
	MysteryProgress int           `json:"mysteryProgress"`
	HintsGiven      []string      `json:"hintsGiven"`
	Solved          []string      `json:"solved"`
	CreatedAt       time.Time     `json:"createdAt"`
	LastMessageAt   time.Time     `json:"lastMessageAt"`
	IsTyping        bool          `json:"isTyping"`
}

type Mystery struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Track       string   `json:"track"` // "educational" | "fun"
	Hints       []string `json:"hints"`
	Solution    string   `json:"solution"`
	Concept     string   `json:"concept,omitempty"` // Temporal concept for educational
	Summary     string   `json:"summary,omitempty"` // Educational summary for the topic
}

type MysteryStatus struct {
	Active      bool     `json:"active"`
	Mystery     *Mystery `json:"mystery,omitempty"`
	Progress    int      `json:"progress"`
	HintsGiven  []string `json:"hintsGiven"`
	TotalHints  int      `json:"totalHints"`
}

// ChatHistoryResponse is returned by the chat history query
type ChatHistoryResponse struct {
	Messages       []ChatMessage  `json:"messages"`
	MysteryStatus  *MysteryStatus `json:"mysteryStatus,omitempty"`
}

func NewChatState(owner string) ChatState {
	return ChatState{
		Owner:     owner,
		Messages:  []ChatMessage{},
		HintsGiven: []string{},
		Solved:    []string{},
		CreatedAt: time.Now(),
	}
}

func (s *ChatState) AddMessage(role, content string, timestamp time.Time) ChatMessage {
	msg := ChatMessage{
		ID:        generateMessageID(len(s.Messages)),
		Role:      role,
		Content:   content,
		Timestamp: timestamp,
	}
	s.Messages = append(s.Messages, msg)
	s.LastMessageAt = timestamp
	return msg
}

func (s *ChatState) GetMysteryStatus() MysteryStatus {
	if s.ActiveMystery == nil {
		return MysteryStatus{Active: false}
	}
	return MysteryStatus{
		Active:     true,
		Mystery:    s.ActiveMystery,
		Progress:   s.MysteryProgress,
		HintsGiven: s.HintsGiven,
		TotalHints: len(s.ActiveMystery.Hints),
	}
}

func generateMessageID(index int) string {
	return fmt.Sprintf("msg-%d-%d", time.Now().UnixNano(), index)
}
