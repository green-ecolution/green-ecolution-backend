package plugin

// Token represents an authentication token used to interact with a plugin host.
//
// Fields:
// - AccessToken: A short-lived token used for authenticating API requests to the plugin host.
// - RefreshToken: A long-lived token used to obtain a new access token when the current one expires.
//
// The `Token` structure is typically used after a plugin registers with the plugin host and
// is issued authentication credentials.
//
// Example usage:
//
//	token := &Token{
//		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//		RefreshToken: "dGVzdF9yZWZyZXNoX3Rva2VuX3ZhbHVl...",
//	}
type Token struct {
	AccessToken  string
	RefreshToken string
}
