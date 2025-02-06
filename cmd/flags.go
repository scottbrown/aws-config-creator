package main

var (
	ssoSession      string // SSO session nickname
	profile         string // AWS profile name
	ssoRegion       string // AWS region
	mapping         string // Mapping of account IDs to nicknames
	filename        string // Output filename
	stdout          bool   // Flag to print output to stdout instead of a file
	ssoFriendlyName string // Optional friendly name for the SSO instance
)
