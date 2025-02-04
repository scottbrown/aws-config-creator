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

var (
	ssoSession      string
	profile         string
	ssoRegion       string
	mapping         string
	filename        string
	stdout          bool
	ssoFriendlyName string
)

func handleRoot(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()

	if err := cmd.MarkFlagRequired(FlagSSOSession); err != nil {
		return err
	}
	if err := cmd.MarkFlagRequired(FlagSSORegion); err != nil {
		return err
	}

	var cfg aws.Config
	var err error

	// check if a profile is specified
	if profile != "" {
		// create a config with the specified profile
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(ssoRegion), config.WithSharedConfigProfile(profile))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		// create a config with static credentials
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(ssoRegion))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	ssoClient := ssoadmin.NewFromConfig(cfg)
	orgClient := organizations.NewFromConfig(cfg)

	instance, err := core.SsoInstance(ctx, ssoClient)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	accounts, err := core.ListAccounts(ctx, orgClient)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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

var rootCmd = &cobra.Command{
	Use:   "setlist",
	Short: "Creates an AWS config file from AWS SSO configuration",
	Long:  "Parses an AWS organizations permission set structure to build a complete .aws/config file with all permission sets provisioned across all AWS member accounts",
	RunE:  handleRoot,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&ssoSession, FlagSSOSession, "s", "", "Nickname to give the SSO Session (e.g. org name)")
	rootCmd.PersistentFlags().StringVarP(&profile, FlagProfile, "p", "", "Profile")
	rootCmd.PersistentFlags().StringVarP(&ssoRegion, FlagSSORegion, "r", "", "AWS region where AWS SSO resides")
	rootCmd.PersistentFlags().StringVarP(&mapping, FlagMapping, "m", "", "Comma-delimited Account Nickname Mapping (id=nickname)")
	rootCmd.PersistentFlags().StringVarP(&filename, FlagOutput, "o", DEFAULT_FILENAME, "Where the AWS config file will be written")
	rootCmd.PersistentFlags().BoolVar(&stdout, FlagStdout, false, "Specify this flag to write the config file to stdout instead of a file")
	rootCmd.PersistentFlags().StringVar(&ssoFriendlyName, FlagSSOFriendlyName, "", "Use this instead of the identity store ID for the start URL")
}
