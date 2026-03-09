package database

import (
	"adb-backup/internal/auth"

	"fmt"

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

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) updatePassword(password string) error {
	pass, err := hashPassword(password)
	if err != nil {
		return err
	}
	err = db.Model(u).Update("password", pass).Error
	if err == nil {
		u.Password = pass
	}
	return err
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

func Authenticate(authSource int, username, password string) (*User, error) {
	if authSource <= 0 {
		user, err := FindUserByName(username)
		if user.SourceType != 0 {
			return nil, fmt.Errorf("user not match auth source")
		}
		if err != nil || !user.CheckPassword(password) {
			return nil, fmt.Errorf("password error")
		}
		return user, nil
	}
	authType := auth.Type(authSource)
	provider, ok := auth.Providers[authType]
	if !ok {
		return nil, fmt.Errorf("not found auth provider")
	}
	ok, err := provider.Authenticate(username, password)
	if err != nil || !ok {
		return nil, fmt.Errorf("provider password error")
	}
	return createUser(username, password, int64(authSource))
}

func createUser(username, password string, sourceType int64) (*User, error) {
	user, err := FindUserByName(username)
	if err != nil {
		return nil, err
	}
	if user != nil && user.Name == username {
		if user.SourceType != sourceType {
			return nil, fmt.Errorf("user already exists %v, type : %d", user, sourceType)
		}
		if !user.CheckPassword(password) {
			err := user.updatePassword(password)
			if err != nil {
				return nil, err
			}
		}
		return user, nil
	}
	user = &User{
		Name:       username,
		Password:   password,
		SourceType: sourceType,
	}
	err = CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
