package setlist

import (
	"fmt"
)

// ConfigFile represents the structure of the configuration file,
// including session details, profiles, and nickname mappings.
type ConfigFile struct {
	SessionName     string            // Name of the SSO session
	IdentityStoreId string            // The unique identity store ID
	FriendlyName    string            // Alt name used for the SSO instance
	Region          string            // AWS region
	Profiles        []Profile         // List of AWS profiles
	NicknameMapping map[string]string // Mapping of account IDs to nicknames
}

// StartURL constructs the AWS SSO start URL based on the IdentityStoreId
// or FriendlyName.
func (c *ConfigFile) StartURL() string {
	subdomain := c.IdentityStoreId

	if c.hasFriendlyName() {
		subdomain = c.FriendlyName
	}

	return fmt.Sprintf("https://%s.awsapps.com/start", subdomain)
}

// hasFriendlyName checks if a friendly name has been set for the SSO
// instance.
func (c *ConfigFile) hasFriendlyName() bool {
	return c.FriendlyName != ""
}

// HasNickname determines whether an account has a mapped nickname.
func (c ConfigFile) HasNickname(accountId string) bool {
	exists := true

	if _, ok := c.NicknameMapping[accountId]; !ok {
		exists = false
	}

	return exists
}
