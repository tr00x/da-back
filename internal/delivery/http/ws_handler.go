package http

import (
	"dubai-auto/internal/model"
	"dubai-auto/internal/service"
	"dubai-auto/internal/utils"
	"dubai-auto/pkg/auth"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type SocketHandler struct {
	service        *service.SocketService
	wsUserConns    map[int]*websocket.Conn
	wsMutex        sync.RWMutex
	tmpMessages    map[string]int
	tmpMutex       sync.Mutex
	connWriteMutex map[*websocket.Conn]*sync.Mutex
	connWriteMuMux sync.RWMutex
}

type SocketHandlerOption func(*SocketHandler)

func NewSocketHandler(service *service.SocketService, opts ...SocketHandlerOption) *SocketHandler {
	handler := &SocketHandler{
		service:        service,
		wsUserConns:    make(map[int]*websocket.Conn),
		tmpMessages:    make(map[string]int),
		connWriteMutex: make(map[*websocket.Conn]*sync.Mutex),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(handler)
		}
	}

	return handler
}

func (h *SocketHandler) closeConnGracefully(c *websocket.Conn) {

	if c == nil {
		return
	}

	h.connWriteMuMux.Lock()
	writeMu, exists := h.connWriteMutex[c]

	if exists && writeMu != nil {
		writeMu.Lock()
		defer writeMu.Unlock()
	}

	h.connWriteMuMux.Unlock()
	_ = c.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(1*time.Second))
	_ = c.Close()
	h.connWriteMuMux.Lock()
	delete(h.connWriteMutex, c)
	h.connWriteMuMux.Unlock()
}

func (h *SocketHandler) safeWriteJSON(c *websocket.Conn, msg any) error {
	if c == nil {
		return fmt.Errorf("connection is nil")
	}

	h.connWriteMuMux.RLock()
	writeMu, exists := h.connWriteMutex[c]
	h.connWriteMuMux.RUnlock()

	if !exists {
		h.connWriteMuMux.Lock()
		writeMu, exists = h.connWriteMutex[c]

		if !exists {
			writeMu = &sync.Mutex{}
			h.connWriteMutex[c] = writeMu
		}
		h.connWriteMuMux.Unlock()
	}

	writeMu.Lock()
	defer writeMu.Unlock()

	return c.WriteJSON(msg)
}

func (h *SocketHandler) handleHeartbeat(conn *websocket.Conn, userID int, heartbeatTimeout chan struct{}, heartbeatCh chan struct{}, done chan struct{}) {
	missCount := 0
	for {
		select {
		case <-done:
			return
		default:
		}

		_ = h.safeWriteJSON(conn, model.WSMessage{Event: "ping"})

		select {
		case <-heartbeatCh:
			missCount = 0
			time.Sleep(3 * time.Second)
		case <-time.After(3 * time.Second):
			missCount++
			if missCount >= 3 {
				log.Printf("⛔ Closing idle websocket for user %d due to heartbeat timeout", userID)
				heartbeatTimeout <- struct{}{}
				return
			}
		}
	}
}

