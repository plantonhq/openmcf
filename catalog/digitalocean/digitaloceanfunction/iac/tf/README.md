# DigitalOcean Function -- Terraform Module

Deploys a `digitalocean_app` with a single functions component from a `DigitalOceanFunction` spec. There is no standalone Functions resource. Provider pin is `~> 2.99`.

`variables.tf` is generated (`planton tofu generate-variables DigitalOceanFunction`). Do not hand-edit it. The API token lives in `credentials.tf`.

Other component families on `digitalocean_app` (service, worker, job, static site, database) are not used — they belong on DigitalOceanApp.

## Prerequisites

- OpenTofu or Terraform 1.0+
- DigitalOcean API token (`digitalocean_token`)

## Usage

```hcl
module "function" {
  source = "./path/to/module"

  metadata = {
    name = "hello"
  }

  spec = {
    function_name    = "hello"
    region           = "nyc3"
    source_directory = "packages"
    git = {
      repo_clone_url = "https://github.com/digitalocean/sample-functions-nodejs-helloworld.git"
      branch         = "master"
    }
  }

  digitalocean_token = var.digitalocean_token
}

output "https_endpoint" {
  value = module.function.https_endpoint
}
```

## Outputs

| Output | Description |
|--------|-------------|
| `function_id` | App UUID (import id for `digitalocean_app`) |
| `https_endpoint` | Public HTTPS URL |
| `default_hostname` | Default `ondigitalocean.app` hostname |

Runtime, memory, timeout, and schedules come from `project.yml` in `source_directory`, not from this module. See the kind [GUIDE](../../GUIDE.md).
