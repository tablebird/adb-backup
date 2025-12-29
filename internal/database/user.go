package database

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id         int    `gorm:"primaryKey;autoIncrement"`
	Name       string `gorm:"unique;type:varchar(255);not null"`
	Password   string `gorm:"type:varchar(255);not null"`
	SourceType int64  `gorm:"type:int;not null;default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (u *User) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CountLocalUsers() (int64, error) {
	var count int64
	err := db.Model(&User{}).Where("source_type = ?", 0).Count(&count).Error
	return count, err
}

func CreateUser(user *User) error {
	res, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = res
	return db.Create(user).Error
}

func FindUserByName(name string) (*User, error) {
	var user User
	err := db.Where("name = ?", name).First(&user).Error
	return &user, err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
