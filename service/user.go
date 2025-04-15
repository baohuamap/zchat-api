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
	GetSentFriendRequests(c context.Context, userID uint64) ([]models.Friendship, error)
	GetReceivedFriendRequests(c context.Context, userID uint64) ([]models.Friendship, error)
	GetFriends(c context.Context, userID uint64) ([]models.User, error)
	// UploadAvatar(c context.Context, userID uint64, file *multipart.File) error
	FindUsers(c context.Context, search string) (*dto.FindUserListRes, error)
}

type service struct {
	repo           repo.UserRepository
	friendshipRepo repo.FriendshipRepository
	// s3Client       aws.S3Client
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

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	res := &dto.CreateUserRes{
		ID:        strconv.FormatUint(u.ID, 10),
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

	friendship := &models.Friendship{
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

func (s *service) GetSentFriendRequests(c context.Context, userID uint64) ([]models.Friendship, error) {
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

func (s *service) GetReceivedFriendRequests(c context.Context, userID uint64) ([]models.Friendship, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	// Check if user exists
	if _, err := s.repo.Get(ctx, userID); err != nil {
		slog.Error("User not found", "userID", userID)
		return nil, err
	}

	friendships, err := s.friendshipRepo.GetByFriendID(ctx, userID)
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

	var friendships []models.Friendship
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
			if friendship.UserID == userID {
				friend, err := s.repo.Get(ctx, friendship.FriendID)
				if err != nil {
					slog.Error("Error getting friend", "friendID", friendship.FriendID)
					continue
				}
				friends = append(friends, friend)
			} else if friendship.FriendID == userID {
				friend, err := s.repo.Get(ctx, friendship.UserID)
				if err != nil {
					slog.Error("Error getting friend", "friendID", friendship.FriendID)
					continue
				}
				friends = append(friends, friend)
			}
		}
	}

	return friends, nil
}

// func (s *service) UploadAvatar(c context.Context, userID uint64, file *multipart.File) error {
// 	ctx, cancel := context.WithTimeout(c, 10*time.Second)
// 	defer cancel()

// 	// Check if user exists
// 	if _, err := s.repo.Get(ctx, userID); err != nil {
// 		slog.Error("User not found", "userID", userID)
// 		return err
// 	}

// 	// Upload file to s3
// 	filePath, err := s.s3Client.UploadFile(ctx, "avatars/"+strconv.FormatUint(userID, 10), *file)

// 	// Update user avatar
// 	user := &models.User{
// 		ID:     userID,
// 		Avatar: filePath,
// 	}

// 	return s.repo.Update(ctx, user)
// }

func (s *service) FindUsers(c context.Context, search string) (*dto.FindUserListRes, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	users, err := s.repo.Search(ctx, search)
	if err != nil {
		slog.Error("User not found", "search", search)
		return &dto.FindUserListRes{}, err
	}

	var res dto.FindUserListRes
	for _, u := range users {
		res.Users = append(res.Users, dto.FindUserRes{
			ID:        strconv.FormatUint(u.ID, 10),
			Username:  u.Username,
			Email:     u.Email,
			Phone:     u.Phone,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Avatar:    u.Avatar,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
			UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &res, nil
}
