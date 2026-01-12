package service

import (
	"context"
	"dubai-auto/internal/model"
	"dubai-auto/internal/repository"
	"dubai-auto/internal/utils"
	"strconv"
	"time"
)

type SocketService struct {
	repo *repository.SocketRepository
}

func NewSocketService(repo *repository.SocketRepository) *SocketService {
	return &SocketService{repo}
}

func (s *SocketService) UpdateUserStatus(userID int, status bool) error {
	err := s.repo.UpdateUserStatus(userID, status)
	return err
}

func (s *SocketService) UpsertConversation(senderUserID, targetUserID int, message string, msgType int, createdAt time.Time) (int, error) {
	return s.repo.UpsertConversation(senderUserID, targetUserID, message, msgType, createdAt)
}

func (s *SocketService) GetNewMessages(userID int) ([]model.UserMessage, error) {
	messages, err := s.repo.GetNewMessages(userID)
	return messages, err
}

func (s *SocketService) GetUserAvatarAndUsername(userID int) (string, string, error) {
	return s.repo.GetUserAvatarAndUsername(userID)
}

func (s *SocketService) MessageWriteToDatabase(senderUserID int, status bool, data model.UserMessage, targetUserID, senderID int) (int, error) {
	data.ID = senderID
	data.Messages[0].SenderID = senderID
	id, err := s.repo.MessageWriteToDatabase(senderUserID, status, data, targetUserID)
	return id, err
}

func (s *SocketService) MarkMessageAsUnreadAndSendPush(senderUserID int, data model.UserMessage, targetUserID int) error {
	err := s.repo.MarkMessageAsUnreadAndSendPush(senderUserID, data, targetUserID)
	return err
}

func (s *SocketService) CheckUserExists(userID int) error {
	err := s.repo.CheckUserExists(userID)
	return err
}

func (s *SocketService) GetActiveAdminWithChatPermission() (int, error) {
	return s.repo.GetActiveAdminWithChatPermission()
}

func (s *SocketService) GetConversations(userID int) model.Response {
	data, err := s.repo.GetConversations(userID)

	if err != nil {
		return model.Response{
			Error:  err,
			Status: 500,
		}
	}

	return model.Response{
		Data:   data,
		Status: 200,
	}
}

func (s *SocketService) GetConversationMessages(ctx context.Context, userID int, conversationID, lastMessageID, limitStr string) model.Response {
	conversationIDInt, err := strconv.Atoi(conversationID)

	if err != nil {
		return model.Response{
			Error:  err,
			Status: 400,
		}
	}

	lastID, limit := utils.CheckLastIDLimit(lastMessageID, limitStr, "chat")
	data, err := s.repo.GetConversationMessages(ctx, userID, conversationIDInt, lastID, limit)

	if err != nil {
		return model.Response{
			Error:  err,
			Status: 500,
		}
	}

	return model.Response{
		Data: data,
	}
}
