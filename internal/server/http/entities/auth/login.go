package auth

type LoginResponse struct {
	LoginURL string `json:"login_url"`
} // @Name LoginResponse

type LoginTokenRequest struct {
	Code string `json:"code"`
} // @Name LoginTokenRequest
