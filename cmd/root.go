package cmd

import (
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ini "gopkg.in/ini.v1"
)

var log = logging.MustGetLogger("logger")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

// CredentialsShortTerm is used to reflect updated credentials
type CredentialsShortTerm struct {
	AssumedRole        string `ini:"assumed_role"`
	AssumedRoleARN     string `ini:"assumed_role_arn,omitempty"`
	AWSAccessKeyID     string `ini:"aws_access_key_id"`
	AWSSecretAccessKey string `ini:"aws_secret_access_key"`
	AWSSessionToken    string `ini:"aws_session_token"`
	AWSSecurityToken   string `ini:"aws_security_token"`
	Expiration         string `ini:"expiration"`
}

var rootCmd = &cobra.Command{
	Use:   "mfaws",
	Short: "AWS Multi-Factor Authentication manager",
	Long:  `AWS Multi-Factor Authentication manager`,

	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := ini.Load(viper.GetString("credentials-file"))
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
			os.Exit(1)
		}
		profileLongTerm := viper.GetString("profile") + viper.GetString("long-term-suffix")
		profileShortTerm := viper.GetString("profile") + viper.GetString("short-term-suffix")

		if cfg.Section(profileShortTerm).HasKey("expiration") && !viper.GetBool("force") {
			expirationUnparsed := cfg.Section(profileShortTerm).Key("expiration").String()
			expiration, _ := time.Parse("2006-01-02 15:04:05", expirationUnparsed)
			secondsRemaining := expiration.Unix() - time.Now().Unix()
			if secondsRemaining > 0 {
				log.Infof("Your credentials are still valid for %d seconds", secondsRemaining)
				os.Exit(1)
			}
		}

		if cfg.Section(profileLongTerm).HasKey("aws_mfa_device") {
			viper.SetDefault("device", cfg.Section(profileLongTerm).Key("aws_mfa_device").String())
		}
		if cfg.Section(profileLongTerm).HasKey("assume_role") {
			viper.SetDefault("assume-role", cfg.Section(profileLongTerm).Key("assume_role").String())
		}

		sess := CreateSession(profileLongTerm)
		var credsShortTerm CredentialsShortTerm
		if len(viper.GetString("assume-role")) == 0 {
			viper.SetDefault("duration", 43200)
			DumpConfig()
			credsShortTerm = GetCredsWithoutRole(sess)
		} else {
			viper.SetDefault("duration", 3600)
			DumpConfig()
			credsShortTerm = GetCredsWithRole(sess)
		}
		log.Debug(credsShortTerm)
		err = cfg.Section(profileShortTerm).ReflectFrom(&credsShortTerm)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		cfg.SaveTo(viper.GetString("credentials-file"))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute: %v", err)
		os.Exit(1)
	}
}

func init() {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
		os.Exit(1)
	}

	rootCmd.PersistentFlags().StringP("credentials-file", "c", "", "Path to the AWS credentials file [AWS_SHARED_CREDENTIALS_FILE]")
	rootCmd.PersistentFlags().StringP("profile", "p", "", "Name of the CLI profile to use [AWS_PROFILE]")
	rootCmd.PersistentFlags().String("long-term-suffix", "", "Suffix appended to long-term profiles")
	rootCmd.PersistentFlags().String("short-term-suffix", "", "Suffix appended to short-term profiles")
	rootCmd.PersistentFlags().StringP("device", "d", "", "MFA Device ARN [MFA_DEVICE]")
	rootCmd.PersistentFlags().StringP("assume-role", "a", "", "ARN of the IAM role to assume [MFA_ASSUME_ROLE]")
	rootCmd.PersistentFlags().IntP("duration", "l", 0, "Duration in seconds for the credentials to remain valid [MFA_STS_DURATION]")
	rootCmd.PersistentFlags().StringP("role-session-name", "s", "", "Session name")
	rootCmd.PersistentFlags().BoolP("force", "f", false, "Refresh credentials even if currently valid")
	rootCmd.PersistentFlags().String("log-level", "false", "Set log level")
	rootCmd.PersistentFlags().Bool("setup", false, "Setup a new long term credentials section")
	rootCmd.PersistentFlags().StringP("token", "t", "", "Provide MFA token as an argument")

	viper.BindPFlags(rootCmd.PersistentFlags())

	viper.BindEnv("credentials-file", "AWS_SHARED_CREDENTIALS_FILE")
	viper.BindEnv("profile", "AWS_PROFILE")
	viper.BindEnv("device", "MFA_DEVICE")
	viper.BindEnv("assume-role", "MFA_ASSUME_ROLE")
	viper.BindEnv("duration", "MFA_STS_DURATION")

	viper.SetDefault("credentials-file", filepath.Join(currentUser.HomeDir, ".aws", "credentials"))
	viper.SetDefault("profile", "default")
	viper.SetDefault("long-term-suffix", "-long-term")
	viper.SetDefault("short-term-suffix", "")
	viper.SetDefault("role-session-name", currentUser.Username)
}

