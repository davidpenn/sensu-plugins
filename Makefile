GOOS = $(shell go env | grep GOOS= | cut -d \" -f2)
GOARCH = $(shell go env | grep GOARCH= | cut -d \" -f2)
SOURCES = $(patsubst ./cmd/%/main.go,%,$(shell find ./cmd -maxdepth 2 -name 'main.go'))
TARGETS = $(patsubst %,bin/$(GOOS)/$(GOARCH)/%,$(SOURCES))
GO_FLAGS = -a -tags netgo -ldflags "-s -w"

get_os = $(shell dirname $(patsubst bin/%/,%,$(dir $@)))
get_arch = $(shell basename $(patsubst bin/%/,%,$(dir $@)))

.PHONY: all
all: build

.PHONY: build
build: $(TARGETS)

bin/%:
	GOOS=$(call get_os,$@) GOARCH=$(call get_arch,$@) CGO_ENABLED=0 go build $(GO_FLAGS) -o $@ cmd/$(notdir $*)/main.go

.PHONY: clean
clean:
	rm -rf bin/ dist/
