package domain

import "github.com/jinzhu/gorm"

type OrganizationDao struct {
	db *gorm.DB
}

func NewOrganizationDao(db *gorm.DB) *OrganizationDao {
	return &OrganizationDao{db: db}
}

func (od *OrganizationDao) GetById(id string) (Organization, error) {
	u := Organization{Id: id}
	err := od.db.First(&u).Error
	return u, err
}

func (od *OrganizationDao) GetByIds(ids []string) ([]Organization, error) {
	u := []Organization{}
	err := od.db.Where("organizations.id in (?)", ids).
		Find(&u).Error
	return u, err
}

func (od *OrganizationDao) GetByName(name string) (Organization, error) {
	u := Organization{}
	err := od.db.Where("organizations.name = ?", name).
		First(&u).Error
	return u, err
}

func (od *OrganizationDao) Create(u Organization) error {
	return od.db.Create(&u).Error
}

func (od *OrganizationDao) Update(u Organization) error {
	return od.db.Save(&u).Error
}
