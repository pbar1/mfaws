package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/viper"
)

// CredentialsShortTerm is used to reflect updated credentials
type CredentialsShortTerm struct {
	AssumedRole        string `ini:"assumed_role"`
	AssumedRoleARN     string `ini:"assumed_role_arn,omitempty"`
	ExternalID         string `ini:"external_id,omitempty"`
	AWSAccessKeyID     string `ini:"aws_access_key_id"`
	AWSSecretAccessKey string `ini:"aws_secret_access_key"`
	AWSSessionToken    string `ini:"aws_session_token"`
	AWSSecurityToken   string `ini:"aws_security_token"`
	Expiration         string `ini:"expiration"`
}

// CheckError is a simple wrapper for log.Fatalln()
func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
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
		fmt.Printf("MFA token code: ")
		_, err := fmt.Scanln(&mfaToken)
		CheckError(err)
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
	CheckError(err)
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
		p.RoleSessionName = viper.GetString("role-session-name")
		p.ExternalID = aws.String(viper.GetString("external-id"))
	})

	credsRepsonse, err := creds.Get()
	CheckError(err)
	expirationTime := time.Now().UTC().Add(time.Duration(viper.GetInt("duration")) * time.Second)

	credsShortTerm := CredentialsShortTerm{
		AssumedRole:        "True",
		AssumedRoleARN:     viper.GetString("assume-role"),
		ExternalID:         viper.GetString("external-id"),
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
	if viper.GetBool("verbose") {
		log.Printf("credentials-file: %s\n", viper.Get("credentials-file"))
		log.Printf("profile: %s\n", viper.Get("profile"))
		log.Printf("long-term-suffix: %s\n", viper.Get("long-term-suffix"))
		log.Printf("short-term-suffix: %s\n", viper.Get("short-term-suffix"))
		log.Printf("device: %s\n", viper.Get("device"))
		log.Printf("assume-role: %s\n", viper.Get("assume-role"))
		log.Printf("duration: %d\n", viper.Get("duration"))
		log.Printf("role-session-name: %s\n", viper.Get("role-session-name"))
		log.Printf("force: %t\n", viper.Get("force"))
		log.Printf("verbose: %t\n", viper.Get("verbose"))
		log.Printf("token: %s\n", viper.Get("token"))
		log.Printf("external-id: %s\n", viper.Get("external-id"))
	}
}
