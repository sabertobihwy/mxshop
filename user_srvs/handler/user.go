package handler

import (
	context "context"
	"crypto/sha512"
	"fmt"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	"mxshop_srvs/user_srvs/global"
	. "mxshop_srvs/user_srvs/model"
	"mxshop_srvs/user_srvs/proto"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func Paginate(pg, pgSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case pgSize > 100:
			pgSize = 100
		case pgSize <= 0:
			pgSize = 10
		}
		offset := (pg - 1) * pgSize
		return db.Offset(offset).Limit(pgSize)
	}
}

func ModelToRsp(usr User) proto.UserInfoRsp {
	rsp := proto.UserInfoRsp{
		Id:       usr.ID,
		PassWord: usr.Password,
		NickName: usr.NickName,
		Gender:   usr.Gender,
		Role:     int32(usr.Role),
	}
	if usr.Birthday != nil {
		rsp.BirthDay = uint64(usr.Birthday.Unix())
	}
	return rsp
}

func (s *UserServer) GetUserList(ctx context.Context, in *proto.PageInfo) (*proto.UserListRsp, error) {
	var users []User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListRsp{}
	rsp.Total = int32(result.RowsAffected)
	global.DB.Scopes(Paginate(int(in.Pn), int(in.PSize))).Find(&users)
	for _, user := range users {
		UserInfoRsp := ModelToRsp(user)
		rsp.Data = append(rsp.Data, &UserInfoRsp)
	}
	return rsp, nil
}
func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoRsp, error) {
	var user User
	result := global.DB.Where(&User{Mobile: req.Mobile}).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not exists")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := ModelToRsp(user)
	return &rsp, nil
}

func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoRsp, error) {
	var user User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not exists")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := ModelToRsp(user)
	return &rsp, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoRsp, error) {
	var user User
	result := global.DB.Where(&User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}
	user.Mobile = req.Mobile
	user.NickName = req.NickName

	options := &password.Options{16, 100, 30, sha512.New}
	salt, encodedPwd := password.Encode("generic password", options)
	user.Password = fmt.Sprintf("$sha512$%s$%s", salt, encodedPwd)
	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	// create user后database生成id -> 返回给调用方
	rsp := ModelToRsp(user)
	return &rsp, nil
}
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	user.NickName = req.NickName
	user.Gender = req.Gender
	unix := time.Unix(int64(req.BirthDay), 0)
	user.Birthday = &unix
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}
func (s *UserServer) CheckPwd(ctx context.Context, req *proto.PwdCheckInfo) (*proto.CheckRsp, error) {
	options := &password.Options{16, 100, 30, sha512.New}
	//newPwd := fmt.Sprintf("$sha512$%s$%s", salt, encodedPwd)
	pwdInfo := strings.Split(req.EncryptedPws, "$")
	check := password.Verify(req.PassWord, pwdInfo[2], pwdInfo[3], options)
	return &proto.CheckRsp{Success: check}, nil
}
