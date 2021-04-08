

.PHONY: all

all:
	CGO_ENABLED=0  GOOS=linux GOARCH=amd64 go build -o dist/ -ldflags "-w -s" -v ./cmd

clean:
	rm -f dist/*