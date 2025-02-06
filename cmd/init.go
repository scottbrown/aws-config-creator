package main

// init initializes command-line flags for the root command.
func init() {
	rootCmd.PersistentFlags().StringVarP(&ssoSession, FlagSSOSession, "s", "", "Nickname to give the SSO Session (e.g. org name) (required)")
	rootCmd.PersistentFlags().StringVarP(&profile, FlagProfile, "p", "", "Profile")
	rootCmd.PersistentFlags().StringVarP(&ssoRegion, FlagSSORegion, "r", "", "AWS region where AWS SSO resides (required)")
	rootCmd.PersistentFlags().StringVarP(&mapping, FlagMapping, "m", "", "Comma-delimited Account Nickname Mapping (id=nickname)")
	rootCmd.PersistentFlags().StringVarP(&filename, FlagOutput, "o", DEFAULT_FILENAME, "Where the AWS config file will be written")
	rootCmd.PersistentFlags().BoolVar(&stdout, FlagStdout, false, "Specify this flag to write the config file to stdout instead of a file")
	rootCmd.PersistentFlags().StringVar(&ssoFriendlyName, FlagSSOFriendlyName, "", "Use this instead of the identity store ID for the start URL")

	rootCmd.MarkPersistentFlagRequired(FlagSSOSession)
	rootCmd.MarkPersistentFlagRequired(FlagSSORegion)
}
