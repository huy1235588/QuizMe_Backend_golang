package websocket

import (
	"encoding/json"
	"time"
)

// MessageType defines the type of WebSocket message
type MessageType string

const (
	// Game events
	MessageTypeGameStart      MessageType = "GAME_START"
	MessageTypeGameEnd        MessageType = "GAME_END"
	MessageTypeQuestion       MessageType = "QUESTION"
	MessageTypeTimer          MessageType = "TIMER"
	MessageTypeAnswer         MessageType = "ANSWER"
	MessageTypeAnswerResult   MessageType = "ANSWER_RESULT"
	MessageTypeQuestionResult MessageType = "QUESTION_RESULT"
	MessageTypeLeaderboard    MessageType = "LEADERBOARD"
	MessageTypeNextQuestion   MessageType = "NEXT_QUESTION"

	// Room events
	MessageTypeJoin        MessageType = "JOIN"
	MessageTypeLeave       MessageType = "LEAVE"
	MessageTypeChat        MessageType = "CHAT"
	MessageTypeParticipant MessageType = "PARTICIPANT"

	// Connection events
	MessageTypeConnect           MessageType = "CONNECT"
	MessageTypeDisconnect        MessageType = "DISCONNECT"
	MessageTypePlayerDisconnect  MessageType = "PLAYER_DISCONNECT"
	MessageTypePlayerReconnect   MessageType = "PLAYER_RECONNECT"
	MessageTypePlayerTimeout     MessageType = "PLAYER_TIMEOUT"
	MessageTypeError             MessageType = "ERROR"
	MessageTypePing              MessageType = "PING"
	MessageTypePong              MessageType = "PONG"
)

// Message represents a WebSocket message
type Message struct {
	Type      MessageType     `json:"type"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

// NewMessage creates a new message with the given type and payload
func NewMessage(msgType MessageType, payload interface{}) (*Message, error) {
	var payloadBytes json.RawMessage
	var err error

	if payload != nil {
		payloadBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	return &Message{
		Type:      msgType,
		Payload:   payloadBytes,
		Timestamp: time.Now().UTC(),
	}, nil
}

// ParsePayload parses the message payload into the given target
func (m *Message) ParsePayload(target interface{}) error {
	if m.Payload == nil {
		return nil
	}
	return json.Unmarshal(m.Payload, target)
}

// JoinPayload is the payload for join messages
type JoinPayload struct {
	RoomID    uint    `json:"roomId"`
	UserID    *uint   `json:"userId,omitempty"`
	GuestName *string `json:"guestName,omitempty"`
	Token     string  `json:"token,omitempty"`
}

// ChatPayload is the payload for chat messages
type ChatPayload struct {
	RoomID    uint   `json:"roomId"`
	Content   string `json:"content"`
	GuestName string `json:"guestName,omitempty"`
}

// AnswerPayload is the payload for answer submission
type AnswerPayload struct {
	QuestionID      uint    `json:"questionId"`
	SelectedOptions []uint  `json:"selectedOptions"`
	AnswerTime      float64 `json:"answerTime"`
}

// TimerPayload is the payload for timer events
type TimerPayload struct {
	RemainingTime int `json:"remainingTime"`
	TotalTime     int `json:"totalTime"`
}

// ErrorPayload is the payload for error messages
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