func (h *SocketHandler) SetupWebSocketHandler() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		token := c.Query("token")
		user, err := auth.ValidateWSJWT(token)

		if err != nil {
			h.sendErrorAndCloseConn(c, "Authentication failed")
			return
		}

		err = h.service.CheckUserExists(user.ID)

		if err != nil {
			h.sendErrorAndCloseConn(c, "User not found")
			return
		}

		user.Conn = c
		h.wsMutex.Lock()

		if old, exists := h.wsUserConns[user.ID]; exists && old != nil && old != c {
			h.closeConnGracefully(old)
		}
		fmt.Println("user connected", user.ID)
		h.wsUserConns[user.ID] = c
		h.wsMutex.Unlock()
		user.Avatar, user.Username, err = h.service.GetUserAvatarAndUsername(user.ID)

		if err != nil {
			h.sendErrorAndCloseConn(c, "Error getting user avatar and username")
			return
		}

		err = h.service.UpdateUserStatus(user.ID, true)

		if err != nil {
			log.Printf("❌ Error updating user status: %v", err)
		}

		defer h.closeConnAndUpdateUserStatus(user.ID, c)
		data, err := h.service.GetNewMessages(user.ID)

		if err != nil {
			log.Printf("❌ Error getting unread messages: %v", err)
		} else if data != nil {
			fmt.Printf("📤 Sending new messages to user %d, message_id: %d \n", user.ID, data[0].Messages[0].ID)
			h.sendToUser(user.ID, "new_message", data)
		}

		heartbeatCh := make(chan struct{}, 1)
		done := make(chan struct{})
		heartbeatTimeout := make(chan struct{})
		defer close(done)
		go h.handleHeartbeat(c, user.ID, heartbeatTimeout, heartbeatCh, done)

		for {
			var msg model.WSMessageReceived

			select {
			case <-heartbeatTimeout:
				log.Printf("ℹ️ Exiting read loop for user %d after heartbeat timeout", user.ID)
				h.closeConnGracefully(c)
				return
			default:
			}

			if err := c.ReadJSON(&msg); err != nil {
				log.Printf("❌ Read error: %v", err)
				break
			}

			select {
			case heartbeatCh <- struct{}{}:
			default:
			}

			switch msg.Event {
			case "ping":

			case "private_message":
				s := false
				senderID := user.ID
				var messageReceived model.MessageReceived
				// todo: add senderID. if 0. this is admin. and emit it as sender id 0

				if b, err := json.Marshal(msg.Data); err == nil {
					if err := json.Unmarshal(b, &messageReceived); err != nil {
						log.Printf("❌ Error decoding private_message data: %v", err)
						break
					}

				} else {
					log.Printf("❌ Error marshaling private_message data: %v", err)
					break
				}

				h.sendToUser(user.ID, "ack", messageReceived.Time)

				if msg.TargetUserID == 0 {
					msg.TargetUserID, _ = h.service.GetActiveAdminWithChatPermission()
				}

				conversationID, err := h.service.UpsertConversation(user.ID, msg.TargetUserID, messageReceived.Message, messageReceived.Type, messageReceived.Time)

				if err != nil {
					log.Printf("❌ Error upserting conversation: %v", err)
					break
				}

				if messageReceived.Admin {
					senderID = 0
				}

				data := []model.UserMessage{
					{
						ID:             senderID,
						Username:       user.Username,
						Avatar:         &user.Avatar,
						LastActiveDate: time.Now(),
						ConversationID: conversationID,
						Messages: []model.Message{
							{
								CreatedAt:      messageReceived.Time,
								Message:        messageReceived.Message,
								Type:           messageReceived.Type,
								ConversationID: conversationID,
								SenderID:       senderID,
							},
						},
					},
				}

				// Send to all receivers
				h.wsMutex.RLock()
				targetC, exists := h.wsUserConns[msg.TargetUserID]
				h.wsMutex.RUnlock()

				data[0].Messages[0].ID, err = h.service.MessageWriteToDatabase(user.ID, s, data[0], msg.TargetUserID, user.ID)

				if err != nil {
					log.Printf("❌ Error message write to db  %d: %v", msg.TargetUserID, err)
				}

				if exists && targetC != nil {
					s = true
					h.sendToUser(msg.TargetUserID, "new_message", data)
					key := fmt.Sprintf("%d|%s", msg.TargetUserID, messageReceived.Time.Format(time.RFC3339Nano))
					h.tmpMutex.Lock()
					h.tmpMessages[key] = 1
					h.tmpMutex.Unlock()
					// Retry mechanism for each receiver
					go h.CheckReceived(msg.TargetUserID, data[0], targetC, user.ID)
				}

			case "ack":

				var t time.Time
				switch v := msg.Data.(type) {
				case string:
					parsed, err := time.Parse(time.RFC3339Nano, v)

					if err != nil {

						parsed, err = time.Parse(time.RFC3339, v)
					}

					if err != nil {
						log.Printf("❌ Error parsing ack time: %v", err)
						break
					}

					t = parsed
					// todo: move this switch case. need just 1st case's inside
				case map[string]any:

					if s, ok := v["time"].(string); ok {
						parsed, err := time.Parse(time.RFC3339Nano, s)

						if err != nil {
							parsed, err = time.Parse(time.RFC3339, s)
						}

						if err != nil {
							log.Printf("❌ Error parsing ack time from map: %v", err)
							break
						}

						t = parsed
					} else {
						log.Printf("❌ Unexpected ack data shape: %#v", v)
						break
					}
				case float64:
					fmt.Println("ack data type float64")

					t = time.Unix(int64(v), 0)
				default:
					fmt.Println("ack data type default")
					log.Printf("❌ Unsupported ack data type: %T", v)
				}

				key := fmt.Sprintf("%d|%s", user.ID, t.Format(time.RFC3339Nano))
				h.tmpMutex.Lock()
				delete(h.tmpMessages, key)
				h.tmpMutex.Unlock()

			default:
				log.Printf("⚠️ Unknown event: %s", msg.Event)
			}
		}
	})
}

