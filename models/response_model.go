package models

type ResponseJson struct {
	Success      bool        `json:"success"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data"`
	AccessToken  *string     `json:"access_token,omitempty"`
	RefreshToken *string     `json:"refresh_token,omitempty"`
	Total        *int        `json:"total,omitempty"`
}
