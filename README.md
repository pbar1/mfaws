# mfaws

AWS Multi-Factor Authentication manager. Inspired by and fully compatible with [`aws-mfa`][1].

[1]: https://github.com/broamski/aws-mfa

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
      --force                      Refresh credentials even if currently valid
  -h, --help                       help for mfaws
      --log-level string           Set log level (default "false")
      --long-term-suffix string    Suffix appended to long-term profiles
  -p, --profile string             Name of the CLI profile to use [AWS_PROFILE]
  -s, --role-session-name string   Session name
      --setup                      Setup a new long term credentials section
      --short-term-suffix string   Suffix appended to short-term profiles
  -t, --token int                  Provide MFA token as an argument
```

## Todo
Flags:
- [ ] `--assume-role`
- [x] `--credentials-file`
- [x] `--device`
- [x] `--duration`
- [ ] `--force`
- [ ] `--help` and `help`
- [ ] `--log-level`
- [x] `--long-term-suffix`
- [x] `--profile`
- [ ] `--role-session-name`
- [ ] `--setup`
- [x] `--short-term-suffix`
- [ ] `--token`

Other:
- [ ] Testing
- [ ] Documentation
- [ ] CICD
