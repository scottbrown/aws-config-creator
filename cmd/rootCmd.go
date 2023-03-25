package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/go-ini/ini"
	"github.com/spf13/cobra"
)

const DEFAULT_FILENAME string = "aws.config"

var (
	ssoSession string
	profile    string
	ssoRegion  string
	mapping    string
	filename   string
	stdout     bool
)

func parseNicknameMapping(mapping string) map[string]string {
	nicknameMapping := make(map[string]string)

	if len(mapping) == 0 {
		return nicknameMapping
	}

	tokens := strings.Split(mapping, ",")
	for _, token := range tokens {
		parts := strings.Split(token, "=")

		nicknameMapping[parts[0]] = parts[1]
	}

	return nicknameMapping
}

func handleRoot(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()

	if err := cmd.MarkFlagRequired("sso-session"); err != nil {
		return err
	}
	if err := cmd.MarkFlagRequired("sso-region"); err != nil {
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

	nicknameMapping := parseNicknameMapping(mapping)

	// get the SSO instance ARN (there's only one allowed)
	ssoClient := ssoadmin.NewFromConfig(cfg)

	resp, err := ssoClient.ListInstances(ctx, nil)
	if err != nil {
		return err
	}

	if len(resp.Instances) == 0 {
		return errors.New("SSO is not enabled.  No SSO instances exist.")
	}

	instanceArn := resp.Instances[0].InstanceArn

	// list all accounts
	orgClient := organizations.NewFromConfig(cfg)

	orgOutput, err := orgClient.ListAccounts(ctx, &organizations.ListAccountsInput{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// print out the SSO session configuration to file
	payload := ini.Empty()

	ssoSection := payload.Section(fmt.Sprintf("sso-session %s", ssoSession))
	if _, err := ssoSection.NewKey("sso_start_url", fmt.Sprintf("https://%s.awsapps.com/start", *resp.Instances[0].IdentityStoreId)); err != nil {
		return err
	}
	if _, err := ssoSection.NewKey("sso_region", ssoRegion); err != nil {
		return err
	}
	if _, err := ssoSection.NewKey("sso_registration_scopes", "sso:account:access"); err != nil {
		return err
	}

	// loop through each account
	for _, account := range orgOutput.Accounts {
		// list-permission-sets-provisioned-to-account
		params := &ssoadmin.ListPermissionSetsProvisionedToAccountInput{
			InstanceArn: instanceArn,
			AccountId:   account.Id,
		}
		resp, err := ssoClient.ListPermissionSetsProvisionedToAccount(ctx, params)
		if err != nil {
			return err
		}

		// loop through permissions sets
		for _, permissionSetArn := range resp.PermissionSets {
			// get permission set name
			params := &ssoadmin.DescribePermissionSetInput{
				InstanceArn:      instanceArn,
				PermissionSetArn: aws.String(permissionSetArn),
			}
			resp, err := ssoClient.DescribePermissionSet(ctx, params)
			if err != nil {
				return err
			}

			// create section for AccountId-PermissionSet profile
			section1 := payload.Section(fmt.Sprintf("profile %s-%s", *account.Id, *resp.PermissionSet.Name))
			section1.Comment = fmt.Sprintf("# %s. Session Duration: %s", *resp.PermissionSet.Description, *resp.PermissionSet.SessionDuration)
			if _, err := section1.NewKey("sso_session", ssoSession); err != nil {
				return err
			}
			if _, err := section1.NewKey("sso_account_id", *account.Id); err != nil {
				return err
			}
			if _, err := section1.NewKey("sso_role_name", *resp.PermissionSet.Name); err != nil {
				return err
			}

			if len(nicknameMapping) == 0 {
				continue
			}

			// create section for AccountNickname-PermissionSet profile
			section2 := payload.Section(fmt.Sprintf("profile %s-%s", nicknameFor(*account.Id, nicknameMapping), *resp.PermissionSet.Name))
			section2.Comment = fmt.Sprintf("# %s. Session Duration: %s", *resp.PermissionSet.Description, *resp.PermissionSet.SessionDuration)
			if _, err := section2.NewKey("sso_session", ssoSession); err != nil {
				return err
			}
			if _, err := section2.NewKey("sso_account_id", *account.Id); err != nil {
				return err
			}
			if _, err := section2.NewKey("sso_role_name", *resp.PermissionSet.Name); err != nil {
				return err
			}
		}
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

func nicknameFor(accountId string, nicknameMapping map[string]string) string {
	v := nicknameMapping[accountId]

	if v == "" {
		v = fmt.Sprintf("NoNickname-%s", accountId)
	}

	return v
}

var rootCmd = &cobra.Command{
	Use:   "aws-config-creator",
	Short: "Creates an AWS config file from AWS SSO configuration",
	Long:  "Parses an AWS organizations permission set structure to build a complete .aws/config file with all permission sets provisioned across all AWS member accounts",
	RunE:  handleRoot,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&ssoSession, "sso-session", "s", "", "Nickname to give the SSO Session (e.g. org name)")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "Profile")
	rootCmd.PersistentFlags().StringVarP(&ssoRegion, "sso-region", "r", "", "AWS region where AWS SSO resides")
	rootCmd.PersistentFlags().StringVarP(&mapping, "mapping", "m", "", "Comma-delimited Account Nickname Mapping (id=nickname)")
	rootCmd.PersistentFlags().StringVarP(&filename, "output", "o", DEFAULT_FILENAME, "Where the AWS config file will be written")
	rootCmd.PersistentFlags().BoolVar(&stdout, "stdout", false, "Specify this flag to write the config file to stdout instead of a file")
}
