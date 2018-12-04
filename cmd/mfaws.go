package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ini "gopkg.in/ini.v1"

	"github.com/pbar1/mfaws/internal"
)

// nolint: gochecknoglobals
var (
	VERSION string
	COMMIT  string
	DATE    string
)

var rootCmd = &cobra.Command{
	Use:   "mfaws",
	Short: "AWS Multi-Factor Authentication manager",
	Long:  `AWS Multi-Factor Authentication manager`,

	Run: func(cmd *cobra.Command, args []string) {
		userFlow()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version, commit, date string) {
	VERSION = version
	COMMIT = commit
	DATE = date

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	homeDirPath, err := homedir.Dir()
	internal.CheckError(err)

	rootCmd.PersistentFlags().StringP("credentials-file", "c", "", "Path to AWS credentials file (default \"~/.aws/credentials\") [AWS_SHARED_CREDENTIALS_FILE]")
	rootCmd.PersistentFlags().StringP("profile", "p", "", "Name of profile to use in AWS credentials file (default \"default\") [AWS_PROFILE]")
	rootCmd.PersistentFlags().String("long-term-suffix", "", "Suffix appended to long-term profiles (default \"-long-term\")")
	rootCmd.PersistentFlags().String("short-term-suffix", "", "Suffix appended to short-term profiles (default \"\")")
	rootCmd.PersistentFlags().StringP("device", "d", "", "ARN of MFA device to use [MFA_DEVICE]")
	rootCmd.PersistentFlags().StringP("assume-role", "a", "", "ARN of IAM role to assume [MFA_ASSUME_ROLE]")
	rootCmd.PersistentFlags().IntP("duration", "l", 0, "Duration in seconds for credentials to remain valid (default assume-role ? 3600 : 43200) [MFA_STS_DURATION]")
	rootCmd.PersistentFlags().StringP("role-session-name", "s", "", "Session name when assuming a role")
	rootCmd.PersistentFlags().BoolP("force", "f", false, "Force credentials to refresh even if not expired")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringP("token", "t", "", "MFA token to use for authentication")

	viper.BindPFlags(rootCmd.PersistentFlags())

	viper.BindEnv("credentials-file", "AWS_SHARED_CREDENTIALS_FILE")
	viper.BindEnv("profile", "AWS_PROFILE")
	viper.BindEnv("device", "MFA_DEVICE")
	viper.BindEnv("assume-role", "MFA_ASSUME_ROLE")
	viper.BindEnv("duration", "MFA_STS_DURATION")

	viper.SetDefault("credentials-file", filepath.Join(homeDirPath, ".aws", "credentials"))
	viper.SetDefault("profile", "default")
	viper.SetDefault("long-term-suffix", "-long-term")
	viper.SetDefault("short-term-suffix", "")
	viper.SetDefault("role-session-name", "mfaws")
}

func userFlow() {
	ini.PrettyFormat = false
	ini.PrettyEqual = true

	cfg, err := ini.Load(viper.GetString("credentials-file"))
	internal.CheckError(err)

	profileLongTerm := viper.GetString("profile") + viper.GetString("long-term-suffix")
	profileShortTerm := viper.GetString("profile") + viper.GetString("short-term-suffix")

	if cfg.Section(profileShortTerm).HasKey("expiration") && !viper.GetBool("force") {
		expirationUnparsed := cfg.Section(profileShortTerm).Key("expiration").String()
		expiration, _ := time.Parse("2006-01-02 15:04:05", expirationUnparsed)
		secondsRemaining := expiration.Unix() - time.Now().Unix()
		if secondsRemaining > 0 {
			fmt.Printf("Credentials for profile `%s` still valid for %d seconds\n", profileShortTerm, secondsRemaining)
			os.Exit(0)
		}
	}

	if cfg.Section(profileLongTerm).HasKey("aws_mfa_device") {
		viper.SetDefault("device", cfg.Section(profileLongTerm).Key("aws_mfa_device").String())
	}
	if cfg.Section(profileLongTerm).HasKey("assume_role") {
		viper.SetDefault("assume-role", cfg.Section(profileLongTerm).Key("assume_role").String())
	}

	sess := internal.CreateSession(profileLongTerm)
	var credsShortTerm internal.CredentialsShortTerm
	if len(viper.GetString("assume-role")) == 0 {
		viper.SetDefault("duration", 43200)
		internal.DumpConfig()
		credsShortTerm = internal.GetCredsWithoutRole(sess)
	} else {
		viper.SetDefault("duration", 3600)
		internal.DumpConfig()
		credsShortTerm = internal.GetCredsWithRole(sess)
	}

	err = cfg.Section(profileShortTerm).ReflectFrom(&credsShortTerm)
	internal.CheckError(err)

	err = cfg.SaveTo(viper.GetString("credentials-file"))
	internal.CheckError(err)

	fmt.Printf("Success! Credentials for profile `%s` valid for %d seconds\n", profileShortTerm, viper.GetInt("duration"))
}
