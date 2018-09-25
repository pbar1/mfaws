<h1 align="center" style="border-bottom: none;">:lock: mfaws :lock:</h1>
<h3 align="center">AWS Multi-Factor Authentication manager</h3>

<p align="center">
  <a href="https://travis-ci.org/pbar1/mfaws">
    <img alt="Build Status" src="https://travis-ci.org/pbar1/mfaws.svg?branch=master">
  </a>
  <a href="https://github.com/pbar1/mfaws/releases/latest">
    <img alt="GitHub release" src="https://img.shields.io/github/release/pbar1/mfaws.svg">
  </a>
  <a href="https://goreportcard.com/report/github.com/pbar1/mfaws">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/pbar1/mfaws">
  </a>
</p>

<p align="center">
  <a href="https://asciinema.org/a/194262" target="_blank">
    <img width="443" height="261" src="https://asciinema.org/a/194262.png"/>
  </a>
</p>

<!-- installation -->
## Installation
 
#### Install script (Linux & macOS)
```sh
curl -sL --proto-redir -all,https https://raw.githubusercontent.com/pbar1/mfaws/master/install.sh | sh
```

#### [Chocolatey][4] (Windows)
```powershell
choco install mfaws
```

#### [AUR][6] (Arch Linux)
```sh
git clone https://aur.archlinux.org/mfaws-bin.git
cd mfaws-bin
makepkg -si
```
Or, if you have an AUR helper like [aurman][5],
```sh
aurman -S mfaws-bin
```

#### Brew (macOS & Linux)
_coming soon!_
<!-- installationstop -->

<!-- usage -->
## Usage
```
AWS Multi-Factor Authentication manager

Usage:
  mfaws [flags]

Flags:
  -a, --assume-role string         ARN of IAM role to assume [MFA_ASSUME_ROLE]
  -c, --credentials-file string    Path to AWS credentials file (default "~/.aws/credentials") [AWS_SHARED_CREDENTIALS_FILE]
  -d, --device string              ARN of MFA device to use [MFA_DEVICE]
  -l, --duration int               Duration in seconds for credentials to remain valid (default assume-role ? 3600 : 43200) [MFA_STS_DURATION]
  -f, --force                      Force credentials to refresh even if not expired
  -h, --help                       help for mfaws
      --long-term-suffix string    Suffix appended to long-term profiles (default "-long-term")
  -p, --profile string             Name of profile to use in AWS credentials file (default "default") [AWS_PROFILE]
  -s, --role-session-name string   Session name when assuming a role
      --short-term-suffix string   Suffix appended to short-term profiles (default "")
  -t, --token string               MFA token to use for authentication
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

Then, simply run
```sh
mfaws
```
to fetch temporary credentials for your **default** AWS profile. More advanced configuration is possible (see [Usage](#usage)).

#### Combine `mfaws` with [`oathtool`][2]
Set an alias for generating your MFA token, then pipe it into `mfaws`:
```sh
alias otp-aws="oathtool --totp --base32 $YOUR_AWS_TOTP_KEY"

otp-aws | mfaws -t -
```
<!-- examplesstop -->

<!-- todo -->
## Todo
Subcommands:
- [ ] `setup`, to configure long term profiles
- [ ] `check` time left on short term creds

Continuous integration
- [ ] Deploy to Homebrew

Other:
- [ ] Documentation
- [ ] Testing
- [ ] Debug and error logging
<!-- todostop -->

[1]: https://github.com/pbar1/mfaws/releases
[2]: https://www.nongnu.org/oath-toolkit/
[3]: https://github.com/go-semantic-release/semantic-release
[4]: https://chocolatey.org/packages/mfaws
[5]: https://github.com/polygamma/aurman
[6]: https://aur.archlinux.org/packages/mfaws-bin/
