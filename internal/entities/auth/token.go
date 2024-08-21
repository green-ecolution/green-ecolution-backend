package auth

type IntroSpectTokenResult struct {
	Exp      *int
	Active   *bool
	AuthTime *int
	Type     *string
}
