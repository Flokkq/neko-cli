VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -X github.com/nekoman-hq/neko-cli/internal/version.Version=$(VERSION) \
           -X github.com/nekoman-hq/neko-cli/internal/version.Commit=$(COMMIT) \
           -X github.com/nekoman-hq/neko-cli/internal/version.Date=$(DATE) \
           -X github.com/nekoman-hq/neko-cli/internal/version.BuiltBy=make

.PHONY: build install clean install-plugins test

build:
	go build -ldflags "$(LDFLAGS)" -o neko

install:
	go install -ldflags "$(LDFLAGS)"

# Plugins bauen und installieren
install-plugins:
	cd plugin/release && $(MAKE) install


# Alles bauen und installieren
all: build install-plugins

clean:
	rm -f neko
	cd plugin/release && $(MAKE) clean || true
	cd plugin/core && $(MAKE) clean || true

test:
	go test ./...