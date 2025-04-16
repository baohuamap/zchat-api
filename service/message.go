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
