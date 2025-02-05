package main

const AppName string = "setlist"

const AppDescShort string = "Creates an AWS config file from AWS SSO configuration"

const AppDescLong string = `
Parses an AWS organizations permission set structure to build a complete .aws/config file with all permission sets provisioned across all AWS member accounts
`

const FlagSSOSession string = "sso-session"

const FlagSSORegion string = "sso-region"

const FlagProfile string = "profile"

const FlagMapping string = "mapping"

const FlagOutput string = "output"

const FlagStdout string = "stdout"

const FlagSSOFriendlyName string = "sso-friendly-name"

const DEFAULT_FILENAME string = "aws.config"
