package dto

type User struct {
	Id        uint32 `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	State     string `json:"state"`
	AvatarUrl string `json:"avatar_url"`
	WebUrl    string `json:"web_url"`
}
