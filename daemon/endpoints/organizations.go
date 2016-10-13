package endpoints

import (
	"errors"
	"fmt"
	"time"

	log "github.com/cihub/seelog"
	"github.com/jinzhu/gorm"
	"github.com/ok-borg/borg/daemon/domain"
	"github.com/satori/go.uuid"
)

func (e Endpoints) CreateOrganization(db *gorm.DB, userId string, name string) error {
	organizationDao := domain.NewOrganizationDao(db)
	// first check if organization with same name exists
	if _, err := organizationDao.GetByName(name); err == nil {
		// no error, we sucessfully get an organization,
		// return an error
		return errors.New(fmt.Sprintf("An organization with the name %s already exists", name))
	}

	// first create organization
	newOrganization := domain.Organization{
		Id:   uuid.NewV4().String(),
		Name: name,
		// by default, the creator user is admin, in the future we may allow to change the
		// admin user
		UserAdminId: userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   userId,
		UpdatedBy:   userId,
	}

	if err := organizationDao.Create(newOrganization); err != nil {
		log.Errorf("[Endpoints.CreateOrganization] unable to create organization: %s", err.Error())
		return err
	}

	// then create association between organization and the creator user
	userOrganizationDao := domain.NewUserOrganizationDao(db)
	newUserOrganization := domain.UserOrganization{
		Id:             uuid.NewV4().String(),
		UserId:         userId,
		OrganizationId: newOrganization.Id,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CreatedBy:      userId,
		UpdatedBy:      userId,
	}

	if err := userOrganizationDao.Create(newUserOrganization); err != nil {
		log.Errorf("[Endpoints.CreateOrganization] unable to associate user to newly create organization: %s", err.Error())
		return err
	}

	return nil
}
