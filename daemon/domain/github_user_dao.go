package domain

import "github.com/jinzhu/gorm"

type GithubUserDao struct {
	db *gorm.DB
}

func NewGithubUserDao(db *gorm.DB) *GithubUserDao {
	return &GithubUserDao{db: db}
}

func (gu *GithubUserDao) Create(model GithubUser) error {
	return gu.db.Create(&model).Error
}

func (gu *GithubUserDao) GetByGithubId(id string) (GithubUser, error) {
	model := GithubUser{}
	err := gu.db.Where("github_users.github_id = ?", id).
		First(&model).Error
	return model, err
}
