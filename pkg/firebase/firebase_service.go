package firebase

import (
	"context"
	"dubai-auto/internal/config"
	"dubai-auto/internal/model"
	"strconv"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

// Firebase service
type FirebaseService struct {
	client *messaging.Client
	ctx    context.Context
}

func InitFirebase(cfg *config.Config) (*FirebaseService, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile(cfg.FIREBASE_ACCOUNT_FILE)
	app, err := firebase.NewApp(ctx, nil, opt)

	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(ctx)

	if err != nil {
		return nil, err
	}

	return &FirebaseService{
		client: client,
		ctx:    ctx,
	}, nil
}

// Send notification to a specific device token
func (fs *FirebaseService) SendToToken(token string, targetUserID int, data model.UserMessage) (string, error) {
	var body string
	switch data.Messages[0].Type {
	case 2:
		body = "car listing shared."
	case 3:
		body = "video shared."
	case 4:
		body = "image shared."
	default:
		body = data.Messages[0].Message
	}
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: data.Username,
			Body:  body,
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ChannelID: config.ENV.FCM_CHANNEL_ID,
				Priority:  messaging.PriorityHigh,
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound: "default",
					Badge: &[]int{1}[0],
					Alert: &messaging.ApsAlert{
						Title: data.Username,
						Body:  data.Messages[0].Message,
					},
				},
			},
		},
		Data: map[string]string{
			"type":            "chat",
			"current_user_id": strconv.Itoa(targetUserID),
			"sender_id":       strconv.Itoa(data.ID), //todo: if admin 0
			"sender_name":     data.Username,
			"sender_avatar":   *data.Avatar,
			"message_id":      strconv.Itoa(data.Messages[0].ID),
			"message":         data.Messages[0].Message,
			"msg_type":        strconv.Itoa(data.Messages[0].Type),
		},
	}

	response, err := fs.client.Send(fs.ctx, message)
	return response, err
}

// Send notification to multiple tokens
func (fs *FirebaseService) SendToMultipleTokens(tokens []string, title, body string, data map[string]string) (*messaging.BatchResponse, error) {
	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	response, err := fs.client.SendEachForMulticast(fs.ctx, message)
	return response, err
}

// Send to topic (for broadcast notifications)
func (fs *FirebaseService) SendToTopic(topic, title, body string, data map[string]string) (string, error) {
	message := &messaging.Message{
		Topic: topic,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	response, err := fs.client.Send(fs.ctx, message)

	if err != nil {
		return "", err
	}

	return response, nil
}

// Subscribe tokens to a topic
func (fs *FirebaseService) SubscribeToTopic(tokens []string, topic string) error {
	_, err := fs.client.SubscribeToTopic(fs.ctx, tokens, topic)
	return err
}
