package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"mxshop_srvs/user_srvs/proto"
)

var (
	userClient proto.UserClient
	conn       *grpc.ClientConn
)

func Init() {
	var err error
	conn, err = grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		panic("conn error")
	}
	userClient = proto.NewUserClient(conn)

}

func main() {
	Init()
	userListRsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    2,
		PSize: 5,
	})
	if err != nil {
		panic(err)
	}
	for _, userInfo := range userListRsp.Data {
		fmt.Println(userInfo.NickName, userInfo.Mobile, userInfo.PassWord)
		rsp, _ := userClient.CheckPwd(context.Background(), &proto.PwdCheckInfo{PassWord: "admin123", EncryptedPws: userInfo.PassWord})
		fmt.Println(rsp.Success)
	}

	err = conn.Close()
	if err != nil {
		return
	}
}
