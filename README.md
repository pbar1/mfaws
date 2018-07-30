# mfaws

AWS Multi-Factor Authentication manager. Drop-in replacement for [`aws-mfa`][1].

## Usage
```
AWS Multi-Factor Authentication manager

Usage:
  mfaws [flags]

Flags:
  -a, --assume-role string         ARN of the IAM role to assume [MFA_ASSUME_ROLE]
  -c, --credentials-file string    Path to the AWS credentials file [AWS_SHARED_CREDENTIALS_FILE]
  -d, --device string              MFA Device ARN [MFA_DEVICE]
  -l, --duration int               Duration in seconds for the credentials to remain valid [MFA_STS_DURATION]
  -f, --force                      Refresh credentials even if currently valid
  -h, --help                       help for mfaws
      --log-level string           Set log level (default "false")
      --long-term-suffix string    Suffix appended to long-term profiles
  -p, --profile string             Name of the CLI profile to use [AWS_PROFILE]
  -s, --role-session-name string   Session name
      --setup                      Setup a new long term credentials section
      --short-term-suffix string   Suffix appended to short-term profiles
  -t, --token string               Provide MFA token as an argument
```

## Examples
Combine `mfaws` with [`oathtool`][2] for super speed
```sh
alias otp-aws="oathtool --topt --base32 $YOUR_AWS_TOTP_KEY"

otp-aws | mfaws -t -
```

## Todo
Flags:
- [x] `--assume-role`
- [x] `--credentials-file`
- [x] `--device`
- [x] `--duration`
- [x] `--force`
- [x] `--help`
- [x] ~~`--log-level`~~ `--verbose`
- [x] `--long-term-suffix`
- [x] `--profile`
- [x] `--role-session-name`
- [x] ~~`--setup`~~ arguably unnecessary, may become `setup`
- [x] `--short-term-suffix`
- [x] `--token`
- [ ] `--check` or `check` time left on short term creds

Other:
- [ ] Testing
- [ ] Documentation
- [ ] CICD

[1]: https://github.com/broamski/aws-mfa
[2]: https://www.nongnu.org/oath-toolkit/
