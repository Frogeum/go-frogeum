# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gfro android ios gfro-cross evm all test clean
.PHONY: gfro-linux gfro-linux-386 gfro-linux-amd64 gfro-linux-mips64 gfro-linux-mips64le
.PHONY: gfro-linux-arm gfro-linux-arm-5 gfro-linux-arm-6 gfro-linux-arm-7 gfro-linux-arm64
.PHONY: gfro-darwin gfro-darwin-386 gfro-darwin-amd64
.PHONY: gfro-windows gfro-windows-386 gfro-windows-amd64

GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run

gfro:
	$(GORUN) build/ci.go install ./cmd/gfro
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gfro\" to launch gfro."

all:
	$(GORUN) build/ci.go install

android:
	$(GORUN) build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gfro.aar\" to use the library."
	@echo "Import \"$(GOBIN)/gfro-sources.jar\" to add javadocs"
	@echo "For more info see https://stackoverflow.com/questions/20994336/android-studio-how-to-attach-javadoc"

ios:
	$(GORUN) build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Gfro.framework\" to use the library."

test: all
	$(GORUN) build/ci.go test

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean:
	env GO111MODULE=on go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go install golang.org/x/tools/cmd/stringer@latest
	env GOBIN= go install github.com/kevinburke/go-bindata/go-bindata@latest
	env GOBIN= go install github.com/fjl/gencodec@latest
	env GOBIN= go install github.com/golang/protobuf/protoc-gen-go@latest
	env GOBIN= go install ./cmd/abigen
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

gfro-cross: gfro-linux gfro-darwin gfro-windows gfro-android gfro-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/gfro-*

gfro-linux: gfro-linux-386 gfro-linux-amd64 gfro-linux-arm gfro-linux-mips64 gfro-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gfro-linux-*

gfro-linux-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/gfro
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/gfro-linux-* | grep 386

gfro-linux-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/gfro
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gfro-linux-* | grep amd64

gfro-linux-arm: gfro-linux-arm-5 gfro-linux-arm-6 gfro-linux-arm-7 gfro-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/gfro-linux-* | grep arm

gfro-linux-arm-5:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/gfro
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/gfro-linux-* | grep arm-5

gfro-linux-arm-6:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/gfro
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/gfro-linux-* | grep arm-6

gfro-linux-arm-7:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/gfro
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/gfro-linux-* | grep arm-7

gfro-linux-arm64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/gfro
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/gfro-linux-* | grep arm64

gfro-linux-mips:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gfro
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/gfro-linux-* | grep mips

gfro-linux-mipsle:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gfro
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/gfro-linux-* | grep mipsle

gfro-linux-mips64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gfro
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/gfro-linux-* | grep mips64

gfro-linux-mips64le:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gfro
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/gfro-linux-* | grep mips64le

gfro-darwin: gfro-darwin-386 gfro-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gfro-darwin-*

gfro-darwin-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/gfro
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/gfro-darwin-* | grep 386

gfro-darwin-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gfro
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gfro-darwin-* | grep amd64

gfro-windows: gfro-windows-386 gfro-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gfro-windows-*

gfro-windows-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/gfro
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/gfro-windows-* | grep 386

gfro-windows-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/gfro
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gfro-windows-* | grep amd64
