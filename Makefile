BINARY_NAME := atf
INSTALL_DIR := /usr/local/bin

.PHONY: build install uninstall clean

build:
	go build -o $(BINARY_NAME) .

install: build
	install -Dm755 $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
