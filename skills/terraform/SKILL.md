---
name: terraform
description: >
  Terraform infrastructure as code patterns. Trigger: When writing .tf files, managing cloud infrastructure, or IaC.
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use
- Writing Terraform configurations (.tf)
- Managing cloud resources (AWS, GCP, Azure)
- Setting up remote state
- Writing modules
- Planning and applying infrastructure changes

## Critical Patterns

### Remote state (REQUIRED)
```hcl
terraform {
  backend "s3" {
    bucket         = "terraform-state"
    key            = "prod/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-locks"
  }
}
```

### Module structure (REQUIRED)
```
modules/
  vpc/
    main.tf
    variables.tf
    outputs.tf
  ecs/
    main.tf
    variables.tf
    outputs.tf
```

### Variables with validation
```hcl
variable "environment" {
  type        = string
  description = "Deployment environment"
  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Must be dev, staging, or prod."
  }
}
```

## Anti-Patterns
### Don't: Hardcode secrets
```hcl
# ❌
password = "my-secret-password"

# ✅
password = var.db_password  # From tfvars or env
```

## Quick Reference
| Task | Command |
|------|---------|
| Init | `terraform init` |
| Plan | `terraform plan -out=plan.tfplan` |
| Apply | `terraform apply plan.tfplan` |
| Destroy | `terraform destroy` |
| Format | `terraform fmt -recursive` |
| Validate | `terraform validate` |
