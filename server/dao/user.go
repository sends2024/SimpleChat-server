package dao

import (
	"simplechat/server/common/pkg/db"
	"simplechat/server/models"
)

func CreateUser(user *models.User) {
	db.DB.Create(user)
}

func GetUserByUsername(username string) *models.User {
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil
	}
	return &user
}

func GetUserByID(id string) *models.User {
	var user models.User
	err := db.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil
	}
	return &user
}

func UpdatePassword(userID, newHashedPassword string) {
	db.DB.Model(&models.User{}).Where("id = ?", userID).Update("password", newHashedPassword)
}

func UpdateAvatar(url, userID string) {
	db.DB.Model(&models.User{}).Where("id = ?", userID).Update("avatar_url", url)
}

func DeleteAIHistory(userID string) error {
	return db.DB.Where("user_id = ?", userID).Delete(&models.AIMessage{}).Error
}

func CreateAIMessage(msg *models.AIMessage) error {
	return db.DB.Create(msg).Error
}

func GetAIHistoryList(userID string) ([]models.AIMessage, error) {
	var list []models.AIMessage

	err := db.DB.
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&list).Error

	return list, err
}
