package handler

import (
	"context"
	"user_srv/utils"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"user_srv/database"
	"user_srv/model"
	"user_srv/proto/gen/user_pb"
)

type UserService struct{}

func convertModelUserToResponseUser(user model.User) *user_pb.UserInfoResponse {
	respUser := user_pb.UserInfoResponse{}
	respUser.Id = int32(user.ID)
	respUser.Password = user.Password
	respUser.Mobile = user.Mobile
	respUser.Nickname = user.Nickname
	return &respUser
}

func (u *UserService) GetUserList(ctx context.Context, request *user_pb.PageInfoRequest) (resp *user_pb.UserListResponse, err error) {
	var users []model.User
	db, err := database.GetDB()
	if err != nil {
		zap.S().Errorf("数据库错误:%s\n", err.Error())
		return nil, status.Error(codes.DataLoss, "获取数据库出错")
	}
	var count int64
	db.Model(&users).Count(&count)
	resp = &user_pb.UserListResponse{}
	resp.Total = int32(count)

	var page uint32 = 1
	var pageNum uint32 = 10
	if request.PageSize > 0 {
		pageNum = request.PageSize
	}
	if request.PageNum > 1 {
		page = request.PageNum
	}
	offset := pageNum * (page - 1)
	db.Offset(int(offset)).Limit(int(pageNum)).Find(&users)
	for _, value := range users {
		userInfoResp := convertModelUserToResponseUser(value)
		resp.Data = append(resp.Data, userInfoResp)
	}

	return resp, nil
}

func (u UserService) CreateUser(ctx context.Context, request *user_pb.CreateUserInfoRequest) (*user_pb.UserInfoResponse, error) {
	var user model.User
	db, err := database.GetDB()
	if err != nil {
		zap.S().Errorf("数据库错误:%s\n", err.Error())
		return nil, status.Error(codes.DataLoss, "获取数据库出错")
	}

	db.Where("mobile = ?", request.Mobile).First(&user)
	if user.ID > 0 {
		return nil, status.Error(codes.AlreadyExists, "用户已存在")
	}

	user = model.User{Nickname: request.Nickname, Password: utils.MD5(request.Password), Mobile: request.Mobile}
	result := db.Create(&user)
	if result.Error != nil {
		return nil, status.Error(codes.Unknown, "创建用户失败,未知原因")
	}
	userInfoResp := convertModelUserToResponseUser(user)
	return userInfoResp, nil
}
