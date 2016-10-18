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

func (e Endpoints) CreateOrganization(
	db *gorm.DB,
	userId string,
	name string,
) (*domain.Organization, error) {
	organizationDao := domain.NewOrganizationDao(db)
	// first check if organization with same name exists
	if _, err := organizationDao.GetByName(name); err == nil {
		// no error, we sucessfully get an organization,
		// return an error
		return nil, errors.New(fmt.Sprintf("An organization with the name %s already exists", name))
	}

	// first create organization
	newOrganization := domain.Organization{
		Id:   uuid.NewV4().String(),
		Name: name,
		// by default, the creator user is admin, in the future we may allow to change the
		// admin user
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: userId,
		UpdatedBy: userId,
	}

	if err := organizationDao.Create(newOrganization); err != nil {
		log.Errorf("[Endpoints.CreateOrganization] unable to create organization: %s", err.Error())
		return nil, err
	}

	// then create association between organization and the creator user
	userOrganizationDao := domain.NewUserOrganizationDao(db)
	newUserOrganization := domain.UserOrganization{
		Id:             uuid.NewV4().String(),
		UserId:         userId,
		OrganizationId: newOrganization.Id,
		IsAdmin:        1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CreatedBy:      userId,
		UpdatedBy:      userId,
	}

	if err := userOrganizationDao.Create(newUserOrganization); err != nil {
		log.Errorf("[Endpoints.CreateOrganization] unable to associate user to newly create organization: %s", err.Error())
		return nil, err
	}

	return &newOrganization, nil
}

func (e Endpoints) CreateOrganizationJoinLink(
	db *gorm.DB,
	userId string,
	organizationId string,
	ttl int64,
) (*domain.OrganizationJoinLink, error) {
	// Get UserOrganization to check if the user is admin or not
	userOrganizationDao := domain.NewUserOrganizationDao(db)
	userOrganization, err := userOrganizationDao.GetByUserAndOrganization(userId, organizationId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"user (id=%s) is not member of the organization (id=%s)",
			userId, organizationId))
	}
	if userOrganization.IsAdmin != 1 {
		return nil, errors.New(fmt.Sprintf(
			"user (id=%s) is not administrator of organization (id=%s)",
			userId, organizationId))
	}

	// ok so here the organization exists, the user is the admin
	// lets check if a join-link already exist, if yes remove it then create a new one hehe.
	organizationJoinLinkDao := domain.NewOrganizationJoinLinkDao(db)
	if ojl, err := organizationJoinLinkDao.GetByOrganizationId(organizationId); err == nil {
		// no error this join link exist for this organization
		// lets remove it
		organizationJoinLinkDao.Delete(ojl.Id)
	}
	// now create a new one
	ojl := domain.OrganizationJoinLink{
		Id:             uuid.NewV4().String(),
		OrganizationId: organizationId,
		Ttl:            ttl,
		CreatedAt:      time.Now(),
		CreatedBy:      userId,
	}

	if err := organizationJoinLinkDao.Create(ojl); err != nil {
		log.Errorf("[Endpoints.CreateOrganizationjoinlink] unable to create organiation join link for organization: %s, %s", organizationId, err.Error())
		return nil, err
	}
	return &ojl, nil
}

func (e Endpoints) DeleteOrganizationJoinLink(
	db *gorm.DB,
	userId string,
	organizationJoinLinkId string,
) error {
	// get the organizastionJoinLink
	organizationJoinLinkDao := domain.NewOrganizationJoinLinkDao(db)
	ojl, err := organizationJoinLinkDao.GetById(organizationJoinLinkId)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot fin organization join link (id=%s)",
			organizationJoinLinkId))
	}

	// then get UserOrganization to check if the user is admin or not
	userOrganizationDao := domain.NewUserOrganizationDao(db)
	userOrganization, err := userOrganizationDao.GetByUserAndOrganization(userId, ojl.OrganizationId)
	if err != nil {
		return errors.New(fmt.Sprintf(
			"user (id=%s) is not member of the organization (id=%s)",
			userId, ojl.OrganizationId))
	}
	if userOrganization.IsAdmin != 1 {
		return errors.New(fmt.Sprintf(
			"user (id=%s) is not administrator of organization (id=%s)",
			userId, ojl.OrganizationId))
	}

	// ok so here the organization exists, the user is the admin
	// DELETE ALL THE SHIT
	return organizationJoinLinkDao.Delete(organizationJoinLinkId)
}

func (e Endpoints) GetOrganizationJoinLink(
	db *gorm.DB,
	organizationJoinLinkId string,
) (domain.OrganizationJoinLink, error) {
	organizationJoinLinkDao := domain.NewOrganizationJoinLinkDao(db)
	return organizationJoinLinkDao.GetById(organizationJoinLinkId)
}

func (e Endpoints) GetOrganizationJoinLinkForOrganization(
	db *gorm.DB,
	userId string,
	organizationId string,
) (*domain.OrganizationJoinLink, error) {
	organizationJoinLinkDao := domain.NewOrganizationJoinLinkDao(db)
	organizationJoinLink, err := organizationJoinLinkDao.GetByOrganizationId(organizationId)

	if err != nil {
		log.Errorf("[Endpoints.GetOrganizationjoinlinkbyorganizationid] cannot get join link from organization: %s", err)
		return nil, err
	}

	// then get UserOrganization to check if the user is admin or not
	userOrganizationDao := domain.NewUserOrganizationDao(db)
	userOrganization, err := userOrganizationDao.GetByUserAndOrganization(userId, organizationJoinLink.OrganizationId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"user (id=%s) is not member of the organization (id=%s)",
			userId, organizationJoinLink.OrganizationId))
	}
	if userOrganization.IsAdmin != 1 {
		return nil, errors.New(fmt.Sprintf(
			"user (id=%s) is not administrator of organization (id=%s)",
			userId, organizationJoinLink.OrganizationId))
	}

	return &organizationJoinLink, err
}

