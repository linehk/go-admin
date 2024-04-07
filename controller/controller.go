package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/linehk/go-admin/config"
	"github.com/linehk/go-admin/errcode"
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

func userModelToResp(userModel model.AppUser) User {
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

func reqToCreateUserParams(req User) (model.CreateUserParams, error) {
	var createUserParams model.CreateUserParams
	createUserParams.Username = req.Username
	password, err := hash(req.Password)
	if err != nil {
		return model.CreateUserParams{}, err
	}
	createUserParams.Password = password
	if req.Name != nil {
		createUserParams.Name = *req.Name
	}
	if req.Email != nil {
		createUserParams.Email = *req.Email
	}
	if req.Phone != nil {
		createUserParams.Phone = *req.Phone
	}
	if req.Remark != nil {
		createUserParams.Remark = *req.Remark
	}
	if req.Status != nil {
		createUserParams.Status = string(*req.Status)
	} else {
		createUserParams.Status = string(Activated)
	}
	if req.Created != nil {
		err = createUserParams.Created.Scan(*req.Created)
		if err != nil {
			return model.CreateUserParams{}, err
		}
	} else {
		err = createUserParams.Created.Scan(time.Now())
		if err != nil {
			return model.CreateUserParams{}, err
		}
	}
	if req.Updated != nil {
		err = createUserParams.Updated.Scan(*req.Updated)
		if err != nil {
			return model.CreateUserParams{}, err
		}
	} else {
		err = createUserParams.Updated.Scan(time.Now())
		if err != nil {
			return model.CreateUserParams{}, err
		}
	}
	return createUserParams, nil
}

func reqToUpdateUserParams(req User) (model.UpdateUserParams, error) {
	var updateUserParams model.UpdateUserParams
	updateUserParams.Username = req.Username
	password, err := hash(req.Password)
	if err != nil {
		return model.UpdateUserParams{}, err
	}
	updateUserParams.Password = password
	if req.Name != nil {
		updateUserParams.Name = *req.Name
	}
	if req.Email != nil {
		updateUserParams.Email = *req.Email
	}
	if req.Phone != nil {
		updateUserParams.Phone = *req.Phone
	}
	if req.Remark != nil {
		updateUserParams.Remark = *req.Remark
	}
	if req.Status != nil {
		updateUserParams.Status = string(*req.Status)
	}
	if req.Created != nil {
		err = updateUserParams.Created.Scan(*req.Created)
		if err != nil {
			return model.UpdateUserParams{}, err
		}
	}
	if req.Updated != nil {
		err = updateUserParams.Updated.Scan(*req.Updated)
		if err != nil {
			return model.UpdateUserParams{}, err
		}
	} else {
		err = updateUserParams.Updated.Scan(time.Now())
		if err != nil {
			return model.UpdateUserParams{}, err
		}
	}
	return updateUserParams, nil
}

func decode(w http.ResponseWriter, r *http.Request, req any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ReturnErr(w, errcode.Parse)
	}
}

func encode(w http.ResponseWriter, resp any) {
	err := json.NewEncoder(w).Encode(&resp)
	if err != nil {
		ReturnErr(w, errcode.Parse)
	}
}

func paging(current, pageSize int) (int32, int32) {
	if current > 0 && pageSize > 0 {
		return int32((current - 1) * pageSize), int32(pageSize)
	}
	return int32(current), int32(pageSize)
}

func ReturnErr(w http.ResponseWriter, e int32) {
	w.Header().Set("Content-Type", "application/json")
	errResp := Error{
		Code:    e,
		Message: errcode.Msg(e),
	}
	err := json.NewEncoder(w).Encode(errResp)
	if err != nil {
		panic(err)
	}
	slog.Error("err: ", err)
	w.WriteHeader(http.StatusBadRequest)
}
