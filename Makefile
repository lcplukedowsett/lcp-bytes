TEST?=$(shell go list ./... | findstr /V vendor)
HOSTNAME=terraform.test.com
NAMESPACE=alpha
NAME=bytes
BINARY=terraform-provider-${NAME}
VERSION=0.2.0
OS_ARCH=windows_amd64

default: install

build:
	go build -o ${BINARY}.exe

release:
	goreleaser release --rm-dist --snapshot --skip-publish --skip-sign

install: build
	if not exist "%APPDATA%\terraform.d\plugins\${HOSTNAME}\${NAMESPACE}\${NAME}\${VERSION}\${OS_ARCH}" mkdir "%APPDATA%\terraform.d\plugins\${HOSTNAME}\${NAMESPACE}\${NAME}\${VERSION}\${OS_ARCH}"
	move ${BINARY}.exe "%APPDATA%\terraform.d\plugins\${HOSTNAME}\${NAMESPACE}\${NAME}\${VERSION}\${OS_ARCH}"

test:
	go test -i $(TEST) || exit 1
	for %%G in ($(TEST)) do go test -timeout=30s -parallel=4 %%G $(TESTARGS)

testacc:
	set TF_ACC=1
	go test $(TEST) -v $(TESTARGS) -timeout 120m