GOARCH = amd64
name = spurgo
bins = $(name) $(name)-linux

all: $(bins)

$(name): spurgo.go
	GOOS=darwin go build -o $@

$(name)-linux: spurgo.go
	GOOS=linux go build -o $@

clean:
	go clean
	rm -f $(bins)
