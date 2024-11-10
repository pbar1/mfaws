<h1 align="center" style="border-bottom: none;">:lock: mfaws :lock:</h1>
<p align="center"><b>AWS multi-factor authentication manager</b></p>

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
  <img src="./.github/assets/example.svg"/>
</p>

## Installation

[![Packaging status](https://repology.org/badge/vertical-allrepos/mfaws.svg)](https://repology.org/project/mfaws/versions)

- **Manual:** Download the binary for your OS/arch from the [releases](https://github.com/pbar1/mfaws/releases) page.
- **Nix (flake)**: To run without installing: `nix run github:pbar1/mfaws --`
- **[Go](https://pkg.go.dev/github.com/pbar1/mfaws):** `go install github.com/pbar1/mfaws@latest`
- **[Docker](https://github.com/pbar1/mfaws/pkgs/container/mfaws):** `docker pull ghcr.io/pbar1/mfaws:latest`
- **[Homebrew](https://github.com/pbar1/homebrew-tap/blob/main/mfaws.rb):** `brew tap pbar1/tap && brew install mfaws`
- **[Scoop](https://github.com/pbar1/scoop-bucket/blob/master/bucket/mfaws.json):** `scoop bucket add pbar1 https://github.com/pbar1/scoop-bucket && scoop install pbar1/mfaws`
- **[Chocolatey](https://chocolatey.org/packages/mfaws):** `choco install mfaws`
- **[AUR](https://aur.archlinux.org/packages/mfaws-bin):** `yay -S mfaws-bin`

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
 [myprofile-long-term]
 aws_access_key_id     = $YOUR_AWS_ACCESS_KEY_ID
 aws_secret_access_key = $YOUR_AWS_SECRET_ACCESS_KEY
 aws_mfa_device        = $YOUR_MFA_DEVICE_ARN

+[myprofile]
+aws_access_key_id     = ...
+aws_secret_access_key = ...
+aws_session_token     = ...
```

## Examples

> [!NOTE]
> Make sure your hardware clock is correct, [especially if dual booting](https://wiki.archlinux.org/index.php/Time#UTC_in_Windows). If your time is out of sync, codes generated on your machine will be wrong and your MFA attempts will fail.

### Combine with [`oathtool`](https://www.nongnu.org/oath-toolkit/)

> [!CAUTION]
> While convenient, it's generally not advisable to save the MFA *secret key* to disk, since it does not expire.

You can use `oathtool` to get TOTP codes directly in the CLI without having to copy them from elsewhere. `mfaws` can receive a TOTP code piped from stdin:

```sh
oathtool --totp --base32 $YOUR_AWS_TOTP_KEY | mfaws
```

### Combine with [1Password CLI](https://developer.1password.com/docs/cli/)

You can get TOTP codes from MFA keys that you've saved in your 1Password account. This has the advantage of not leaking the secret to disk. In this example, we're requesting a TOTP code from an item called "AWS" in our 1Password account and piping it into `mfaws`:

```sh
op item get 'AWS' --otp | mfaws
```

### Combine with [HashiCorp Vault](https://developer.hashicorp.com/vault/docs/secrets/totp) TOTP secrets engine

Similar to the above examples, you can request a TOTP code from HashiCorp Vault. In this example, we've enabled the TOTP secret engine and previously saved our MFA secret as an item called `my-aws-totp-secret`. Simply use the Vault CLI to read just the `code` field from that secret: 

```
vault read -field=code totp/code/my-aws-totp-scret | mfaws
```