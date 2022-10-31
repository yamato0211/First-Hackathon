package cruds

import (
	"context"
	"first-hackathon/db"
	"first-hackathon/graph/model"
	"first-hackathon/middlewares"
	"first-hackathon/utils"
	"strings"

	"github.com/google/uuid"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

func UserCreate(ctx context.Context, input model.NewUser) (*model.User, error) {
	input.Password = utils.HashPassword(input.Password)

	user := model.User{
		ID:       uuid.New().String(),
		Name:     input.Name,
		Email:    strings.ToLower(input.Email),
		Password: input.Password,
	}

	if err := db.Psql.Model(user).Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func UserGetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := db.Psql.Model(user).Where("id = ?", id).Take(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func UserGetMe(ctx context.Context) (*model.User, error) {
	var user model.User
	uid := middlewares.CtxValue(ctx).ID
	if err := db.Psql.Model(user).Where("id = ?", uid).Take(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAllUser(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	if err := db.Psql.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func UserGetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := db.Psql.Model(user).Where("email LIKE ?", email).Take(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func UserRegister(ctx context.Context, input model.NewUser) (interface{}, error) {
	// Check Email
	_, err := UserGetByEmail(ctx, input.Email)
	if err == nil {
		// if err != record not found
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	createdUser, err := UserCreate(ctx, input)
	if err != nil {
		return nil, err
	}

	token, err := utils.JwtGenerate(ctx, createdUser.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
	}, nil
}

func UserLogin(ctx context.Context, email string, password string) (interface{}, error) {
	getUser, err := UserGetByEmail(ctx, email)
	if err != nil {
		// if user not found
		if err == gorm.ErrRecordNotFound {
			return nil, &gqlerror.Error{
				Message: "Email not found",
			}
		}
		return nil, err
	}

	if err := utils.ComparePassword(getUser.Password, password); err != nil {
		return nil, err
	}

	token, err := utils.JwtGenerate(ctx, getUser.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
	}, nil
}