func (h *SocketHandler) closeConnAndUpdateUserStatus(userID int, conn *websocket.Conn) {
	h.wsMutex.Lock()
	delete(h.wsUserConns, userID)
	h.wsMutex.Unlock()
	h.service.UpdateUserStatus(userID, false)
	h.closeConnGracefully(conn)
}

func (h *SocketHandler) CheckReceived(receiverID int, data model.UserMessage, conn *websocket.Conn, senderID int) {
	data.ID = senderID
	data.Messages[0].SenderID = senderID

	for {
		time.Sleep(3 * time.Second)
		k := fmt.Sprintf("%d|%s", receiverID, data.Messages[0].CreatedAt.Format(time.RFC3339Nano))
		h.tmpMutex.Lock()
		attempts, exists := h.tmpMessages[k]

		if !exists {
			h.tmpMutex.Unlock()
			return
		}

		if attempts >= 3 {
			delete(h.tmpMessages, k)
			h.tmpMutex.Unlock()
			// update message status and send to token
			err := h.service.MarkMessageAsUnreadAndSendPush(data.ID, data, receiverID)

			if err != nil {
				log.Printf("❌ Error marking message as unread and sending push to user %d: %v", receiverID, err)
			}

			h.closeConnGracefully(conn)
			return
		}

		h.tmpMessages[k] = attempts + 1
		h.tmpMutex.Unlock()
	}
}

func (h *SocketHandler) sendToUser(userID int, event string, data any) {
	h.wsMutex.RLock()
	userConn, exists := h.wsUserConns[userID]

	if !exists || userConn == nil {
		h.wsMutex.RUnlock()
		log.Printf("❌ User %d not connected", userID)
		return
	}

	h.wsMutex.RUnlock()
	log.Printf("📤 Sending %s to user %d", event, userID)

	msg := model.WSMessage{
		Event:        event,
		TargetUserID: userID,
		Data:         data,
	}

	go h.safeWriteJSON(userConn, msg)
}

func (h *SocketHandler) sendErrorAndCloseConn(conn *websocket.Conn, errMsg string) {
	if conn != nil {
		msg := model.WSMessage{
			Event: "error",
			Data:  errMsg,
		}

		h.safeWriteJSON(conn, msg)
		h.closeConnGracefully(conn)
	}
}

// GetConversations godoc
// @Summary      Get conversations
// @Description  Returns a list of conversations for the user
// @Tags         ws
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  model.Conversation
// @Failure      400  {object}  model.ResultMessage
// @Failure      401  {object}  auth.ErrorResponse
// @Failure      403  {object}  auth.ErrorResponse
// @Failure      500  {object}  model.ResultMessage
// @Router       /ws/conversations [get]
func (h *SocketHandler) GetConversations(c *fiber.Ctx) error {
	userID := c.Locals("id").(int)
	data := h.service.GetConversations(userID)
	return utils.FiberResponse(c, data)
}

// GetConversationMessages godoc
// @Summary      Get messages
// @Description  Returns a list of messages for the user
// @Tags         ws
// @Produce      json
// @Security     BearerAuth
// @Param        conversation_id  path      string  true  "Conversation ID"
// @Param        last_id  query      string  false  "Last message ID"
// @Param        limit  query      string  false  "Limit"
// @Success      200  {array}  model.Message
// @Failure      400  {object}  model.ResultMessage
// @Failure      401  {object}  auth.ErrorResponse
// @Failure      403  {object}  auth.ErrorResponse
// @Failure      500  {object}  model.ResultMessage
// @Router       /ws/conversations/{conversation_id}/messages [get]
func (h *SocketHandler) GetConversationMessages(c *fiber.Ctx) error {
	conversationID := c.Params("conversation_id")
	lastMessageID := c.Query("last_id")
	limit := c.Query("limit")
	userID := c.Locals("id").(int)
	data := h.service.GetConversationMessages(c.Context(), userID, conversationID, lastMessageID, limit)
	return utils.FiberResponse(c, data)
}
