package auth

import "net/url"

type LoginRequest struct {
	RedirectURL *url.URL
}

type LoginResp struct {
	LoginURL *url.URL
}

type LoginCallback struct {
  Code string `validate:"required"`
  RedirectURL *url.URL 
}

