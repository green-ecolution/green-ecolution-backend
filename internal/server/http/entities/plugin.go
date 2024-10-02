package entities

type PluginRegisterRequest struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Auth struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"auth"`
}

type PluginRegisterResponse struct {
	Success bool                `json:"success"`
	Token   ClientTokenResponse `json:"token"`
}
