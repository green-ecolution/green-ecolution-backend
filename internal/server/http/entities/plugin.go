package entities

type PluginResponse struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	HostPath    string `json:"host_path"`
}

type PluginRegisterRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Path        string `json:"path"`
	Auth        struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"auth"`
}

type PluginRegisterResponse struct {
	Success bool                `json:"success"`
	Token   ClientTokenResponse `json:"token"`
}

type PluginListResponse struct {
	Plugins []PluginResponse `json:"plugins"`
}
