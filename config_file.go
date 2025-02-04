package awsconfigcreator

import (
	"fmt"
)

type ConfigFile struct {
	SessionName     string
	IdentityStoreId string
	FriendlyName    string
	Region          string
	Profiles        []Profile
	NicknameMapping map[string]string
}

func (c *ConfigFile) StartURL() string {
	subdomain := c.IdentityStoreId

	if c.hasFriendlyName() {
		subdomain = c.FriendlyName
	}

	return fmt.Sprintf("https://%s.awsapps.com/start", subdomain)
}

func (c *ConfigFile) hasFriendlyName() bool {
	return c.FriendlyName != ""
}

func (c ConfigFile) HasNickname(accountId string) bool {
	exists := true

	if _, ok := c.NicknameMapping[accountId]; !ok {
		exists = false
	}

	return exists
}
