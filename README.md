# proxmox-oidc-credential-helper

This small utility allows to obtain credentials from Proxmox's UI using OIDC and exports them to the shell as env variables.

## Prerequisites

1. Should be able to properly login to proxmox UI with OIDC provider
2. OIDC provider should have additional allowed callback set to: `http://localhost:8996/oidc/callback` (default can be changed with configuration flags). 

## Usage: 

Download binary and put somewhere in $PATH.

Run `proxmox-oidc-credentials-helper -proxmox-url https://proxmox.example.com:8006 -realm OIDC-REALM-NAME`. By default this application will output export commands to be executed to set credentials in the shell. 
Alternatively run:

```shell
eval $(proxmox-oidc-credentials-helper -proxmox-url https://proxmox.example.com:8006 -realm OIDC-REALM-NAME)
```

For additional configuration options please check `proxmox-oidc-credentials-helper -h`

# Compatibility

This application should work well on all OSes, but due to limited resources it is only tested on MacOS.
