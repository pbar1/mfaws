<h1 align="center" style="border-bottom: none;">:lock: mfaws :lock:</h1>
<h3 align="center">AWS Multi-Factor Authentication manager</h3>

<p align="center">
  <a href="https://travis-ci.org/pbar1/mfaws">
    <img alt="Build Status" src="https://travis-ci.org/pbar1/mfaws.svg?branch=develop">
  </a>
  <a href="https://github.com/pbar1/mfaws/releases/latest">
    <img alt="GitHub release" src="https://img.shields.io/github/release/pbar1/mfaws.svg">
  </a>
  <a href="https://goreportcard.com/report/github.com/pbar1/mfaws">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/pbar1/mfaws">
  </a>
  <a href="https://hub.docker.com/r/pbar1/mfaws">
    <img alt="Docker pulls" src="https://img.shields.io/docker/pulls/pbar1/mfaws.svg">
  </a>
</p>

<p align="center">
  <img src="./assets/example.svg"/>
</p>

<!-- installation -->
## Installation
 
#### Install script (Linux & macOS)
Download the appropriate binary from the releases page, `chmod +x`, and drop it into your `PATH`.

#### [Chocolatey][4] (Windows)
```powershell
choco install mfaws
```

**Note**: Make sure your hardware clock is correct! [Especially if dual booting][7]. If your time is out of sync, your MFA attempts will fail _and_ the codes `oathtool` generates will be wrong (if you use it).
<!-- installationstop -->

<!-- usage -->
## Usage
```
AWS Multi-Factor Authentication manager

Usage:
  mfaws [flags]
  mfaws [command]

Available Commands:
  help        Help about any command
  version     Prints mfaws version information

Flags:
  -a, --assume-role string         ARN of IAM role to assume [MFA_ASSUME_ROLE]
  -c, --credentials-file string    Path to AWS credentials file (default "~/.aws/credentials") [AWS_SHARED_CREDENTIALS_FILE]
  -d, --device string              ARN of MFA device to use [MFA_DEVICE]
  -l, --duration int               Duration in seconds for credentials to remain valid (default assume-role ? 3600 : 43200) [MFA_STS_DURATION]
  -e, --external-id string         Unique ID used by third parties to assume a role in their customers' accounts [AWS_EXTERNAL_ID]
  -f, --force                      Force credentials to refresh even if not expired
  -h, --help                       help for mfaws
      --long-term-suffix string    Suffix appended to long-term profiles (default "-long-term")
  -p, --profile string             Name of profile to use in AWS credentials file (default "default") [AWS_PROFILE]
  -s, --role-session-name string   Session name when assuming a role
      --short-term-suffix string   Suffix appended to short-term profiles (default "")
  -t, --token string               MFA token to use for authentication
  -v, --verbose                    Enable verbose output

Use "mfaws [command] --help" for more information about a command.
```
<!-- usagestop -->

<!-- examples -->
## Examples

#### Using the default profile
Make sure you have the following in your `$HOME/.aws/credentials` file:
*NOTE:* the profile name must be `[default-long-term]`
```
[default-long-term]
aws_access_key_id     = $YOUR_AWS_ACCESS_KEY_ID
aws_secret_access_key = $YOUR_AWS_SECRET_ACCESS_KEY
aws_mfa_device        = $YOUR_MFA_DEVICE_ARN
```

Then, simply run
```sh
mfaws
```
to fetch temporary credentials for your **default** AWS profile. More advanced configuration is possible (see [Usage](#usage)).

#### Combine `mfaws` with [`oathtool`][2]
Set an alias for generating your MFA token, then pipe it into `mfaws`:
```sh
alias otp-aws="oathtool --totp --base32 $YOUR_AWS_TOTP_KEY"

otp-aws | mfaws
# or
otp-aws | mfaws -p some-profile
```
<!-- examplesstop -->


[1]: https://github.com/pbar1/mfaws/releases
[2]: https://www.nongnu.org/oath-toolkit/
[3]: https://github.com/go-semantic-release/semantic-release
[4]: https://chocolatey.org/packages/mfaws
[5]: https://github.com/polygamma/aurman
[6]: https://aur.archlinux.org/packages/mfaws-bin/
[7]: https://wiki.archlinux.org/index.php/Time#UTC_in_Windows
