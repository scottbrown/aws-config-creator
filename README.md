# aws-config-creator

Command line tool to automatically generate a `.aws/config` file based on AWS SSO permission sets in your org.

When working in an organization that uses AWS SSO, there are often multiple
permission sets (like IAM roles) that can be assumed by people.  This tool
looks through all permission sets and account assignments and generates
an `.aws/config` file containing these roles that can be assumed.

## Permissions Required

This tool requires some readonly permissions from your AWS organization account.  They are:

1. `organizations:ListAccounts`
1. `sso:ListInstances`
1. `sso:ListPermissionSetsProvisionedToAccount`
1. `sso:DescribePermissionSet`

## Usage

```bash
$ aws-config-creator --sso-session td --sso-region us-east-1 --profile admin
Writing 43 entries to aws.config...done.
```

Add `-v` or `--verbose` to see more output about which operations are
happening.

## Contributing

1. Fork the repository.
1. Make your change.
1. `make fmt`
1. `make test`
1. `make build`
1. Make a Pull Request.

## License

[MIT](LICENSE)
