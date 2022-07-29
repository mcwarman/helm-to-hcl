# helm-to-hcl

Utility for converting helm configuration to HCL for use in Terraform.

## Usage

```shell
helm-to-hcl -f ./values.yaml -o main.tf
```

## Install

### Windows

```shell
scoop bucket add mcwarman https://github.com/mcwarman/scoop-bucket
scoop install helm-to-hcl
```

### Mac & Linux

```shell
brew install mcwarman/tap/helm-to-hcl
```
