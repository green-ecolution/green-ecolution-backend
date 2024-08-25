package auth

type IntroSpectTokenResult struct {
	Exp      *int
	Active   *bool
	AuthTime *int
	Type     *string
}

type ClientToken struct {
	AccessToken      string
	IDToken          string
	ExpiresIn        int
	RefreshExpiresIn int
	RefreshToken     string
	TokenType        string
	NotBeforePolicy  int
	SessionState     string
	Scope            string
}
