package setlist

// Profile represents an AWS SSO profile configuration.
type Profile struct {
	Name            string
	Description     string
	SessionDuration string
	SessionName     string
	AccountId       string
	RoleName        string
}
