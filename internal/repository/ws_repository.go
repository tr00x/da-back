package repository

import (
	"context"
	"dubai-auto/internal/config"
	"dubai-auto/internal/model"
	"dubai-auto/pkg/firebase"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SocketRepository struct {
	db              *pgxpool.Pool
	firebaseService *firebase.FirebaseService
	config          *config.Config
}

func NewSocketRepository(db *pgxpool.Pool, firebaseService *firebase.FirebaseService, config *config.Config) *SocketRepository {
	return &SocketRepository{db, firebaseService, config}
}

func (r *SocketRepository) UpdateUserStatus(userID int, status bool) error {
	q := `
		UPDATE users 
		SET online = $1, last_active_date = now() 
		WHERE id = $2
	`
	_, err := r.db.Exec(context.Background(), q, status, userID)
	return err
}

func (r *SocketRepository) GetNewMessages(userID int) ([]model.UserMessage, error) {
	q := `
		WITH updated_messages AS (
			UPDATE messages
			SET status = 2
			WHERE sender_id != $1 AND status = 1 AND conversation_id in (select id from conversations where user_id_1 = $1 or user_id_2 = $1)
			RETURNING id, sender_id, message, type, created_at, conversation_id
		)
		SELECT 
			u.id,
			u.username,
			u.last_active_date,
			p.avatar,
			c.id as conversation_id,
			json_agg(
				json_build_object(
					'id', m.id,
					'message', m.message,
					'type', m.type,
					'conversation_id', m.conversation_id,
					'created_at', to_char(m.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
					'sender_id', m.sender_id
				)
			) as messages
		FROM updated_messages m
		LEFT JOIN users u ON m.sender_id = u.id 
		left join conversations c on (c.user_id_1 = $1 and c.user_id_2 = u.id) or (c.user_id_2 = $1 and c.user_id_1 = u.id)
		LEFT JOIN profiles p ON u.id = p.user_id
		GROUP BY u.id, p.avatar, c.id;
	`
	rows, err := r.db.Query(context.Background(), q, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var messages []model.UserMessage

	for rows.Next() {
		var message model.UserMessage
		err := rows.Scan(&message.ID, &message.Username, &message.LastActiveDate, &message.Avatar, &message.ConversationID, &message.Messages)

		if err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}

	return messages, err
}

func (r *SocketRepository) GetUserAvatarAndUsername(userID int) (string, string, error) {
	q := `
		SELECT $2 || avatar, username FROM profiles WHERE user_id = $1
	`
	var avatar, username string
	var avatarP, usernameP *string
	err := r.db.QueryRow(context.Background(), q, userID, r.config.IMAGE_BASE_URL).Scan(&avatarP, &usernameP)

	if avatarP == nil {
		avatar = ""
	} else {
		avatar = *avatarP
	}

	if usernameP == nil {
		username = ""
	} else {
		username = *usernameP
	}

	return avatar, username, err
}

func (r *SocketRepository) MessageWriteToDatabase(senderUserID int, status bool, data model.UserMessage, targetUserID int) (int, error) {
	s := 1

	if status {
		s = 2
	}

	q := `
		WITH new_message AS (
			INSERT INTO messages (
				conversation_id,
				sender_id,
				status,
				message,
				type,
				created_at
			)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		)
		UPDATE conversations c
		SET
			last_message_id = nm.id,
			last_message = $4,
			last_message_type = $5,
			updated_at = NOW(),
			user_1_unread_messages = CASE
				WHEN c.user_id_1 <> $2 THEN c.user_1_unread_messages + 1
				ELSE c.user_1_unread_messages
			END,
			user_2_unread_messages = CASE
				WHEN c.user_id_2 <> $2 THEN c.user_2_unread_messages + 1
				ELSE c.user_2_unread_messages
			END
		FROM new_message nm
		WHERE c.id = $1
		RETURNING nm.id
	`
	var id int
	err := r.db.QueryRow(context.Background(), q,
		data.Messages[0].ConversationID, senderUserID,
		s, data.Messages[0].Message, data.Messages[0].Type,
		data.Messages[0].CreatedAt).Scan(&id)

	if err != nil {
		return id, err
	}

	if !status {
		userFcmToken := ""
		q = `
			select device_token from user_tokens where user_id = $1
		`
		err = r.db.QueryRow(context.Background(), q, targetUserID).Scan(&userFcmToken)

		if err != nil {
			fmt.Println("error getting fcm token: ", err)
			return id, nil
		}

		_, err = r.firebaseService.SendToToken(userFcmToken, targetUserID, data)

		if err != nil {
			fmt.Println("error sending notification: ", err)
		}
	}

	return id, nil
}

func (r *SocketRepository) MarkMessageAsUnreadAndSendPush(senderUserID int, data model.UserMessage, targetUserID int) error {

	q := `
		UPDATE messages
		SET status = 1
		WHERE sender_id = $1 and conversation_id = $2 and created_at = $3
	`
	_, err := r.db.Exec(context.Background(), q, senderUserID, data.Messages[0].ConversationID, data.Messages[0].CreatedAt)

	if err != nil {
		return err
	}

	userFcmToken := ""
	q = `
		select device_token from user_tokens where user_id = $1
	`
	err = r.db.QueryRow(context.Background(), q, targetUserID).Scan(&userFcmToken)

	if err != nil {
		return err
	}

	_, err = r.firebaseService.SendToToken(userFcmToken, targetUserID, data)

	return err
}

func (r *SocketRepository) GetUserAvatarName(userID int) (string, string, error) {
	q := `
		SELECT 
			username,
			$2 || avatar
		FROM profile 
		WHERE user_id = $1
	`
	var username string
	var avatar string
	err := r.db.QueryRow(context.Background(), q, userID, r.config.IMAGE_BASE_URL).Scan(&username, &avatar)
	return username, avatar, err
}

func (r *SocketRepository) CheckUserExists(userID int) error {
	q := `
		SELECT id FROM users WHERE id = $1
	`
	var id int
	err := r.db.QueryRow(context.Background(), q, userID).Scan(&id)
	return err
}

func (r *SocketRepository) GetUserToken(userID int) (string, error) {
	q := `
		SELECT device_token FROM user_tokens WHERE user_id = $1
	`
	var token string
	err := r.db.QueryRow(context.Background(), q, userID).Scan(&token)
	return token, err
}

// GetActiveAdminWithChatPermission returns ID of active admin user with "chat" permission
func (r *SocketRepository) GetActiveAdminWithChatPermission() (int, error) {
	var id int
	q := `
		SELECT id 
		FROM users 
		WHERE role_id = 0 
		AND status = 1 
		AND permissions @> '["chat"]'::jsonb
		limit 1
	`
	err := r.db.QueryRow(context.Background(), q).Scan(&id)
	fmt.Println(id, err)
	return id, err
}

func (r *SocketRepository) GetConversations(userID int) ([]model.Conversation, error) {
	q := `
		select 
			c.updated_at,
			case 
				when c.user_id_1 = $1 then c.user_1_unread_messages
				else c.user_2_unread_messages
			end as unread_messages,
			c.last_message,
			c.last_message_type,
			c.last_message_id,
			u.username,
			p.avatar,
			u.id user_id,
			c.id
		from conversations c
		join users u on u.id = 
			case 
				when c.user_id_1 = $1 then c.user_id_2 
				else c.user_id_1 
			end
		join profiles p on p.user_id = u.id
		where c.user_id_1 = $1 or c.user_id_2 = $1
		order by c.updated_at desc
	`
	rows, err := r.db.Query(context.Background(), q, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	conversations := make([]model.Conversation, 0)

	for rows.Next() {
		var conversation model.Conversation
		err := rows.Scan(&conversation.LastActiveDate, &conversation.UnreadMessages,
			&conversation.LastMessage, &conversation.LastMessageType, &conversation.LastMessageID,
			&conversation.Username, &conversation.Avatar,
			&conversation.UserID, &conversation.ID)

		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conversation)
	}

	return conversations, nil
}

func (r *SocketRepository) UpsertConversation(userID1, userID2 int, message string, messageType int, createdAt time.Time) (int, error) {
	fmt.Println(userID1, userID2)
	if userID1 > userID2 {
		userID1, userID2 = userID2, userID1
	}

	var id int
	q := `
		select id from conversations 
		where user_id_1 = $1 and user_id_2 = $2
	`
	err := r.db.QueryRow(context.Background(), q, userID1, userID2).Scan(&id)

	if err == pgx.ErrNoRows {
		q = `
			insert into conversations (user_id_1, user_id_2, updated_at) 
			values ($1, $2, $3)
			returning id
		`
		err = r.db.QueryRow(context.Background(), q, userID1, userID2, createdAt).Scan(&id)

		if err != nil {
			return 0, err
		}
	}

	return id, err
}

func (r *SocketRepository) GetConversationMessages(ctx context.Context, userID, conversationID, lastID, limit int) ([]model.ConversationMessage, error) {
	q := `
		select 
			id, 
			sender_id, 
			status,
			message, 
			type, 
			created_at 
		from messages 
		where conversation_id = $1 and id < $2 
		order by id desc 
		limit $3
	`
	rows, err := r.db.Query(ctx, q, conversationID, lastID, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	mes := make([]model.ConversationMessage, 0)

	for rows.Next() {
		var message model.ConversationMessage
		err := rows.Scan(&message.ID, &message.SenderID, &message.Status, &message.Message, &message.Type, &message.CreatedAt)

		if err != nil {
			return nil, err
		}

		mes = append(mes, message)
	}

	q = `
		UPDATE conversations 
		SET 
			user_1_unread_messages = CASE 
				WHEN user_id_1 = $2 THEN 0 
				ELSE user_1_unread_messages 
			END,
			user_2_unread_messages = CASE 
				WHEN user_id_1 != $2 THEN 0 
				ELSE user_2_unread_messages 
			END,
			updated_at = NOW()
		WHERE id = $1;
	`
	_, err = r.db.Exec(ctx, q, conversationID, userID)

	return mes, err
}
