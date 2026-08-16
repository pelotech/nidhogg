# Changelog

## [0.8.2](https://github.com/pelotech/nidhogg/compare/v0.8.1...v0.8.2) (2026-08-16)


### Chores

* **deps:** pin dependencies ([#114](https://github.com/pelotech/nidhogg/issues/114)) ([ebd48cf](https://github.com/pelotech/nidhogg/commit/ebd48cf61690e8a57a17d9371589530bc51f86a4))
* **deps:** update actions/checkout action to v7 ([#110](https://github.com/pelotech/nidhogg/issues/110)) ([b69392d](https://github.com/pelotech/nidhogg/commit/b69392df2a1047fa7ab2823277884e0646a8dd72))
* **deps:** update actions/setup-go action to v7 ([#111](https://github.com/pelotech/nidhogg/issues/111)) ([f5b5228](https://github.com/pelotech/nidhogg/commit/f5b5228def1f098d50091a10922700deb575f827))
* **deps:** update googleapis/release-please-action action to v5 ([#108](https://github.com/pelotech/nidhogg/issues/108)) ([8266ffc](https://github.com/pelotech/nidhogg/commit/8266ffc1ebabd9ccf8707805f8f7d506b5798ea1))
* **deps:** update pre-commit hook adrienverge/yamllint to v1.38.0 ([#116](https://github.com/pelotech/nidhogg/issues/116)) ([a877a4d](https://github.com/pelotech/nidhogg/commit/a877a4d144802f011708b84c2b44af80b37caf1c))
* **deps:** update pre-commit hook alessandrojcm/commitlint-pre-commit-hook to v9.26.0 ([#117](https://github.com/pelotech/nidhogg/issues/117)) ([b96706b](https://github.com/pelotech/nidhogg/commit/b96706bfc5b712c07c855e01e779767f0971f447))
* **deps:** update pre-commit hook gruntwork-io/pre-commit to v0.1.30 ([#115](https://github.com/pelotech/nidhogg/issues/115)) ([71fb867](https://github.com/pelotech/nidhogg/commit/71fb86778c2a079cb02c57831d9c3a0d1080ee99))
* update renovate config ([1c996c9](https://github.com/pelotech/nidhogg/commit/1c996c902b769901e641a586ff46d7f35bf5e724))

## [0.8.1](https://github.com/pelotech/nidhogg/compare/v0.8.0...v0.8.1) (2026-03-27)


### Bug Fixes

* remove log message about missing/unknown taint effect ([#98](https://github.com/pelotech/nidhogg/issues/98)) ([#99](https://github.com/pelotech/nidhogg/issues/99)) ([af89b30](https://github.com/pelotech/nidhogg/commit/af89b3053150d67270b1a6c2183bb4fec55718c9))


### Chores

* **deps:** update docker/build-push-action action to v7 ([#105](https://github.com/pelotech/nidhogg/issues/105)) ([d5572ea](https://github.com/pelotech/nidhogg/commit/d5572ea5dc092c888a3b13c24617efaecc8f2e73))
* **deps:** update docker/login-action action to v4 ([#102](https://github.com/pelotech/nidhogg/issues/102)) ([66d5e63](https://github.com/pelotech/nidhogg/commit/66d5e63c0f6099b907853fb67e0d186c91a11424))
* **deps:** update docker/metadata-action action to v6 ([#104](https://github.com/pelotech/nidhogg/issues/104)) ([1ce4b71](https://github.com/pelotech/nidhogg/commit/1ce4b71029b2e11f47a7fc0bd499f8d353993c46))
* **deps:** update docker/setup-buildx-action action to v4 ([#103](https://github.com/pelotech/nidhogg/issues/103)) ([611040a](https://github.com/pelotech/nidhogg/commit/611040a7b8373a2e8821c701149ce6320f1af1cb))
* **deps:** update docker/setup-qemu-action action to v4 ([#101](https://github.com/pelotech/nidhogg/issues/101)) ([e3716a0](https://github.com/pelotech/nidhogg/commit/e3716a094ed4e1c964daf4d13510b172f18004b8))

## [0.8.0](https://github.com/pelotech/nidhogg/compare/v0.7.1...v0.8.0) (2025-12-22)


### Features

* add option for configuring taint effects ([7b20a26](https://github.com/pelotech/nidhogg/commit/7b20a26f09ad4b4c22ff019f2a4b4d654c98b408))

## [0.7.1](https://github.com/pelotech/nidhogg/compare/v0.7.0...v0.7.1) (2025-12-22)


### Bug Fixes

* **deps:** update kubernetes packages to v0.35.0 ([#80](https://github.com/pelotech/nidhogg/issues/80)) ([f4b9fe0](https://github.com/pelotech/nidhogg/commit/f4b9fe08b772708c3fbacc37fbdf0826ba592de6))
* **deps:** update module github.com/onsi/gomega to v1.38.3 ([#82](https://github.com/pelotech/nidhogg/issues/82)) ([b34cf62](https://github.com/pelotech/nidhogg/commit/b34cf629ff98e604f308e1daf609834e712f0cfb))
* **deps:** update module github.com/prometheus/client_golang to v1.23.2 ([#83](https://github.com/pelotech/nidhogg/issues/83)) ([f0d3d09](https://github.com/pelotech/nidhogg/commit/f0d3d09cf9cde2631ad4ea2a7c52e8537f77bbc0))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.22.4 ([#86](https://github.com/pelotech/nidhogg/issues/86)) ([8767981](https://github.com/pelotech/nidhogg/commit/87679818b198724de360921794b960834067493f))


### Chores

* **deps:** update actions/checkout action to v6 ([#89](https://github.com/pelotech/nidhogg/issues/89)) ([68f9dc6](https://github.com/pelotech/nidhogg/commit/68f9dc6c82600f14233a44959746d587f702b450))
* **deps:** update actions/setup-go action to v6 ([#88](https://github.com/pelotech/nidhogg/issues/88)) ([af5cc36](https://github.com/pelotech/nidhogg/commit/af5cc36c8a7de252537c831bc87efe66b1d92745))
* **deps:** update appany/helm-oci-chart-releaser action to v0.5.0 ([#81](https://github.com/pelotech/nidhogg/issues/81)) ([e675a37](https://github.com/pelotech/nidhogg/commit/e675a37173dd68f4d20fad1f13872f17053a1a78))
* **deps:** update dependency go to 1.25.x ([#84](https://github.com/pelotech/nidhogg/issues/84)) ([dea9479](https://github.com/pelotech/nidhogg/commit/dea94793cef9072d5d36c29a2356c6ba69524969))
* **deps:** update golang docker tag to v1.25.5 ([#79](https://github.com/pelotech/nidhogg/issues/79)) ([ebbc71c](https://github.com/pelotech/nidhogg/commit/ebbc71c7a079b14ce1fe7c0249a60cdb4367ce45))
* upgrade to go 1.25 ([4696d43](https://github.com/pelotech/nidhogg/commit/4696d434cf59b3912ffba9e609029c6120f91e1e))
* use go-version-file instead of hard coded value in action ([900b5b4](https://github.com/pelotech/nidhogg/commit/900b5b405d0f44fff1c9fb0159d0b83d167710a8))


### Docs

* update release please config to include more types ([8a8546a](https://github.com/pelotech/nidhogg/commit/8a8546a25c53289979d8b91d87e7a54b4030aac6))

## [0.7.0](https://github.com/pelotech/nidhogg/compare/v0.6.6...v0.7.0) (2025-05-28)


### Features

* adding ratelimiting parameters ([#74](https://github.com/pelotech/nidhogg/issues/74)) ([c7373e7](https://github.com/pelotech/nidhogg/commit/c7373e723f37531057f45a67d521d28871f2435a))


### Bug Fixes

* **deps:** update kubernetes packages to v0.33.0 ([#72](https://github.com/pelotech/nidhogg/issues/72)) ([10ac8e3](https://github.com/pelotech/nidhogg/commit/10ac8e38a272ef74c181217a14780650cd1e73fc))
* **deps:** update kubernetes packages to v0.33.1 ([#76](https://github.com/pelotech/nidhogg/issues/76)) ([3ffaf23](https://github.com/pelotech/nidhogg/commit/3ffaf23801e96d6c94c910cd562a19997f2f1082))
* **deps:** update module github.com/onsi/gomega to v1.36.3 ([#65](https://github.com/pelotech/nidhogg/issues/65)) ([4d593fd](https://github.com/pelotech/nidhogg/commit/4d593fd4d0c4e2a0e529fa56092abd6847d9d251))
* **deps:** update module github.com/onsi/gomega to v1.37.0 ([#69](https://github.com/pelotech/nidhogg/issues/69)) ([570f695](https://github.com/pelotech/nidhogg/commit/570f6952e4dc1e7dc82e35dac094de75d4f331d8))
* **deps:** update module github.com/prometheus/client_golang to v1.22.0 ([#70](https://github.com/pelotech/nidhogg/issues/70)) ([b661379](https://github.com/pelotech/nidhogg/commit/b661379e6177a1f6d855ea891907a46aa39b2101))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.20.4 ([#67](https://github.com/pelotech/nidhogg/issues/67)) ([8eda4fe](https://github.com/pelotech/nidhogg/commit/8eda4febe58e6858b9cfadca7d3d76f795aff4ab))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.21.0 ([#77](https://github.com/pelotech/nidhogg/issues/77)) ([0e5c680](https://github.com/pelotech/nidhogg/commit/0e5c680e1d093df6251be7c7752c15f0f5f1e9eb))

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
