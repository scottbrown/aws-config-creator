package awsconfigcreator

import (
	"fmt"
	"strings"

	"github.com/go-ini/ini"
)

type FileBuilder struct {
	Config ConfigFile
}

func NewFileBuilder(configFile ConfigFile) FileBuilder {
	return FileBuilder{
		Config: configFile,
	}
}

func (f *FileBuilder) Build() (*ini.File, error) {
	payload := ini.Empty()

	if err := f.addDefaultSection(payload); err != nil {
		return payload, err
	}

	if err := f.addSSOSection(payload); err != nil {
		return payload, err
	}

	for _, p := range f.Config.Profiles {
		p.Name = fmt.Sprintf("%s-%s", p.AccountId, p.RoleName)
		if err := f.addProfileSection(p, payload); err != nil {
			return payload, err
		}

		if !f.Config.HasNickname(p.AccountId) {
			continue
		}

		// create section for AccountNickname-PermissionSet profile
		p.Name = fmt.Sprintf("%s-%s", f.Config.NicknameMapping[p.AccountId], p.RoleName)
		if err := f.addProfileSection(p, payload); err != nil {
			return payload, err
		}
	}

	return payload, nil
}

func (f *FileBuilder) addDefaultSection(file *ini.File) error {
	section := file.Section("default")

	if _, err := section.NewKey(SSOSessionAttrKey, f.Config.SessionName); err != nil {
		return err
	}

	return nil
}

func (f *FileBuilder) addSSOSection(file *ini.File) error {
	section := file.Section(strings.Join([]string{SSOSessionSectionKey, f.Config.SessionName}, " "))

	if _, err := section.NewKey(SSOStartUrlKey, f.Config.StartURL()); err != nil {
		return err
	}

	if _, err := section.NewKey(SSORegionKey, f.Config.Region); err != nil {
		return err
	}

	if _, err := section.NewKey(SSORegistrationScopesKey, SSORegistrationScopesValue); err != nil {
		return err
	}

	return nil
}

func (f *FileBuilder) addProfileSection(p Profile, file *ini.File) error {
	section := file.Section(fmt.Sprintf("profile %s", p.Name))

	section.Comment = fmt.Sprintf("# %s. Session Duration: %s", p.Description, p.SessionDuration)

	if _, err := section.NewKey(SSOSessionAttrKey, p.SessionName); err != nil {
		return err
	}

	if _, err := section.NewKey(SSOAccountIdKey, p.AccountId); err != nil {
		return err
	}

	if _, err := section.NewKey(SSORoleNameKey, p.RoleName); err != nil {
		return err
	}

	return nil
}

func (f *FileBuilder) ssoStartUrl() string {
	subdomain := f.Config.IdentityStoreId

	if f.hasFriendlyName() {
		subdomain = f.Config.FriendlyName
	}

	return fmt.Sprintf("https://%s.awsapps.com/start", subdomain)
}

func (f *FileBuilder) hasFriendlyName() bool {
	return f.Config.FriendlyName != ""
}
