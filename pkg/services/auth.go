package services

import (
	"context"
	"github.com/hovhannesyan/RiskIndex-AuthSVC/pkg/db"
	"github.com/hovhannesyan/RiskIndex-AuthSVC/pkg/models"
	"github.com/hovhannesyan/RiskIndex-AuthSVC/pkg/pb"
	"github.com/hovhannesyan/RiskIndex-AuthSVC/pkg/utils"
	"net/http"
)

type Server struct {
	DbHandler db.Handler
	Jwt       utils.JwtWrapper
	pb.UnimplementedAuthServiceServer
}

func (s *Server) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	var user models.User
	var err error

	if result := s.DbHandler.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error == nil {
		return &pb.SignUpResponse{
			Status: http.StatusConflict,
			Error:  "e-mail already exists",
		}, nil
	}

	user.Email = req.Email
	user.Password, err = utils.HashPassword(req.Password)

	if err != nil {
		return &pb.SignUpResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	s.DbHandler.DB.Create(&user)

	return &pb.SignUpResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	var user models.User

	if result := s.DbHandler.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error != nil {
		return &pb.SignInResponse{
			Status: http.StatusNotFound,
			Error:  "user not found",
		}, nil
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)

	if !match {
		return &pb.SignInResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	token, err := s.Jwt.GenerateToken(user)

	if err != nil {
		return &pb.SignInResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	return &pb.SignInResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := s.Jwt.ValidateToken(req.Token)

	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	var user models.User

	if result := s.DbHandler.DB.Where(&models.User{Email: claims.Email}).First(&user); result.Error != nil {
		return &pb.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: user.Id,
	}, nil
}
