package handler

import (
	"context"

	"user_srv/proto/gen/user_pb"
)

type UserService struct{}

func (u UserService) GetUserList(ctx context.Context, request *user_pb.PageInfoRequest) (*user_pb.UserListResponse, error) {
	panic("implement me")
}

func (u UserService) CreateUser(ctx context.Context, request *user_pb.CreateUserInfoRequest) (*user_pb.UserInfoResponse, error) {
	panic("implement me")
}
