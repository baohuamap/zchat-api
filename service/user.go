package service

import (
	"context"
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
}

type service struct {
	repo repo.UserRepository
}

func NewUserService(r repo.UserRepository) User {
	return &service{
		r,
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
