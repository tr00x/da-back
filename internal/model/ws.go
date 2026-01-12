package model

import (
	"time"

	"github.com/gofiber/contrib/websocket"
)

type WSMessage struct {
	Event        string `json:"event"`
	TargetUserID int    `json:"target_user_id"`
	Data         any    `json:"data"`
}

type MessageReceived struct {
	Time         time.Time `json:"time"`
	Admin        bool      `json:"admin"` // true if message is from admin
	Message      string    `json:"message"`
	TargetUserID int       `json:"target_user_id"` // todo: remove this field, it is have in WSMessageReceived
	Type         int       `json:"type"`
}

type WSMessageReceived struct {
	Event        string `json:"event"`
	TargetUserID int    `json:"target_user_id"`
	Data         any    `json:"data"`
}

type WSMessageAck struct {
	Event        string    `json:"event"`
	TargetUserID int       `json:"target_user_id"`
	Data         time.Time `json:"data"`
}

type WSUser struct {
	ID       int             `json:"id"`
	Username string          `json:"username"`
	Avatar   string          `json:"avatar"`
	RoleID   int             `json:"role_id"`
	Conn     *websocket.Conn `json:"-"`
}

type Message struct {
	CreatedAt      time.Time `json:"created_at"`
	Message        string    `json:"message"`
	ID             int       `json:"id"`
	ConversationID int       `json:"conversation_id"`
	SenderID       int       `json:"sender_id"`
	Type           int       `json:"type"`
}

type ConversationMessage struct {
	CreatedAt time.Time `json:"created_at"`
	Message   string    `json:"message"`
	ID        int       `json:"id"`
	SenderID  int       `json:"sender_id"`
	Status    int       `json:"status"`
	Type      int       `json:"type"`
}

type PushNotificationBody struct {
	Avatar    *string   `json:"avatar"`
	Name      string    `json:"name"`
	Type      int       `json:"type"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int       `json:"user_id"`
}

type UserMessage struct {
	Messages       []Message `json:"messages"`
	LastActiveDate time.Time `json:"last_active_date"`
	Username       string    `json:"username"`
	Avatar         *string   `json:"avatar"`
	ConversationID int       `json:"conversation_id"`
	ID             int       `json:"id"`
}

type FirebaseMessageData struct {
	Time     time.Time `json:"time"`
	Message  string    `json:"message"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	Type     int       `json:"type"`
}

type FirebaseMessage struct {
	Data FirebaseMessageData `json:"data"`
	Type int                 `json:"type"`
}

type Conversation struct {
	LastActiveDate  time.Time `json:"last_active_date"`
	UnreadMessages  int       `json:"unread_messages"`
	LastMessageID   int       `json:"last_message_id"`
	LastMessage     string    `json:"last_message"`
	LastMessageType int       `json:"last_message_type"`
	Username        string    `json:"username"`
	Avatar          *string   `json:"avatar"`
	UserID          int       `json:"user_id"`
	ID              int       `json:"id"`
}
