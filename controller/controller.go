package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/linehk/go-admin/config"
	"github.com/linehk/go-admin/model"
	"golang.org/x/crypto/bcrypt"
)

func Setup() *http.ServeMux {
	mux := http.NewServeMux()
	ctx := context.Background()
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Raw.String("HOST"), config.Raw.String("POSTGRES_USER"), config.Raw.String("POSTGRES_PASSWORD"),
		config.Raw.String("POSTGRES_DB"), config.Raw.String("PORT"), config.Raw.String("SSL_MODE"),
		config.Raw.String("TIMEZONE"))
	userImpl := &UserImpl{DB: model.Setup(ctx, DSN)}
	options := StdHTTPServerOptions{
		BaseRouter: mux,
	}
	HandlerWithOptions(userImpl, options)
	return mux
}

func hash(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(h), err
}

func GetUserRowToResp(userModel model.GetUserRow) User {
	const format = "2006-01-02 15:04:05.999999999"
	var userResp User
	userResp.Username = userModel.Username
	userResp.Name = &userModel.Name
	userResp.Email = &userModel.Email
	userResp.Phone = &userModel.Phone
	userResp.Remark = &userModel.Remark
	userStatus := UserStatus(userModel.Status)
	userResp.Status = &userStatus
	createdStr := userModel.Created.Time.Format(format)
	userResp.Created = &createdStr
	updatedStr := userModel.Updated.Time.Format(format)
	userResp.Updated = &updatedStr
	return userResp
}

func CreateUserRowToResp(userModel model.CreateUserRow) User {
	const format = "2006-01-02 15:04:05.999999999"
	var userResp User
	userResp.Username = userModel.Username
	userResp.Name = &userModel.Name
	userResp.Email = &userModel.Email
	userResp.Phone = &userModel.Phone
	userResp.Remark = &userModel.Remark
	userStatus := UserStatus(userModel.Status)
	userResp.Status = &userStatus
	createdStr := userModel.Created.Time.Format(format)
	userResp.Created = &createdStr
	updatedStr := userModel.Updated.Time.Format(format)
	userResp.Updated = &updatedStr
	return userResp
}
