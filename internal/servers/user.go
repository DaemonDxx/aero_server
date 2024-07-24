package servers

import (
	"context"
	"errors"
	"fmt"
	userpb "github.com/daemondxx/lks_back/gen/pb/go/user"
	"github.com/daemondxx/lks_back/internal/dao"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct {
	uServ UserService
	log   *zerolog.Logger
	userpb.UnimplementedUserServiceServer
}

func NewUserServer(uServ UserService, log *zerolog.Logger) *UserServer {
	l := log.
		With().
		Str("grpc_server_name", "user_server").
		Logger()

	return &UserServer{
		uServ: uServ,
		log:   &l,
	}
}

func (u *UserServer) GetUserInfo(ctx context.Context, r *userpb.GetUserInfoRequest) (*userpb.UserInfo, error) {
	log := u.getLogger("get_user_info").With().Uint("user_id", uint(r.UserId)).Logger()
	log.Debug().Msg("get user info start...")

	if uInfo, err := u.uServ.GetUserByID(ctx, uint(r.UserId)); err != nil {
		if errors.Is(err, dao.ErrUserNotFound) {
			log.Warn().Msg("user not found")
			st := status.New(codes.NotFound, "user not found")
			return nil, st.Err()
		} else {
			log.Err(err).Msg("get user info error")
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

func (u *UserServer) ChangeUserStatus(ctx context.Context, r *userpb.ChangeStatusRequest) (*emptypb.Empty, error) {
	log := u.getLogger("change_user_status").With().Uint("user_id", uint(r.UserId)).Logger()
	log.Debug().Msg("change user status start...")

	if err := u.uServ.UpdateActiveStatus(ctx, uint(r.UserId), r.ActiveStatus); err != nil {
		if errors.Is(err, dao.ErrUserNotFound) {
			log.Warn().Msg("user not found")
			st := status.New(codes.InvalidArgument, fmt.Sprintf("user with id %d not found", r.UserId))
			return nil, st.Err()
		} else {
			log.Err(err).Msg("change user status error")
			return nil, ErrInternal
		}
	}
	return &emptypb.Empty{}, nil
}

func (u *UserServer) UpdateAccord(ctx context.Context, r *userpb.UpdateRequest) (*emptypb.Empty, error) {
	log := u.getLogger("update_accord").With().Uint("user_id", uint(r.UserId)).Logger()
	log.Debug().Msg("update accord info start...")

	if err := u.uServ.UpdateAccord(ctx, uint(r.UserId), r.Login, r.Password); err != nil {
		if errors.Is(err, dao.ErrUserNotFound) {
			log.Warn().Msg("user not found")
			st := status.New(codes.InvalidArgument, fmt.Sprintf("user with id %d not found", r.UserId))
			return nil, st.Err()
		} else {
			log.Err(err).Msg("update accord info error")
			return nil, ErrInternal
		}
	}
	return &emptypb.Empty{}, nil
}

func (u *UserServer) UpdateLks(ctx context.Context, r *userpb.UpdateRequest) (*emptypb.Empty, error) {
	log := u.getLogger("update_lks").With().Uint("user_id", uint(r.UserId)).Logger()
	log.Debug().Msg("update lks info start...")

	if err := u.uServ.UpdateLKS(ctx, uint(r.UserId), r.Login, r.Password); err != nil {
		if errors.Is(err, dao.ErrUserNotFound) {
			log.Warn().Msg("user not found")
			st := status.New(codes.InvalidArgument, fmt.Sprintf("user with id %d not found", r.UserId))
			return nil, st.Err()
		} else {
			log.Err(err).Msg("update lks info error")
			return nil, ErrInternal
		}
	}
	return &emptypb.Empty{}, nil
}

func (u *UserServer) mustEmbedUnimplementedUserServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (u *UserServer) getLogger(method string) zerolog.Logger {
	return u.log.With().Str("method", method).Logger()
}
