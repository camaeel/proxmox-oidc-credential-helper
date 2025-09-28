# proxmox-oidc-credential-helper

> **Note:** This repository has been moved to https://github.com/proxmox-oidc-credential-helper/proxmox-oidc-credential-helper/


This small utility allows to obtain credentials from Proxmox's UI using OIDC and exports them to the shell as env variables.

## Instalation

### With Homebrew
```shell
brew tap proxmox-oidc-credential-helper/tap
brew install proxmox-oidc-credential-helper --cask
```

### MacOS / Linux - download with curl

```shell
architecture=`uname -m`
if [[ $architecture == "x86_64" ]]; then 
  architecture=amd64; 
elif [[ $architecture == "aarch64" ]]; then 
  architecture=arm64;
fi
LATEST_VERSION=` curl -s https://api.github.com/repos/proxmox-oidc-credential-helper/proxmox-oidc-credential-helper/releases/latest | jq -r '.tag_name'`

wget "https://github.com/proxmox-oidc-credential-helper/proxmox-oidc-credential-helper/releases/download/${LATEST_VERSION}/proxmox-oidc-credential-helper_`uname -o`_${architecture}.tar.gz" -O proxmox-oidc-credential-helper.tar.gz
tar -xzvf proxmox-oidc-credential-helper.tar.gz
```

### Windows

Download binary from releases page & unzip or use Homebrew

## Prerequisites

1. Should be able to properly login to proxmox UI with OIDC provider
2. OIDC provider should have additional allowed callback set to: `http://localhost:8996/oidc/callback` (default can be changed with configuration flags). 

## Usage: 

Download binary and put somewhere in $PATH.

### Shell 

Run `proxmox-oidc-credential-helper -proxmox-url https://proxmox.example.com:8006 -realm OIDC-REALM-NAME`. By default this application will output export commands to be executed to set credentials in the shell. 
Alternatively run:

```shell
eval $(proxmox-oidc-credential-helper -proxmox-url https://proxmox.example.com:8006 -realm OIDC-REALM-NAME)
```

### Terragrunt

For automated use with Terragrunt and [bpg/proxmox terraform provider](https://registry.terraform.io/providers/bpg/proxmox/latest) add following code:
```hcl
terragrunt {
  extra_arguments "proxmox_oidc_helper" {
    commands = [
      "apply",
      "refresh",
      "import",
      "plan",
      "taint",
      "untaint",
      "destroy",
      "state",
      "output",
      "console",
    ]
    env_vars = {
      PROXMOX_VE_AUTH_TICKET = jsondecode(local.proxmox_creds)["data"]["ticket"]
      PROXMOX_VE_CSRF_PREVENTION_TOKEN = jsondecode(local.proxmox_creds)["data"]["CSRFPreventionToken"]
    }
  }
}

locals {
  proxmox_creds = run_cmd("--terragrunt-quiet", "proxmox-oidc-credential-helper","-proxmox-url=https://proxmox.example.com:8006","-realm=REALM_NAME", "-output=json")
}
```

### Proxmoxer

[Proxmoxer](https://proxmoxer.github.io/docs/latest/) can utilize proxmox-oidc-credential-helper to obtain credentials from OIDC using web browser.
Example script can be found in this repo in [this example](examples/proxmoxer/proxmoxer_example.py)


For additional configuration options please check `proxmox-oidc-credential-helper -h`

# Compatibility

This application should work well on all OSes, but due to limited resources it is only tested on MacOS.
