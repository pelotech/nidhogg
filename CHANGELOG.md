# Changelog

## [0.6.6](https://github.com/pelotech/nidhogg/compare/v0.6.5...v0.6.6) (2025-03-18)


### Bug Fixes

* **deps:** update kubernetes packages to v0.32.3 ([#63](https://github.com/pelotech/nidhogg/issues/63)) ([134bbbc](https://github.com/pelotech/nidhogg/commit/134bbbc9acf2dd5d4c7cf021ae4626a3cf916093))
* **deps:** update module github.com/prometheus/client_golang to v1.21.0 ([#58](https://github.com/pelotech/nidhogg/issues/58)) ([51d7a7b](https://github.com/pelotech/nidhogg/commit/51d7a7b01d59db2376177c21d52815b5aa2a8d5b))
* **deps:** update module github.com/prometheus/client_golang to v1.21.1 ([#60](https://github.com/pelotech/nidhogg/issues/60)) ([6c2a3ff](https://github.com/pelotech/nidhogg/commit/6c2a3ffa2cc12fcdbc8b9f7df22c53f0939b17c7))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.20.3 ([#62](https://github.com/pelotech/nidhogg/issues/62)) ([35e6bd6](https://github.com/pelotech/nidhogg/commit/35e6bd62c4df8f0d67a88f487666dba852cdf7c6))


### Reverts

* verbose logging ([#64](https://github.com/pelotech/nidhogg/issues/64)) ([67bd868](https://github.com/pelotech/nidhogg/commit/67bd8688d12cb7223f461d324906fa19781fdfa5))

## [0.6.5](https://github.com/pelotech/nidhogg/compare/v0.6.4...v0.6.5) (2025-02-17)


### Bug Fixes

* add logs to help detect re-added taints ([#56](https://github.com/pelotech/nidhogg/issues/56)) ([dac2d16](https://github.com/pelotech/nidhogg/commit/dac2d16c989630d62d247463d5c43cf5c12115d8))

## [0.6.4](https://github.com/pelotech/nidhogg/compare/v0.6.3...v0.6.4) (2025-02-15)


### Bug Fixes

* **deps:** update kubernetes packages to v0.32.2 ([#7](https://github.com/pelotech/nidhogg/issues/7)) ([b603fb7](https://github.com/pelotech/nidhogg/commit/b603fb73d4bc939e463894b585fa08d3faa6eef6))
* **deps:** update module github.com/prometheus/client_golang to v1.20.5 ([#9](https://github.com/pelotech/nidhogg/issues/9)) ([4e414c0](https://github.com/pelotech/nidhogg/commit/4e414c0ccfed7e67a5127bc4c60c00ab36052ea7))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.20.2 ([#51](https://github.com/pelotech/nidhogg/issues/51)) ([e57c21d](https://github.com/pelotech/nidhogg/commit/e57c21d963816f6baa4b07cc07d955e902258e20))
* upgrade controller runtime to v0.20.1 and conforming to new TypedEventHandler interface ([#46](https://github.com/pelotech/nidhogg/issues/46)) ([9d3f9c1](https://github.com/pelotech/nidhogg/commit/9d3f9c1b787f5f7a11fda6b6b3b64d514d394300))

## [0.6.3](https://github.com/pelotech/nidhogg/compare/v0.6.2...v0.6.3) (2025-02-14)


### Bug Fixes

* update template for release please to also update appVersion in the template ([cc23dbd](https://github.com/pelotech/nidhogg/commit/cc23dbd5ea5ad2f8f38626699c4d36038b420d6b))

## [0.6.2](https://github.com/pelotech/nidhogg/compare/v0.6.1...v0.6.2) (2025-02-14)


### Bug Fixes

* update chart values to update app and chart version in sync ([2204fda](https://github.com/pelotech/nidhogg/commit/2204fda99ea41abb3ff5748c9e5d774738dae7a4))

## [0.6.1](https://github.com/pelotech/nidhogg/compare/v0.6.0...v0.6.1) (2025-02-14)


### Bug Fixes

* update to remove fromJson from job outputs ([b262da8](https://github.com/pelotech/nidhogg/commit/b262da8b4a061ade45cc476771ee7c074cee6a91))

## [0.6.0](https://github.com/pelotech/nidhogg/compare/v0.5.3...v0.6.0) (2025-02-13)


### Features

* add pr-title check and pre-commit action ([#29](https://github.com/pelotech/nidhogg/issues/29)) ([e5775b6](https://github.com/pelotech/nidhogg/commit/e5775b6639c8866cb946d159926d9530ba08ee0a))
* **selectors:** Extract selectors from daemonsets if not provided through config ([#33](https://github.com/pelotech/nidhogg/issues/33)) ([dbdb572](https://github.com/pelotech/nidhogg/commit/dbdb5727ff2e986c73ce7fae492dc6ba9f662d3e))


### Bug Fixes

* **badges:** Fix GH workflow badges in docs/README.md ([3850e21](https://github.com/pelotech/nidhogg/commit/3850e2119e8559b7d621ed54e87c70326f40c904))
* **deps:** update module github.com/onsi/gomega to v1.36.2 ([#8](https://github.com/pelotech/nidhogg/issues/8)) ([fc4103c](https://github.com/pelotech/nidhogg/commit/fc4103c9514175cbcd555e7b3c283a4f05f0500d))
* **deps:** update module github.com/stretchr/testify to v1.10.0 ([#10](https://github.com/pelotech/nidhogg/issues/10)) ([eabd5fa](https://github.com/pelotech/nidhogg/commit/eabd5faed9c6d855250ca7e7f1ff52eda1a789c2))
* removed duplicated config for pre-commit ([#30](https://github.com/pelotech/nidhogg/issues/30)) ([e9d08b8](https://github.com/pelotech/nidhogg/commit/e9d08b8dad1ce7ce9b204cb3f79c67b20b748008))
* using container image directly in resources.yaml ([#37](https://github.com/pelotech/nidhogg/issues/37)) ([1f05333](https://github.com/pelotech/nidhogg/commit/1f053339642edb083decdf03e92709de433eec06))
