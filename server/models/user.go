package models

type User struct {
	ID        string `gorm:"type:char(26);primaryKey" json:"user_id"`
	Username  string `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password  string `gorm:"size:128;not null" json:"-"`
	AvatarURL string `gorm:"size:256" json:"avatar_url"`
	Email     string `gorm:"size:128" json:"email"`
}

type AIMessage struct {
	ID         string `gorm:"type:char(26);primaryKey" json:"id"`
	UserID     string `gorm:"type:char(26);index;not null" json:"user_id"`
	SenderRole string `gorm:"type:varchar(16);not null" json:"sender_role"`
	Content    string `gorm:"type:text;not null" json:"content"`
}
