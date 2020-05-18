package api

import (
	"fmt"
	"path"

	"gitlab-metrics/gitlab/dto"
	"gitlab-metrics/log"
)

const (
	usersPath = "/users"
	userPath  = "/users/%d"
)

type UsersApi interface {
	GetAllUsers() ([]dto.User, error)
	GetAllUsersWithOpts(opts UserOpts) ([]dto.User, error)
	GetAllUsersChanneled(users chan<- dto.User) (uint32, error)

	GetUser(userId uint32) (*dto.User, error)
	GetUserWithOpts(userId uint32, opts UserOpts) (*dto.User, error)
}

type usersApi struct {
	api
	log.Loggable
}

func (p *usersApi) GetAllUsersWithOpts(opts UserOpts) ([]dto.User, error) {
	urlPath := path.Join(apiPath, usersPath)
	urlValues := opts.ToValues()
	reqUrl := fmt.Sprintf("%s%s?%s", p.GetBaseUrl(), urlPath, urlValues.Encode())

	var users []dto.User

	err := p.FetchData(reqUrl, &users)
	if err != nil {
		return nil, fmt.Errorf("fetch data: %v", err)
	}

	return users, nil
}

func (p *usersApi) GetAllUsers() ([]dto.User, error) {
	return p.GetAllUsersWithOpts(UserOpts{})
}

func (p *usersApi) GetAllUsersChanneled(fetchedUsers chan<- dto.User) (uint32, error) {
	fetchedCount := uint32(0)
	pageNumber := uint32(1)

	for {
		users, err := p.GetAllUsersWithOpts(UserOpts{
			RequestOpts: RequestOpts{
				PerPage: 100,
				Page:    pageNumber,
			}})

		if err != nil {
			return fetchedCount, err
		}

		p.Log().Tracef("fetched %d users on page %d", len(users), pageNumber)
		if len(users) == 0 {
			return fetchedCount, nil
		}

		for _, user := range users {
			fetchedUsers <- user
			fetchedCount++
		}

		pageNumber++
	}
}

func (p *usersApi) GetUserWithOpts(userId uint32, opts UserOpts) (*dto.User, error) {
	urlPath := fmt.Sprintf(path.Join(apiPath, userPath), userId)
	urlValues := opts.ToValues()
	reqUrl := fmt.Sprintf("%s%s?%s", p.GetBaseUrl(), urlPath, urlValues.Encode())

	user := &dto.User{}

	err := p.FetchData(reqUrl, user)
	if err != nil {
		return nil, fmt.Errorf("fetch data: %v", err)
	}

	return user, nil
}

func (p *usersApi) GetUser(userId uint32) (*dto.User, error) {
	return p.GetUserWithOpts(userId, UserOpts{})
}
