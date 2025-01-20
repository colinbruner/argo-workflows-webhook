# TODO

- Enable Prometheus Metrics
- Enable Tracing
- Implement Mutate/Validate for Workflows, and/or WorkflowTemplates
- Refactor Router package for better testing, mTLS, [authn/authz with Kubelet][auth]

## Maybe TODO

- Utilize [Zap][logging] for structured Logging?
- Use an HTTP framework, use as Chi or Gin?

[auth]: https://kubernetes.io/docs/reference/access-authn-authz/kubelet-authn-authz/
[logging]: https://betterstack.com/community/guides/logging/go/zap/
