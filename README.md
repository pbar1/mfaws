# mfaws

AWS Multi-Factor Authentication manager. Drop-in replacement for [`aws-mfa`][1].

[![Build Status](https://travis-ci.org/pbar1/mfaws.svg?branch=master)](https://travis-ci.org/pbar1/mfaws)

<!-- toc -->
* [Installation](#installation)
* [Usage](#usage)
* [Examples](#examples)
* [Todo](#todo)
<!-- tocstop -->

<!-- installation -->
## Installation
Executables for Linux, macOS, and Windows can be found on the [releases](https://github.com/pbar1/mfaws/releases) page.
<!-- installationstop -->

<!-- usage -->
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
      --long-term-suffix string    Suffix appended to long-term profiles
  -p, --profile string             Name of the CLI profile to use [AWS_PROFILE]
  -s, --role-session-name string   Session name
      --short-term-suffix string   Suffix appended to short-term profiles
  -t, --token string               Provide MFA token as an argument
  -v, --verbose                    Enable verbose output
```
<!-- usagestop -->

<!-- examples -->
## Examples

#### Using the default profile
Make sure you have the following in your `$HOME/.aws/credentials` file:
```
[default-long-term]
aws_access_key_id     = $YOUR_AWS_ACCESS_KEY_ID
aws_secret_access_key = $YOUR_AWS_SECRET_ACCESS_KEY
aws_mfa_device        = $YOUR_MFA_DEVICE_ARN
```

Then, simply run `mfaws` to generate/update your **default** AWS profile. More advanced configuration is possible - see [Usage](#usage).
```sh
mfaws
```

#### Combine `mfaws` with [`oathtool`][2] for super speed
```sh
alias otp-aws="oathtool --topt --base32 $YOUR_AWS_TOTP_KEY"

otp-aws | mfaws -t -
```
<!-- examplesstop -->

<!-- todo -->
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
- [x] CICD
<!-- todostop -->

[1]: https://github.com/broamski/aws-mfa
[2]: https://www.nongnu.org/oath-toolkit/

