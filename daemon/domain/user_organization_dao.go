package domain

import (
	"github.com/jinzhu/gorm"
)

type UserOrganizationDao struct {
	db *gorm.DB
}

func NewUserOrganizationDao(db *gorm.DB) *UserOrganizationDao {
	return &UserOrganizationDao{db: db}
}

func (ud *UserOrganizationDao) GetById(id string) (UserOrganization, error) {
	u := UserOrganization{Id: id}
	err := ud.db.First(&u).Error
	return u, err
}

func (ud *UserOrganizationDao) GetByUserAndOrganization(
	user_id string,
	organization_id string,
) (UserOrganization, error) {
	u := UserOrganization{}
	err := ud.db.Where("user_organizations.user_id = ? AND user_organizations.organization_id = ?",
		user_id, organization_id).First(&u).Error
	return u, err
}

// return list of user ids
func (ud *UserOrganizationDao) ListUsersInOrganization(organization_id string) ([]string, error) {
	uos := []UserOrganization{}
	err := ud.db.Where("user_organizations.organization_id = ?", organization_id).
		Find(&uos).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		} else {
			return []string{}, nil
		}
	}
	// get only ids
	uids := []string{}
	for _, uo := range uos {
		uids = append(uids, uo.UserId)
	}

	return uids, nil
}

// return list of user ids
func (ud *UserOrganizationDao) ListOrganizationsForUser(user_id string) ([]string, error) {
	uos := []UserOrganization{}
	err := ud.db.Where("user_organizations.user_id = ?", user_id).
		Find(&uos).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		} else {
			return []string{}, nil
		}
	}
	// get only ids
	uids := []string{}
	for _, uo := range uos {
		uids = append(uids, uo.OrganizationId)
	}

	return uids, nil
}

func (ud *UserOrganizationDao) Create(u UserOrganization) error {
	return ud.db.Create(&u).Error
}

func (ud *UserOrganizationDao) Update(u UserOrganization) error {
	return ud.db.Save(&u).Error
}

func (ud *UserOrganizationDao) Delete(id string) error {
	return ud.db.Delete(&UserOrganization{Id: id}).Error
}

func (ud *UserOrganizationDao) GetAdmins(organizastionId string) ([]UserOrganization, error) {
	models := []UserOrganization{}
	return models, ud.db.Where("user_organizations.organization_id = ?", organizastionId).
		Where("user_organizations.is_admin = 1").
		Find(&models).Error
}
