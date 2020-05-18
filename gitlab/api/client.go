package api

type Client interface {
	ProjectsApi
	MergeRequestsApi
	UsersApi
	EventsApi
}

type client struct {
	url   string
	token string
	projectsApi
	mergeRequestsApi
	usersApi
	eventsApi
}

func NewClient(url string, token string) *client {
	api := api{
		url:   url,
		token: token,
	}
	projectsApi := projectsApi{api: api}
	usersApi := usersApi{api: api}
	return &client{
		url:              url,
		token:            token,
		projectsApi:      projectsApi,
		mergeRequestsApi: mergeRequestsApi{api: api},
		usersApi:         usersApi,
		eventsApi: eventsApi{
			api:         api,
			projectsApi: projectsApi,
			usersApi:    usersApi,
		},
	}
}
