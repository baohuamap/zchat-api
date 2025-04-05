package service

import (
	"context"
	"io"
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
	UploadImage(c context.Context, file io.Reader, filename string) (string, error)
	GetProfileImage(c context.Context, userID string) ([]byte, error)
}

type service struct {
	repo     repo.UserRepository
	uploader Uploader
}

func NewUserService(r repo.UserRepository, uploader Uploader) User {
	return &service{
		repo:     r,
		uploader: uploader,
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

		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.repo.Create(ctx, *u); err != nil {
		return nil, err
	}

	res := &dto.CreateUserRes{
		ID:       strconv.Itoa(int(u.ID)),
		Username: u.Username,
		Email:    u.Email,
	}

	return res, nil
}

func (s *service) Login(c context.Context, req *dto.LoginUserReq) (*dto.LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	u, err := s.repo.GetByEmail(ctx, req.Phone)
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

func (s *service) UploadImage(c context.Context, file io.Reader, filename string) (string, error) {
	// Sử dụng uploader để đọc file và chuyển thành dữ liệu binary
	fileData, err := s.uploader.UploadFileToDB(c, file, filename)
	if err != nil {
		return "", err
	}

	// Lưu dữ liệu binary vào cơ sở dữ liệu thông qua repository
	userID := uint(1) // Ví dụ: userID được lấy từ context hoặc JWT
	err = s.repo.SaveProfileImageData(c, userID, fileData)
	if err != nil {
		return "", err
	}

	return "Image uploaded successfully", nil
}

func (s *service) GetProfileImage(c context.Context, userID string) ([]byte, error) {
	// Gọi repository để lấy dữ liệu ảnh từ cơ sở dữ liệu
	return s.repo.GetProfileImageData(c, userID)
}