// CreateSession creates an AWS session from the given profile
func CreateSession(profileLongTerm string) *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: profileLongTerm,
	}))
	return sess
}

// GetMFAToken retrieves MFA token codes from either stdin or the "token" flag
func GetMFAToken() string {
	var mfaToken string
	if viper.GetString("token") == "" || viper.GetString("token") == "-" {
		mfaToken, _ = stscreds.StdinTokenProvider()
	} else {
		mfaToken = viper.GetString("token")
	}
	return mfaToken
}

// GetCredsWithoutRole is used to get temporary AWS credentials when NOT assuming a role
func GetCredsWithoutRole(sess *session.Session) CredentialsShortTerm {

	mfaToken := GetMFAToken()

	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(viper.GetInt64("duration")),
		SerialNumber:    aws.String(viper.GetString("device")),
		TokenCode:       aws.String(mfaToken),
	}
	svc := sts.New(sess)
	result, err := svc.GetSessionToken(input)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	creds := result.Credentials
	credsShortTerm := CredentialsShortTerm{
		AssumedRole:        "False",
		AssumedRoleARN:     "",
		AWSAccessKeyID:     *creds.AccessKeyId,
		AWSSecretAccessKey: *creds.SecretAccessKey,
		AWSSessionToken:    *creds.SessionToken,
		AWSSecurityToken:   *creds.SessionToken,
		Expiration:         creds.Expiration.Format("2006-01-02 15:04:05"),
	}
	return credsShortTerm
}

// GetCredsWithRole is used to get temporary AWS credentials when assuming a role
func GetCredsWithRole(sess *session.Session) CredentialsShortTerm {

	mfaToken := GetMFAToken()

	creds := stscreds.NewCredentials(sess, viper.GetString("assume-role"), func(p *stscreds.AssumeRoleProvider) {
		p.Duration = time.Duration(viper.GetInt("duration")) * time.Second
		p.SerialNumber = aws.String(viper.GetString("device"))
		p.TokenCode = aws.String(mfaToken)
	})
	credsRepsonse, err := creds.Get()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	expirationTime := time.Now().UTC().Add(time.Duration(viper.GetInt("duration")) * time.Second)
	credsShortTerm := CredentialsShortTerm{
		AssumedRole:        "True",
		AssumedRoleARN:     viper.GetString("assume-role"),
		AWSAccessKeyID:     credsRepsonse.AccessKeyID,
		AWSSecretAccessKey: credsRepsonse.SecretAccessKey,
		AWSSessionToken:    credsRepsonse.SessionToken,
		AWSSecurityToken:   credsRepsonse.SessionToken,
		Expiration:         expirationTime.Format("2006-01-02 15:04:05"),
	}
	return credsShortTerm
}

// DumpConfig logs the current viper configuration for debugging
func DumpConfig() {
	log.Debugf("credentials-file: %s", viper.Get("credentials-file"))
	log.Debugf("profile: %s", viper.Get("profile"))
	log.Debugf("long-term-suffix: %s", viper.Get("long-term-suffix"))
	log.Debugf("short-term-suffix: %s", viper.Get("short-term-suffix"))
	log.Debugf("device: %s", viper.Get("device"))
	log.Debugf("assume-role: %s", viper.Get("assume-role"))
	log.Debugf("duration: %d", viper.Get("duration"))
	log.Debugf("role-session-name: %s", viper.Get("role-session-name"))
	log.Debugf("force: %t", viper.Get("force"))
	log.Debugf("log-level: %s", viper.Get("log-level"))
	log.Debugf("setup: %t", viper.Get("setup"))
	log.Debugf("token: %s", viper.Get("token"))
}