func (e Endpoints) ListUserOrganizations(db *gorm.DB, userId string) ([]domain.Organization, error) {
	organizationIds, err := domain.NewUserOrganizationDao(db).ListOrganizationsForUser(userId)
	if err != nil {
		log.Errorf("[Endpoint.ListUserOrganizations]cannot list organizations for user %s", userId)
		return nil, errors.New("cannot read organizations")
	}
	return domain.NewOrganizationDao(db).GetByIds(organizationIds)
}

func (e Endpoints) JoinOrganization(
	db *gorm.DB,
	userId string,
	organizationJoinLinkId string,
) error {
	// get the organizationJoinLink
	organizationJoinLinkDao := domain.NewOrganizationJoinLinkDao(db)
	ojl, err := organizationJoinLinkDao.GetById(organizationJoinLinkId)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot fin organization join link (id=%s)",
			organizationJoinLinkId))
	}

	// if not expired continue
	if ojl.IsExpired() {
		return errors.New("join link expired")
	}

	userOrganizationDao := domain.NewUserOrganizationDao(db)
	// if already member returnn error
	if _, err := userOrganizationDao.GetByUserAndOrganization(userId, ojl.OrganizationId); err == nil {
		return errors.New("you already joined this organization")
	}

	userOrganization := domain.UserOrganization{
		Id:             uuid.NewV4().String(),
		UserId:         userId,
		OrganizationId: ojl.OrganizationId,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CreatedBy:      userId,
		UpdatedBy:      userId,
	}

	return userOrganizationDao.Create(userOrganization)
}

func (e Endpoints) LeaveOrganization(
	db *gorm.DB,
	userId string,
	organizationId string,
) error {

	userOrganizationDao := domain.NewUserOrganizationDao(db)
	if userOrganization, err := userOrganizationDao.GetByUserAndOrganization(userId, organizationId); err != nil {
		// user is not part of this organization
		return errors.New(fmt.Sprintf("[Endpoints.LeaveOrganization] User (id=%s), is not part of organization (id=%s)", userId, organizationId))
	} else {
		if userOrganization.IsAdmin == 1 {
			admins, _ := userOrganizationDao.GetAdmins(organizationId)
			// check if you are last admin, is so you can leave only if you
			// are the last user of the org`
			if len(admins) == 1 {
				users, _ := userOrganizationDao.ListUsersInOrganization(organizationId)
				if len(users) != 1 {
					// you are not the last user, you cannot leave this shit
					// lets transfer ownership before
					return errors.New(fmt.Sprintf("[Endpoints.LeaveOrganization] User (id=%s), is the last admin for the organization, he can leave only if there is 1 member in the organization (id=%s)",
						userId, organizationId))
				}
			}
		}
		return userOrganizationDao.Delete(userOrganization.Id)
	}
}

func (e Endpoints) ExpelUserFromOrganization(
	db *gorm.DB,
	userId string,
	userIdToExpel string,
	organizationId string,
) error {
	// first check if the user is admin
	userOrganizationDao := domain.NewUserOrganizationDao(db)
	adminOjl, err := userOrganizationDao.GetByUserAndOrganization(userId, organizationId)
	if err != nil {
		return errors.New(fmt.Sprintf(
			"user (id=%s) is not member of the organization (id=%s)",
			userId, organizationId))
	}
	if adminOjl.IsAdmin != 1 {
		return errors.New(fmt.Sprintf(
			"user (id=%s) is not administrator of organization (id=%s)",
			userId, organizationId))
	}

	// then just remove the user organization assorciation
	userOrganization, err := userOrganizationDao.GetByUserAndOrganization(userIdToExpel, organizationId)
	if err != nil {
		return errors.New(fmt.Sprintf("[Endpoints.ExpelUserFromOrganization] user (id=%s) is not part of organization (id=%s)", userIdToExpel, organizationId))
	}

	return userOrganizationDao.Delete(userOrganization.Id)
}

func (e Endpoints) GrantAdminRightToUser(
	db *gorm.DB,
	userId string,
	userIdToAdmin string,
	organizationId string,
) error {
	// first check if the user is admin
	userOrganizationDao := domain.NewUserOrganizationDao(db)
	adminOjl, err := userOrganizationDao.GetByUserAndOrganization(userId, organizationId)
	if err != nil {
		return errors.New(fmt.Sprintf(
			"user (id=%s) is not member of the organization (id=%s)",
			userId, organizationId))
	}
	if adminOjl.IsAdmin != 1 {
		return errors.New(fmt.Sprintf(
			"user (id=%s) is not administrator of organization (id=%s)",
			userId, organizationId))
	}

	// then just get  the user organization assorciation and update it
	userOrganization, err := userOrganizationDao.GetByUserAndOrganization(userIdToAdmin, organizationId)
	if err != nil {
		return errors.New(fmt.Sprintf("[Endpoints.ExpelUserFromOrganization] user (id=%s) is not part of organization (id=%s)", userIdToAdmin, organizationId))
	}
	userOrganization.IsAdmin = 1
	return userOrganizationDao.Update(userOrganization)
}
