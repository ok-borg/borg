package domain

import "time"

type User struct {
	Id        string
	Username  string
	Email     string
	GithubId  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Organization struct {
	Id          string
	Name        string
	UserAdminId string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   string
	UpdatedBy   string
}

type UserOrganization struct {
	Id             string
	UserId         string
	OrganizationId string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      string
	UpdatedBy      string
}
