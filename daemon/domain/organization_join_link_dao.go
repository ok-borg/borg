package domain

import "github.com/jinzhu/gorm"

type OrganizationJoinLinkDao struct {
	db *gorm.DB
}

func NewOrganizationJoinLinkDao(db *gorm.DB) *OrganizationJoinLinkDao {
	return &OrganizationJoinLinkDao{db: db}
}

func (od *OrganizationJoinLinkDao) GetById(id string) (OrganizationJoinLink, error) {
	u := OrganizationJoinLink{Id: id}
	err := od.db.First(&u).Error
	return u, err
}

func (od *OrganizationJoinLinkDao) GetByOrganizationId(id string) (OrganizationJoinLink, error) {
	o := OrganizationJoinLink{}
	err := od.db.Where("organization_id = ?", id).
		First(&o).Error
	return o, err
}

func (od *OrganizationJoinLinkDao) Create(u OrganizationJoinLink) error {
	return od.db.Create(&u).Error
}

func (od *OrganizationJoinLinkDao) Update(u OrganizationJoinLink) error {
	return od.db.Save(&u).Error
}

func (od *OrganizationJoinLinkDao) Delete(id string) error {
	return od.db.Delete(&OrganizationJoinLink{Id: id}).Error
}
