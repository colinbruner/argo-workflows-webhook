# Argo Workflow Webhooks

Mutating webhooks for manipulating CRDs created by Argo Workflows.

## Goal

I wanted to get more acquainted with Kubernetes [Dynamic Admission Control][dac] through mutating and validating webhooks. I chose Argo Workflow webhooks as I use this within my home lab, which makes testing and experimentation very easy.

### Design Principles

- Simple Framework
- Utilzing Golang STDLIB as much as possible (http, logging, etc)
- Attempting to use Golang idiomatic practices
- Testing!

[dac]: https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/
