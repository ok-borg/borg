package domain

import "github.com/jinzhu/gorm"

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

func (ud *UserDao) GetByEmailOrUsername(str string) (User, error) {
	u := User{}
	err := ud.db.Where("users.email = ? OR users.username = ?", str, str).
		First(&u).Error
	return u, err
}

func (ud *UserDao) GetByEmail(email string) (User, error) {
	u := User{}
	err := ud.db.Where("users.email = ?", email).
		First(&u).Error
	return u, err
}

func (ud *UserDao) GetByLogin(login string) (User, error) {
	u := User{}
	err := ud.db.Where("users.login = ?", login).
		First(&u).Error
	return u, err
}

func (ud *UserDao) Create(u User) error {
	return ud.db.Create(&u).Error
}

func (ud *UserDao) Update(u User) error {
	return ud.db.Save(&u).Error
}
