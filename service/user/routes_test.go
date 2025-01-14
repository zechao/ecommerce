package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/zechao158/ecomm/service/user"
	"github.com/zechao158/ecomm/types"
	"github.com/zechao158/ecomm/types/mocks"
)

func TestUserServiceHandlers(t *testing.T) {
	mockUserRep := &mocks.MockUserRepository{
		CreateUserFunc: func(ctx context.Context, user *types.User) error {
			return nil
		},
		GetUserByEmailFunc: func(ctx context.Context, email string) (*types.User, error) {
			return nil, nil
		},
	}
	mockUserRep.GetUserByEmailFunc(context.Background(), "a")

	handler := user.NewHandler(mockUserRep)

	t.Run("fail by sser invalid payload", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "user",
			LastName:  "123",
			Email:     "",
			Password:  "asd",
		}
		data, err := json.Marshal(payload)
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(data))

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.HandleRegister)
		router.ServeHTTP(rr, req)
		fmt.Println(rr.Body.String())
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
