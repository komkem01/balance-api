# go-temp
💻 Project structure golang

## 🚀 Starter Project

```bash
# Init Project
go mod vendor

# HTTP/1.1 Server
go run . http
```

## Version Tag Workflow

Use semantic version tags like `v1.2.3`.

```bash
# 1) run tests first
make test

# 2) create + push release tag
make tag-release VERSION=v1.2.3

# or use script (validates format + clean git state)
./scripts/release-tag.sh v1.2.3

# 3) build binary with injected version
make build-release VERSION=v1.2.3
```

Notes:

- App version is injected via ldflags into `balance/internal/config.version`.
- `GET /system/manifest` reads this value and returns it as `version`.
- Runtime version resolution order:
	1) `APP_VERSION` env
	2) ldflags-injected build version
	3) latest git tag (`git describe --tags --abbrev=0`)
	4) fallback `0.0.0`
