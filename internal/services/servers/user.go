package servers

import (
	"context"
	"errors"
	"fmt"
	userpb "github.com/daemondxx/lks_back/gen/pb/go/user"
	"github.com/daemondxx/lks_back/internal/dao"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct {
	uServ UserService
	userpb.UnimplementedUserServiceServer
}

func NewUserServer(uServ UserService) *UserServer {
	return &UserServer{
		uServ: uServ,
	}
}

func (u UserServer) GetUserInfo(ctx context.Context, r *userpb.GetUserInfoRequest) (*userpb.UserInfo, error) {
	if uInfo, err := u.uServ.GetUserByID(ctx, uint(r.UserId)); err != nil {
		if errors.Is(err, dao.ErrUserNotFound) {
			st := status.New(codes.NotFound, "user not found")
			return nil, st.Err()
		} else {
			return nil, ErrInternal
		}
	} else {
		return &userpb.UserInfo{
			AccordLogin:    uInfo.AccordLogin,
			AccordPassword: uInfo.AccordPassword,
			LksLogin:       uInfo.LKSLogin,
			LksPassword:    uInfo.LKSPassword,
		}, nil
	}
}

func (u UserServer) ChangeUserStatus(ctx context.Context, r *userpb.ChangeStatusRequest) (*emptypb.Empty, error) {
	if err := u.uServ.UpdateActiveStatus(ctx, uint(r.UserId), r.ActiveStatus); err != nil {
		if errors.Is(err, dao.ErrUserNotFound) {
			st := status.New(codes.InvalidArgument, fmt.Sprintf("user with id %d not found", r.UserId))
			return nil, st.Err()
		} else {
			return nil, ErrInternal
		}
	}
	return &emptypb.Empty{}, nil
}

func (u UserServer) UpdateAccord(ctx context.Context, r *userpb.UpdateRequest) (*emptypb.Empty, error) {
	if err := u.uServ.UpdateAccord(ctx, uint(r.UserId), r.Login, r.Password); err != nil {
		if errors.Is(err, dao.ErrUserNotFound) {
			st := status.New(codes.InvalidArgument, fmt.Sprintf("user with id %d not found", r.UserId))
			return nil, st.Err()
		} else {
			return nil, ErrInternal
		}
	}
	return &emptypb.Empty{}, nil
}

func (u UserServer) UpdateLks(ctx context.Context, r *userpb.UpdateRequest) (*emptypb.Empty, error) {
	if err := u.uServ.UpdateLKS(ctx, uint(r.UserId), r.Login, r.Password); err != nil {
		if errors.Is(err, dao.ErrUserNotFound) {
			st := status.New(codes.InvalidArgument, fmt.Sprintf("user with id %d not found", r.UserId))
			return nil, st.Err()
		} else {
			return nil, ErrInternal
		}
	}
	return &emptypb.Empty{}, nil
}

func (u UserServer) mustEmbedUnimplementedUserServiceServer() {
	//TODO implement me
	panic("implement me")
}
