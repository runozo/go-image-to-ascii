all: windows linux darwin

windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -v" -o ./bin/go-image-to-ascii.exe main.go

linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -v" -o ./bin/go-image-to-ascii-linux main.go

# darwin:
# 	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -v" -o ./bin/go-image-to-ascii-darwin main.go
