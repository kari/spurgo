GOARCH = amd64 /* FIXME: ARM for Mac
name = spurgo
bins = $(name) $(name)-linux

all: $(bins)

$(name): spurgo.go
	GOOS=darwin GOARCH=arm64 go build -o $@

$(name)-linux: spurgo.go
	GOOS=linux GOARCH=amd64 go build -o $@

clean:
	go clean
	rm -f $(bins)
