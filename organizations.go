package awsconfigcreator

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
)

func ListAccounts(ctx context.Context, client *organizations.Client) ([]types.Account, error) {
	var accounts []types.Account

	var token *string
	for {
		orgOutput, err := client.ListAccounts(ctx, &organizations.ListAccountsInput{NextToken: token})
		if err != nil {
			return accounts, err
		}

		for _, v := range orgOutput.Accounts {
			accounts = append(accounts, v)
		}

		if orgOutput.NextToken == nil {
			break
		}

		token = orgOutput.NextToken
	}

	return accounts, nil
}
