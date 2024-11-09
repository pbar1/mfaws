<h1 align="center" style="border-bottom: none;">:lock: mfaws :lock:</h1>
<h3 align="center">AWS Multi-Factor Authentication manager</h3>

<p align="center">
  <a href="https://github.com/pbar1/mfaws/actions/workflows/build.yml">
    <img alt="Build Status" src="https://github.com/pbar1/mfaws/actions/workflows/build.yml/badge.svg">
  </a>
  <a href="https://github.com/pbar1/mfaws/releases/latest">
    <img alt="GitHub release" src="https://img.shields.io/github/release/pbar1/mfaws.svg">
  </a>
  <a href="https://goreportcard.com/report/github.com/pbar1/mfaws">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/pbar1/mfaws">
  </a>
</p>

<p align="center">
  <img src="./assets/example.svg"/>
</p>

## Installation

### Manual

Download the appropriate binary for your OS/arch from the [releases][1] page.

### [Docker](https://github.com/pbar1/mfaws/pkgs/container/mfaws)

```sh
docker pull ghcr.io/pbar1/mfaws:latest
```

### [Homebrew](https://github.com/pbar1/homebrew-tap/blob/main/mfaws.rb)

```sh
brew tap pbar1/tap && brew update
brew install mfaws
```

### [Scoop](https://github.com/pbar1/scoop-bucket/blob/master/bucket/mfaws.json)

```sh
scoop bucket add pbar1 https://github.com/pbar1/scoop-bucket
scoop install pbar1/mfaws
```

### [Chocolatey](https://chocolatey.org/packages/mfaws)

```powershell
choco install mfaws
```

### [AUR](https://aur.archlinux.org/packages/mfaws-bin)

```sh
yay -S mfaws-bin
```

## Usage

`mfaws` works by looking for AWS credentials and an MFA device ARN in profiles suffixed with `-long-term`. It uses those credentials as well as a TOTP code supplied by the user to make an `AssumeRole` call. The outcome of this is another set of short-lived credentials scoped to the role session. These short lived credentials are stored in a separate profile in the credentials file without the `-long-term` suffix.

For example, your `~/.aws/credentials` file should look similar to this. Here we are using the profile `default-long-term`:

```ini
[default-long-term]
aws_access_key_id     = $YOUR_AWS_ACCESS_KEY_ID
aws_secret_access_key = $YOUR_AWS_SECRET_ACCESS_KEY
aws_mfa_device        = $YOUR_MFA_DEVICE_ARN
```

Then, simply run the following, and enter the MFA token when prompted:

```sh
$ mfaws
```

If that is sucessful, it will create a another profile in the credentials file called `default` that contains the session-scoped creds:

```diff
 [default-long-term]
 aws_access_key_id     = $YOUR_AWS_ACCESS_KEY_ID
 aws_secret_access_key = $YOUR_AWS_SECRET_ACCESS_KEY
 aws_mfa_device        = $YOUR_MFA_DEVICE_ARN

+[default]
+aws_access_key_id     = ...
+aws_secret_access_key = ...
+aws_session_token     = ...
```

In this example we used `default` because it is what tools such as the AWS SDK and `aws` CLI load by default when no profile is specified. Using other profiles is also like so: `mfaws -p myprofile`, which will result in the following:

```diff
 [myprofile-term]
 aws_access_key_id     = $YOUR_AWS_ACCESS_KEY_ID
 aws_secret_access_key = $YOUR_AWS_SECRET_ACCESS_KEY
 aws_mfa_device        = $YOUR_MFA_DEVICE_ARN

+[myprofile]
+aws_access_key_id     = ...
+aws_secret_access_key = ...
+aws_session_token     = ...
```

### Combine with [`oathtool`][2]

> [!NOTE]
> Make sure your hardware clock is correct, [especially if dual booting][7]. If your time is out of sync, codes generated on your machine will be wrong and your MFA attempts will fail.

Set an alias for generating your MFA token, then pipe it into `mfaws`:

```sh
alias otp-aws="oathtool --totp --base32 $YOUR_AWS_TOTP_KEY"

otp-aws | mfaws
# or
otp-aws | mfaws -p some-profile
```

### Combine with `aws-vault`

TODO

<!-- Sources -->

[1]: https://github.com/pbar1/mfaws/releases
[2]: https://www.nongnu.org/oath-toolkit/
[3]: https://github.com/go-semantic-release/semantic-release
[5]: https://github.com/polygamma/aurman
[6]: https://aur.archlinux.org/packages/mfaws-bin/
[7]: https://wiki.archlinux.org/index.php/Time#UTC_in_Windows
