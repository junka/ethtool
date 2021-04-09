

.PHONY: all

all:
	CGO_ENABLED=0  GOOS=linux GOARCH=amd64 go build -o dist/ -ldflags "-w -s" -v ./cmd/ethtool

clean:
	rm -f dist/*