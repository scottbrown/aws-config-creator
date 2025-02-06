package main

import (
	"context"
	"fmt"
	"os"

	core "github.com/scottbrown/setlist"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     AppName,
	Short:   AppDescShort,
	Long:    AppDescLong,
	RunE:    handleRoot,
	Version: core.VERSION,
}

// handleRoot executes the main logic of the command-line application.
func handleRoot(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()

	var cfg aws.Config
	var err error

	// check if a profile is specified
	if profile != "" {
		// create a config with the specified profile
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(ssoRegion), config.WithSharedConfigProfile(profile))
		if err != nil {
			return err
		}
	} else {
		// create a config with static credentials
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(ssoRegion))
		if err != nil {
			return err
		}
	}

	ssoClient := ssoadmin.NewFromConfig(cfg)
	orgClient := organizations.NewFromConfig(cfg)

	instance, err := core.SsoInstance(ctx, ssoClient)
	if err != nil {
		return err
	}

	accounts, err := core.ListAccounts(ctx, orgClient)
	if err != nil {
		return err
	}

	configFile := core.ConfigFile{
		SessionName:     ssoSession,
		IdentityStoreId: *instance.IdentityStoreId,
		FriendlyName:    ssoFriendlyName,
		Region:          ssoRegion,
		NicknameMapping: core.ParseNicknameMapping(mapping),
	}

	profiles := []core.Profile{}
	for _, account := range accounts {
		permissionSets, err := core.PermissionSets(ctx, ssoClient, *instance.InstanceArn, *account.Id)

		if err != nil {
			return err
		}

		for _, p := range permissionSets {
			profile := core.Profile{
				Description:     *p.Description,
				SessionDuration: *p.SessionDuration,
				SessionName:     ssoSession,
				AccountId:       *account.Id,
				RoleName:        *p.Name,
			}

			profiles = append(profiles, profile)
		}
	}

	configFile.Profiles = profiles

	builder := core.NewFileBuilder(configFile)
	payload, err := builder.Build()
	if err != nil {
		return err
	}

	if stdout {
		if _, err := payload.WriteTo(os.Stdout); err != nil {
			return err
		}
	} else {
		if err := payload.SaveTo(filename); err != nil {
			return err
		}
		fmt.Printf("Wrote to %s\n", filename)
	}

	return nil
}
