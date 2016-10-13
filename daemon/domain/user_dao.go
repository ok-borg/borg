package domain

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (ud *UserDao) GetById(id string) (User, error) {
	u := User{Id: id}
	err := ud.db.First(&u).Error
	return u, err
}

func (ud *UserDao) GetByIds(ids []string) ([]User, error) {
	us := []User{}
	err := ud.db.Where("users.id in (?)", ids).
		Find(&us).Error
	return us, err
}

func (ud *UserDao) GetByGithubUserId(id string) (User, error) {
	u := User{GithubId: id}
	err := ud.db.First(&u).Error
	return u, err
}

func (ud *UserDao) GetOrCreateFromRaw(username, email, githubId string) (User, error) {
	// try to get the user
	u, err := ud.GetByGithubUserId(githubId)
	if err != nil {
		// user do not exist creating new one
		newUser := User{
			Id:        uuid.NewV4().String(),
			Username:  username,
			Email:     email,
			GithubId:  githubId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return newUser, ud.Create(newUser)
	}

	return u, nil
}

func (ud *UserDao) GetByEmailOrUsername(str string) (User, error) {
	u := User{}
	err := ud.db.Where("users.email = ? OR users.username = ?", str, str).
		First(&u).Error
	return u, err
}

func (ud *UserDao) GetByEmail(email string) (User, error) {
	u := User{Email: email}
	err := ud.db.First(&u).Error
	return u, err
}

func (ud *UserDao) GetByUsername(username string) (User, error) {
	u := User{Username: username}
	err := ud.db.First(&u).Error
	return u, err
}

func (ud *UserDao) Create(u User) error {
	return ud.db.Create(&u).Error
}

func (ud *UserDao) Update(u User) error {
	return ud.db.Save(&u).Error
}
