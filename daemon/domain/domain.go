package domain

import "time"

const (
	AccountTypeGithub = "GITHUB"
)

type User struct {
	Id          string
	Login       string
	Name        string
	Email       string
	AvatarUrl   string
	AccountType string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GithubUser struct {
	Id         string
	GithubId   string
	BorgUserId string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type AccessToken struct {
	Id        string
	Token     string
	UserId    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Organization struct {
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string
}

type UserOrganization struct {
	Id             string
	UserId         string
	OrganizationId string
	IsAdmin        int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      string
	UpdatedBy      string
}

type OrganizationJoinLink struct {
	Id             string
	OrganizationId string
	Ttl            int64
	CreatedAt      time.Time
	CreatedBy      string
}

func (o OrganizationJoinLink) IsExpired() bool {
	if o.CreatedAt.Unix()+o.Ttl < time.Now().Unix() {
		return true
	}
	return false
}
