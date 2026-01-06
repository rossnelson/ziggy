package chat

import (
	"fmt"
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type State struct {
	Owner           string    `json:"owner"`
	Messages        []Message `json:"messages"`
	ActiveMystery   *Mystery  `json:"activeMystery,omitempty"`
	MysteryProgress int       `json:"mysteryProgress"`
	HintsGiven      []string  `json:"hintsGiven"`
	Solved          []string  `json:"solved"`
	CreatedAt       time.Time `json:"createdAt"`
	LastMessageAt   time.Time `json:"lastMessageAt"`
	IsTyping        bool      `json:"isTyping"`
}

type Mystery struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Track       string   `json:"track"`
	Hints       []string `json:"hints"`
	Solution    string   `json:"solution"`
	Concept     string   `json:"concept,omitempty"`
	Summary     string   `json:"summary,omitempty"`
}

type MysteryStatus struct {
	Active     bool     `json:"active"`
	Mystery    *Mystery `json:"mystery,omitempty"`
	Progress   int      `json:"progress"`
	HintsGiven []string `json:"hintsGiven"`
	TotalHints int      `json:"totalHints"`
}

type HistoryResponse struct {
	Messages      []Message      `json:"messages"`
	MysteryStatus *MysteryStatus `json:"mysteryStatus,omitempty"`
	IsTyping      bool           `json:"isTyping"`
}

func NewState(owner string) State {
	return State{
		Owner:      owner,
		Messages:   []Message{},
		HintsGiven: []string{},
		Solved:     []string{},
		CreatedAt:  time.Now(),
	}
}

func (s *State) AddMessage(role, content string, timestamp time.Time) Message {
	msg := Message{
		ID:        generateMessageID(len(s.Messages)),
		Role:      role,
		Content:   content,
		Timestamp: timestamp,
	}
	s.Messages = append(s.Messages, msg)
	s.LastMessageAt = timestamp
	return msg
}

func (s *State) GetMysteryStatus() MysteryStatus {
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
