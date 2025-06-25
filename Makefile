BINARY=zb
CMD_DIR=cmd/zb

.PHONY: build clean

build:
	go build -o $(BINARY) ./$(CMD_DIR)

clean:
	rm -f $(BINARY)