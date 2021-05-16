package main

import (
	"context"
	pbProject "todo/proto/project"
	pbUser "todo/proto/user"
	"todo/shared/md"
	userDomain "todo/user/domain"

	"github.com/golang/protobuf/ptypes"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type userService struct {
	db            *gorm.DB
	projectClient pbProject.ProjectServiceClient
	pbUser.UnimplementedUserServiceServer
}

const DefaultProjectName = "Default Project"

func (s *userService) CreateUser(
	ctx context.Context,
	req *pbUser.CreateUserRequest,
) (*pbUser.CreateUserResponse, error) {
	if req.GetEmail() == "" || len(req.GetPassword()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "email or password is empty")
	}

	hash, err := bcrypt.GenerateFromPassword(req.GetPassword(), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	tx := s.db.Begin()
	user := userDomain.NewUser(req.GetEmail(), hash)
	user = userDomain.NewUserRepository(tx).Create(user)

	ctx = md.AddUserIDToContext(ctx, user.ID)
	projectReq := &pbProject.CreateProjectRequest{Name: DefaultProjectName}
	if _, err = s.projectClient.CreateProject(ctx, projectReq); err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		panic(err)
	}
	return &pbUser.CreateUserResponse{User: buildPbUser(user)}, nil
}

func (s *userService) FindUser(ctx context.Context, req *pbUser.FindUserRequest) (*pbUser.FindUserResponse, error) {
	u := userDomain.NewUserRepository(s.db).FindByID(req.GetUserId())
	if u == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pbUser.FindUserResponse{User: buildPbUser(u)}, nil
}

func (s *userService) VerifyUser(ctx context.Context, req *pbUser.VerifyUserRequest) (*pbUser.VerifyUserResponse, error) {
	u := userDomain.NewUserRepository(s.db).FindByEmail(req.GetEmail())
	if u == nil {
		return nil, status.Error(codes.Unauthenticated, "wrong email or password")
	}

	if err := bcrypt.CompareHashAndPassword(u.Password, req.GetPassword()); err != nil {
		return nil, status.Error(codes.Unauthenticated, "wrong email or password")
	}

	return &pbUser.VerifyUserResponse{User: buildPbUser(u)}, nil
}

func buildPbUser(u *userDomain.User) *pbUser.User {
	var createdAt *timestamppb.Timestamp
	if u.CreatedAt != nil {
		createdAt, _ = ptypes.TimestampProto(*u.CreatedAt)
	}

	return &pbUser.User{
		Id:           u.ID,
		Email:        u.Email,
		PasswordHash: u.Password,
		CreatedAt:    createdAt,
	}
}
