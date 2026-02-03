PLUGIN_NAME := release
BINARY := plugin-$(PLUGIN_NAME)
INSTALL_DIR := $(HOME)/.neko/plugins/$(PLUGIN_NAME)

.PHONY: build install clean

build:
	go build -o $(BINARY) main.go

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY) $(INSTALL_DIR)/
	cp manifest.json $(INSTALL_DIR)/

clean:
	rm -f $(BINARY)
	rm -rf $(INSTALL_DIR)