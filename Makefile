.PHONY: all
all: build generate test

.PHONY: build
build:
	go build -tags=netgo -ldflags="-s -w" -trimpath -o makelens ./cmd/makelens
	go build -tags=netgo -ldflags="-s -w" -trimpath -o exp/makedblens ./exp/cmd/makedblens
	go build -tags=netgo -ldflags="-s -w" -trimpath -o makecolops ./cmd/makecolops
	go build -tags=netgo -ldflags="-s -w" -trimpath -o vet ./cmd/vet	
	GOOS=js GOARCH=wasm go build -ldflags "-w -s" -trimpath  -o docs/static/playground.wasm ./internal/playground
	gzip -v -f -n --best -k docs/static/playground.wasm
.PHONY: build-generate
build-generate:
	go build -tags=netgo,makecolops -ldflags="-s -w" -trimpath -o makelens ./cmd/makelens
	go build -tags=netgo,makecolops -ldflags="-s -w" -trimpath -o exp/makedblens ./exp/cmd/makedblens
	go build -tags=netgo,makecolops -ldflags="-s -w" -trimpath -o makecolops ./cmd/makecolops


.PHONY: generate
generate: build-generate
	go generate ./...

.PHONY: revert-generate
revert-generate: 
	git restore string_ops.go
	git restore map_ops.go
	git restore slice_ops.go

.PHONY: makelens
makelens: build
	go generate ./cmd/makelens

.PHONY: godoc
godoc:
	go run golang.org/x/tools/cmd/godoc@latest -http=:6060


.PHONY: gocleancache
gocleancache:
	go clean -cache

.PHONY: hugodoc
hugodoc:
	cd docs;hugo


.PHONY: hugodoc-server
hugodoc-server:
	cd docs;hugo server --buildDrafts --disableFastRender

.PHONY: test
test:
	go test -p 4 -timeout 30s ./...

.PHONY: test-cover
test-cover:
	go test -p 4 -timeout 30s -coverpkg=./... -coverprofile=cover/coverage.out ./...