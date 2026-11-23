# Changelog

## [1.38.1](https://github.com/env0/terraform-provider-env0/compare/v1.38.0...v1.38.1) (2026-11-23)


### Bug Fixes

* Name the field in pattern validation errors (APO-1641) ([#1184](https://github.com/env0/terraform-provider-env0/issues/1184)) ([b76ea96](https://github.com/env0/terraform-provider-env0/commit/b76ea96af2a88a5507cb6ea3083ea93c26b03827))

## [1.38.0](https://github.com/env0/terraform-provider-env0/compare/v1.37.0...v1.38.0) (2026-11-19)


### Features

* Add env0_approval_policy data source (APO-1657) ([#1181](https://github.com/env0/terraform-provider-env0/issues/1181)) ([3c39902](https://github.com/env0/terraform-provider-env0/commit/3c399024c01b2a082d5e0008209a5a2311e7b109))


### Bug Fixes

* Mark data source id and lookup attributes as Computed (APO-1657) ([#1182](https://github.com/env0/terraform-provider-env0/issues/1182)) ([059a7fb](https://github.com/env0/terraform-provider-env0/commit/059a7fb57d9a7b392f2a39fe1d626d66e8495692))

## [1.37.0](https://github.com/env0/terraform-provider-env0/compare/v1.36.0...v1.37.0) (2026-11-19)


### Features

* Accept LATEST_ALLOWED_FROM_CODE for Terraform and OpenTofu versions (APO-1641) ([#1180](https://github.com/env0/terraform-provider-env0/issues/1180)) ([f952276](https://github.com/env0/terraform-provider-env0/commit/f9522760b8ab1b71b401f027439747c56759ccd2))


### Bug Fixes

* Reset policy booleans to their defaults when deleting env0_organization_policy (ATL-702) ([#1172](https://github.com/env0/terraform-provider-env0/issues/1172)) ([e6cc195](https://github.com/env0/terraform-provider-env0/commit/e6cc195ea5a1d67e75161a4f7d50431ccb3f52bd))
* Stop the provider crashing on parallel env0_team lookups (APO-1537) ([#1178](https://github.com/env0/terraform-provider-env0/issues/1178)) ([3e01528](https://github.com/env0/terraform-provider-env0/commit/3e015281a4f1e21e0d040406740a20d96f64e893))

## [1.36.0](https://github.com/env0/terraform-provider-env0/compare/v1.35.1...v1.36.0) (2026-11-08)


### Features

* Set per-agent IP allowlists on env0_agent_pool (APO-1371) ([#1166](https://github.com/env0/terraform-provider-env0/issues/1166)) ([a0b249e](https://github.com/env0/terraform-provider-env0/commit/a0b249ed326fe3fc20450ae4be1dee416c641a1f))


### Bug Fixes

* **deps:** Update module github.com/cenkalti/backoff/v4 to v7 ([#1162](https://github.com/env0/terraform-provider-env0/issues/1162)) ([2ca7949](https://github.com/env0/terraform-provider-env0/commit/2ca79495fc028d928925e8e94b9e8473bcc16d13))
* **deps:** Update module github.com/onsi/ginkgo to v2 ([#1163](https://github.com/env0/terraform-provider-env0/issues/1163)) ([b42a171](https://github.com/env0/terraform-provider-env0/commit/b42a1714bcb04745dccc1ed66640f771cfd463b4))
* **deps:** Update module google.golang.org/grpc to v1.83.2 [SECURITY] ([#1169](https://github.com/env0/terraform-provider-env0/issues/1169)) ([c2b5ce8](https://github.com/env0/terraform-provider-env0/commit/c2b5ce8a68f882bc79fb0006e6dd979d722305b4))
* **deps:** Update non-major ([#1160](https://github.com/env0/terraform-provider-env0/issues/1160)) ([4ff8f83](https://github.com/env0/terraform-provider-env0/commit/4ff8f83a061864125a31bf41ec7a1fd0a65def0c))
* Send organizationId when reading variables from a repository ([#1173](https://github.com/env0/terraform-provider-env0/issues/1173)) ([5133b67](https://github.com/env0/terraform-provider-env0/commit/5133b67784eb292669079dc362a2296eeefbbb7e))

## [1.35.1](https://github.com/env0/terraform-provider-env0/compare/v1.35.0...v1.35.1) (2026-11-02)


### Bug Fixes

* **deps:** Update module google.golang.org/grpc to v1.83.1 [SECURITY] ([#1157](https://github.com/env0/terraform-provider-env0/issues/1157)) ([b6b62a2](https://github.com/env0/terraform-provider-env0/commit/b6b62a2180919212b32161ac2d73234ad38a0927))
