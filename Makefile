TAG?=0.1.12
NAME:=sas
VERSION:=$(shell grep 'VERSION' pkg/version/version.go | awk '{ print $$4 }' | tr -d '"')

build:
	GIT_COMMIT=$$(git rev-list -1 HEAD) && \
	CGO_ENABLED=0 \
	go build -a -ldflags "-s -w \
	-X github.com/SonicCloudOrg/sonic-android-supply/pkg/version.REVISION=$(GIT_COMMIT)" \
	-o sas \
	./cmd/*

build-linux:
	GIT_COMMIT=$$(git rev-list -1 HEAD) && \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -a -ldflags "-s -w \
	-X github.com/SonicCloudOrg/sonic-android-supply/pkg/version.REVISION=$(GIT_COMMIT)" \
	-o sas \
	./cmd/*

version-set:
	@next="$(TAG)" && \
	current="$(VERSION)" && \
	/usr/bin/sed -i '' "s/$$current/$$next/g" pkg/version/version.go && \
	echo "Version $$next set in code"

release:
	git tag $(VERSION)
	git push origin $(VERSION)