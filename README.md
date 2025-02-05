![setlist](icon.png)

# SetList (originally aws-config-creator)

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

Add `-v` or `--verbose` to see more output about which operations are
happening.

### Create a basic `.aws/config` file

```bash
$ setlist --sso-session acme \
          --sso-region us-east-1 \
          --profile admin
Writing 43 entries to aws.config...done.
```

The resulting file will contain a set of profiles in the format: `[profile AWS_ACCOUNT_ID:PERMISSION_SET_NAME]`

For example: `[profile 0123456789012-AdministratorAccess]`

### Create a friendly `.aws/config` file

```bash
$ setlist --sso-session acme \
          --sso-region us-east-1 \
          --profile admin \
          --mapping "0123456789012=acme,98765432101=acmelite"
Writing 86 entries to aws.config...done.
```

By supplying a `--mapping` flag with a comma-delimited list of key=value pairs corresponding to AWS Account ID and its nickname, the tool will create the basic `.aws/config` profiles and then create a separate set of profiles that follow the format `[profile NICKNAME-PERMISSIONSETNAME]`.  For example: `[profile acme-AdministratorAccess]`.  This removes the need for your users to remember the 12-digit AWS Account ID, but also allows for backward-compatibility for those people that like using the AWS Account ID in the profile name.

## Contributing

1. Fork the repository.
1. Make your change.
1. `task fmt`
1. `task test`
1. `task build`
1. Make a Pull Request.

## Releases

Each release comes with a software bill of materials (SBOM).  It is
generated using [CycloneDX-GoMod](https://github.com/CycloneDX/cyclonedx-gomod) using the following command:

```bash
cyclonedx-gomod mod -licenses -json -output bom.json
```

Releases are typically automated via Github Actions whenever a new tag is
pushed to the default branch.

## Roadmap

- Create a Lambda function artifact that can be run on a schedule, outputting the latest `.aws/config` to an S3 bucket so that it is always available, especially to those users without the permissions to run this tool.

## License

[MIT](LICENSE)
