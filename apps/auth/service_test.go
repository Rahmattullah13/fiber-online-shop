package auth

import (
	"context"
	"fiber-online-shop/external/database"
	"fiber-online-shop/infra/response"
	"fiber-online-shop/internal/config"
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var svc service

func init() {
	filename := "../../cmd/api/config.yaml"
	err := config.LoadConfig(filename)
	if err != nil {
		panic(err)
	}

	db, err := database.ConnectPostgres(config.Cfg.DB)
	if err != nil {
		panic(err)
	}

	repo := newRepository(db)
	svc = newService(repo)
}

func TestRegister(t *testing.T) {
	req := RegisterRequestPayload{
		Email:    fmt.Sprintf("%v@gmail.com", uuid.NewString()),
		Password: "secret123",
	}
	err := svc.register(context.Background(), req)
	require.Nil(t, err)
}

func TestRegister_Fail(t *testing.T) {
	t.Run("error email already used", func(t *testing.T) {
		//preparation for duplicate email
		email := fmt.Sprintf("%v@gmail.com", uuid.NewString())
		req := RegisterRequestPayload{
			Email:    email,
			Password: "secret123",
		}
		err := svc.register(context.Background(), req)
		require.Nil(t, err)
		// end preparation

		err = svc.register(context.Background(), req)
		require.NotNil(t, err)
		require.Equal(t, response.ErrEmailAlreadyUsed, err)
	})
}

func TestLogin_Success(t *testing.T) {
	email := fmt.Sprintf("%v@gmail.com", uuid.NewString())
	password := "secret123"
	request := RegisterRequestPayload{
		Email:    email,
		Password: password,
	}
	err := svc.register(context.Background(), request)
	require.Nil(t, err)

	reqLogin := LoginRequestPayload{
		Email:    email,
		Password: password,
	}

	token, err := svc.login(context.Background(), reqLogin)
	require.Nil(t, err)
	require.NotEmpty(t, token)
	log.Println(token)
}
