# helm-to-hcl

Utility for converting helm configuration to HCL for use in Terraform.

## Usage

```shell
helm-to-hcl -f ./test/values.yaml -o main.tf
```

```terraform
# main.tf
resource "helm_release" "default" {
  values = [
    yamlencode(
      {
        image = {
          registry = "docker.io"
          repository = "bitnami/postgresql"
          tag = "14.4.0-debian-11-r13"
        }
      }
    )
  ]
}
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
