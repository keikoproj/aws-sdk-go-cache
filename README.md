# aws-sdk-go-cache

> **DEPRECATED — This library is no longer maintained and will not receive any further updates.** There will be no v0.2.0 release. Users should migrate to [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2) and rely on its built-in retry and rate-limiting capabilities instead (see [Migration Guide](#migration-guide) below).

[![Go Report Card][GoReportImg]][GoReportUrl]
[![Build Status][BuildStatusImg]][BuildMasterUrl]
[![codecov][CodecovImg]][CodecovUrl]

This package provided a response-caching layer for [AWS SDK for Go v1](https://github.com/aws/aws-sdk-go), which reached end-of-support on July 31, 2025. It was forked from [ticketmaster/aws-sdk-go-cache](https://github.com/ticketmaster/aws-sdk-go-cache).

## Migration Guide

This library was originally built to mitigate AWS API rate limiting (throttling) by caching read-only API responses. AWS SDK for Go v2 significantly reduces the need for client-side response caching:

### AWS SDK v2 Built-in Rate Limiting

**Standard Retry Mode** (enabled by default):
- Token-bucket rate limiter (500 tokens) that automatically throttles requests when errors are detected
- Exponential backoff with jitter between retry attempts
- 3 retry attempts by default, configurable up to any number
- Safe for multi-tenant applications

**Adaptive Retry Mode** (opt-in):
- Dynamically adjusts client-side request rate based on throttle responses from AWS
- Automatically restricts attempt rate when throttling is detected
- Expands on standard mode with a second token bucket for request rate limiting

### Recommended Approach

1. **Remove this library** and migrate to AWS SDK for Go v2
2. **Use standard retry mode** (the default) — no extra configuration needed
3. **For truly stable data** (e.g., `DescribeInstanceTypes` which rarely changes), use a simple `sync.Once` or package-level cache variable rather than a full caching library
4. **If you need custom response caching**, AWS SDK v2's [middleware framework](https://aws.github.io/aws-sdk-go-v2/docs/middleware/) makes it straightforward to add a caching layer without a dedicated library

### References

- [AWS SDK for Go v2 — Retries and Timeouts](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/retries-timeouts/)
- [AWS SDK for Go v2 — Migration Guide](https://aws.github.io/aws-sdk-go-v2/docs/migrating/)
- [AWS SDKs Retry Behavior Reference](https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html)


<!-- Markdown link -->
[GoReportImg]: https://goreportcard.com/badge/github.com/keikoproj/aws-sdk-go-cache
[GoReportUrl]: https://goreportcard.com/report/github.com/keikoproj/aws-sdk-go-cache

[BuildStatusImg]: https://github.com/keikoproj/aws-sdk-go-cache/actions/workflows/unit-test.yaml/badge.svg?branch=master
[BuildMasterUrl]: https://github.com/keikoproj/aws-sdk-go-cache/actions/workflows/unit-test.yaml

[CodecovImg]: https://codecov.io/gh/keikoproj/aws-sdk-go-cache/branch/master/graph/badge.svg
[CodecovUrl]: https://codecov.io/gh/keikoproj/aws-sdk-go-cache
