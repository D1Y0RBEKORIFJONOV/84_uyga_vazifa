package userusecase

import (
	"context"
	userentity "cors/internal/entity/user"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, req *userentity.CreateUser) (*userentity.Status, error)
	VeryFyUser(ctx context.Context, req *userentity.VerifyRequest) (*userentity.Status, error)
	LoginUser(ctx context.Context, req *userentity.LoginRequest) (*userentity.Token, error)
}

type User struct {
	user UserUseCase
}

func NewUserUseCase(user UserUseCase) *User {
	return &User{
		user: user,
	}
}

func (u *User) CreateUser(ctx context.Context, req *userentity.CreateUser) (*userentity.Status, error) {
	return u.user.CreateUser(ctx, req)
}
func (u *User) VeryFyUser(ctx context.Context, req *userentity.VerifyRequest) (*userentity.Status, error) {
	return u.user.VeryFyUser(ctx, req)
}

func (u *User) LoginUser(ctx context.Context, req *userentity.LoginRequest) (*userentity.Token, error) {
	return u.user.LoginUser(ctx, req)
}
