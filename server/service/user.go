package service

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"simplechat/server/common/pkg/jwt"
	"simplechat/server/common/response"
	"simplechat/server/common/utils"
	"simplechat/server/dao"
	"simplechat/server/models"
)

func Register(username, password, email string) *response.AppError {
	exist := dao.GetUserByUsername(username)
	if exist != nil {
		return response.NewAppError(http.StatusBadRequest, "Username already exists")
	}

	userID := utils.NewULID()
	hashed, _ := utils.HashPassword(password)

	dao.CreateUser(&models.User{
		ID:        userID,
		Username:  username,
		Password:  hashed,
		AvatarURL: "",
		Email:     email,
	})

	return nil
}

func Login(username, password string) (map[string]interface{}, *response.AppError) {
	user := dao.GetUserByUsername(username)
	if user == nil || !utils.CheckPassword(user.Password, password) {
		return nil, response.NewAppError(http.StatusUnauthorized, "Incorrect username or password")
	}

	token, _ := jwt.GenerateToken(user.ID)

	return map[string]interface{}{
		"avatar_url": user.AvatarURL,
		"token":      token,
	}, nil
}

func ChangePassword(OldPassword, NewPassword, userID string) *response.AppError {
	user := dao.GetUserByID(userID)

	if utils.CheckPassword(OldPassword, user.Password) {
		return response.NewAppError(http.StatusBadRequest, "Old password is incorrect")
	}
	hashed, _ := utils.HashPassword(NewPassword)

	dao.UpdatePassword(userID, hashed)

	return nil
}

func ChangeAvatar(userID, newURL string) error {
	user := dao.GetUserByID(userID)
	oldURL := user.AvatarURL

	u, _ := url.Parse(oldURL)
	filename := filepath.Base(u.Path)
	localPath := filepath.Join("uploads", filename)

	if err := os.Remove(localPath); err != nil && !os.IsNotExist(err) {
		log.Printf("failed to delete old avatar: %v", err)
	}

	dao.UpdateAvatar(newURL, userID)

	return nil
}

func DeleteHistory(userID string) *response.AppError {
	if err := dao.DeleteAIHistory(userID); err != nil {
		return response.NewAppError(http.StatusInternalServerError, "Failed to delete history")
	}
	return nil
}

func NewMessage(role, content, userID string) *response.AppError {
	msg := models.AIMessage{
		ID:         utils.NewULID(),
		UserID:     userID,
		SenderRole: role,
		Content:    content,
	}

	if err := dao.CreateAIMessage(&msg); err != nil {
		return response.NewAppError(http.StatusInternalServerError, "Failed to save message")
	}

	return nil
}

func GetAIHistory(userID string) ([]models.AIMessage, *response.AppError) {
	list, err := dao.GetAIHistoryList(userID)
	if err != nil {
		return nil, response.NewAppError(http.StatusInternalServerError, "Failed to get history")
	}
	return list, nil
}
