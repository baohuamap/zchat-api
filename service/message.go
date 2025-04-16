package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/baohuamap/zchat-api/dto"
	repo "github.com/baohuamap/zchat-api/repository"
)

type Message interface {
	LoadConversations(c context.Context, userID uint64) (*dto.ConversationListRes, error)
	LoadMessages(c context.Context, conversationID uint64) (*dto.MessageListRes, error)
	SeenMessages(c context.Context, conversationID uint64, req *dto.SeenMessagesReq) error
}

type msgService struct {
	cRepo repo.ConversationRepository
	mRepo repo.MessageRepository
	pRepo repo.ParticipantRepository
}

func NewMessageService(convRepo repo.ConversationRepository, msgRepo repo.MessageRepository, participantRepo repo.ParticipantRepository) Message {
	return &msgService{
		cRepo: convRepo,
		mRepo: msgRepo,
		pRepo: participantRepo,
	}
}

func (s *msgService) LoadConversations(c context.Context, userID uint64) (*dto.ConversationListRes, error) {
	conversations, err := s.pRepo.GetConversationByParticipants(c, userID)
	if err != nil {
		slog.Error("Failed to get conversations", "error", err)
		return nil, err
	}

	var convRes dto.ConversationListRes
	for _, conv := range conversations {
		participants, err := s.pRepo.GetByConversationID(c, conv.ID)
		if err != nil {
			slog.Error("Failed to get participants", "error", err)
			return nil, err
		}
		participantIDs := make([]uint64, len(participants))
		for i, p := range participants {
			participantIDs[i] = p.UserID
		}
		latestMessage, err := s.mRepo.GetLatestByConversationID(c, conv.ID)
		if err != nil && err.Error() != "NotFound" {
			slog.Error("Failed to get latest message", "error", err)
			return nil, err
		}

		convRes.Conversations = append(convRes.Conversations, dto.ConversationRes{
			ID:           conv.ID,
			Type:         conv.Type,
			CreatorID:    conv.CreatorID,
			Participants: participantIDs,
			Seen:         conv.Seen,
			LatestMessageCreatedAt: func() *time.Time {
				if latestMessage != nil {
					return &latestMessage.CreatedAt
				}
				return nil
			}(),
		})
	}

	return &convRes, nil
}

func (s *msgService) LoadMessages(c context.Context, conversationID uint64) (*dto.MessageListRes, error) {
	messages, err := s.mRepo.GetByConversationID(c, conversationID)
	if err != nil {
		return nil, err
	}

	var msgRes dto.MessageListRes
	for _, msg := range messages {
		msgRes.Messages = append(msgRes.Messages, dto.MessageRes{
			Content:        msg.Content,
			SenderID:       msg.SenderID,
			CreateAt:       msg.CreatedAt,
			ConversationID: msg.ConversationID,
		})
	}

	return &msgRes, nil
}

func (s *msgService) SeenMessages(c context.Context, conversationID uint64, req *dto.SeenMessagesReq) error {

	// Check if the conversation exists
	conversation, err := s.cRepo.Get(c, conversationID)
	if err != nil {
		slog.Error("Failed to get conversation", "error", err)
		return err
	}

	// Check if the user is a participant in the conversation
	_, err = s.pRepo.GetByUserIDAndConversationID(c, req.UserID, conversationID)
	if err != nil {
		slog.Error("Failed to get participant", "error", err)
		return err
	}

	// Get the latest message for the conversation
	latestMessage, err := s.mRepo.GetLatestByConversationID(c, conversationID)
	if err != nil && err.Error() != "NotFound" {
		slog.Error("Failed to get latest message", "error", err)
		return err
	}

	// If there are no messages, return early
	if latestMessage == nil {
		slog.Info("No messages found for conversation", "conversationID", conversationID)
		return nil
	}

	if req.UserID == latestMessage.SenderID {
		slog.Info("User is the sender of the latest message, no need to update seen status", "userID", req.UserID)
		return nil
	}
	// Update the seen status of the latest message for the participant

	conversation.Seen = true
	err = s.cRepo.Update(c, conversation)
	if err != nil {
		slog.Error("Failed to update conversation seen status", "error", err)
		return err
	}

	return nil

}
