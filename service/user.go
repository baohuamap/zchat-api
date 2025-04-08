package service

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/baohuamap/zchat-api/util"

	"github.com/baohuamap/zchat-api/dto"
	"github.com/baohuamap/zchat-api/models"
	repo "github.com/baohuamap/zchat-api/repository"
	"github.com/golang-jwt/jwt"
)

const (
	secretKey = "secret"
)

type User interface {
	CreateUser(c context.Context, req *dto.CreateUserReq) (*dto.CreateUserRes, error)
	Login(c context.Context, req *dto.LoginUserReq) (*dto.LoginUserRes, error)
	AddFriend(c context.Context, userID uint64, friendID uint64) error
	AcceptFriend(c context.Context, userID uint64, friendID uint64) error
	RejectFriend(c context.Context, userID uint64, friendID uint64) error
	GetFriendRequests(c context.Context, userID uint64) ([]models.Friendship, error)
	GetFriends(c context.Context, userID uint64) ([]models.User, error)
}

type service struct {
	repo           repo.UserRepository
	friendshipRepo repo.FriendshipRepository
}

func NewUserService(r repo.UserRepository, f repo.FriendshipRepository) User {
	return &service{
		r, f,
	}
}

func (s *service) CreateUser(c context.Context, req *dto.CreateUserReq) (*dto.CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Phone:     req.Phone,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := s.repo.Create(ctx, *u); err != nil {
		return nil, err
	}

	res := &dto.CreateUserRes{
		ID:        strconv.Itoa(int(u.ID)),
		Username:  u.Username,
		Email:     u.Email,
		Phone:     u.Phone,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}

	return res, nil
}

func (s *service) Login(c context.Context, req *dto.LoginUserReq) (*dto.LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	u, err := s.repo.GetByPhone(ctx, req.Phone)
	if err != nil {
		return &dto.LoginUserRes{}, err
	}

	err = util.CheckPassword(req.Password, u.Password)
	if err != nil {
		return &dto.LoginUserRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, dto.MyJWTClaims{
		ID:       strconv.Itoa(int(u.ID)),
		Username: u.Username,
		MapClaims: jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &dto.LoginUserRes{}, err
	}

	return &dto.LoginUserRes{AccessToken: ss, Username: u.Username, ID: strconv.Itoa(int(u.ID))}, nil
}

func (s *service) AddFriend(c context.Context, userID uint64, friendID uint64) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	// Check if user exists
	if _, err := s.repo.Get(ctx, userID); err != nil {
		slog.Error("User not found", "userID", userID)
		return err
	}
	// Check if friend exists
	if _, err := s.repo.Get(ctx, friendID); err != nil {
		slog.Error("Friend not found", "friendID", friendID)
		return err
	}

	friendship := models.Friendship{
		UserID:   userID,
		FriendID: friendID,
		Status:   "pending",
	}

	return s.friendshipRepo.Create(ctx, friendship)
}

func (s *service) AcceptFriend(c context.Context, userID uint64, friendID uint64) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	// Check if user exists
	if _, err := s.repo.Get(ctx, userID); err != nil {
		slog.Error("User not found", "userID", userID)
		return err
	}
	// Check if friend exists
	if _, err := s.repo.Get(ctx, friendID); err != nil {
		slog.Error("Friend not found", "friendID", friendID)
		return err
	}

	friendship, err := s.friendshipRepo.GetByUserIDAndFriendID(ctx, userID, friendID)
	if err != nil {
		slog.Error("Friendship not found", "userID", userID, "friendID", friendID)
		return err
	}

	friendship.Status = "accepted"

	return s.friendshipRepo.Update(ctx, friendship)
}

func (s *service) RejectFriend(c context.Context, userID uint64, friendID uint64) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	// Check if user exists
	if _, err := s.repo.Get(ctx, userID); err != nil {
		slog.Error("User not found", "userID", userID)
		return err
	}
	// Check if friend exists
	if _, err := s.repo.Get(ctx, friendID); err != nil {
		slog.Error("Friend not found", "friendID", friendID)
		return err
	}

	friendship, err := s.friendshipRepo.GetByUserIDAndFriendID(ctx, userID, friendID)
	if err != nil {
		slog.Error("Friendship not found", "userID", userID, "friendID", friendID)
		return err
	}

	friendship.Status = "rejected"

	return s.friendshipRepo.Update(ctx, friendship)
}

func (s *service) GetFriendRequests(c context.Context, userID uint64) ([]models.Friendship, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	// Check if user exists
	if _, err := s.repo.Get(ctx, userID); err != nil {
		slog.Error("User not found", "userID", userID)
		return nil, err
	}

	friendships, err := s.friendshipRepo.GetByUserID(ctx, userID)
	if err != nil {
		slog.Error("Error getting friend requests", "userID", userID)
		return nil, err
	}

	var requests []models.Friendship
	for _, friendship := range friendships {
		if friendship.Status == "pending" {
			requests = append(requests, friendship)
		}
	}

	return requests, nil
}

func (s *service) GetFriends(c context.Context, userID uint64) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	// Check if user exists
	if _, err := s.repo.Get(ctx, userID); err != nil {
		slog.Error("User not found", "userID", userID)
		return nil, err
	}

	friendships, err := s.friendshipRepo.GetByUserID(ctx, userID)
	if err != nil {
		slog.Error("Error getting friends", "userID", userID)
		return nil, err
	}

	_friendships, err := s.friendshipRepo.GetByFriendID(ctx, userID)
	if err != nil {
		slog.Error("Error getting friends", "userID", userID)
		return nil, err
	}
	friendships = append(friendships, _friendships...)

	var friends []models.User
	for _, friendship := range friendships {
		if friendship.Status == "accepted" {
			friend, err := s.repo.Get(ctx, friendship.FriendID)
			if err != nil {
				slog.Error("Error getting friend", "friendID", friendship.FriendID)
				continue
			}
			friends = append(friends, friend)
		}
	}

	return friends, nil
}
