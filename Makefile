build-linux:
	CGO_ENABLED=0 GOOS=linux go build -o bin/linux/godown


build-win:
	CGO_ENABLED=0 GOOS=windows go build -o bin/win/godown.exe


build-darwin:
	CGO_ENABLED=0 GOOS=darwin go build -o bin/darwin/godown


build-all: build-win build-darwin build-linux


release-linux:
	tar -c bin/linux -f release/godown_0.1.0_linux.tar


release-darwin:
	tar -c bin/linux -f release/godown_0.1.0_darwin.tar


release-win:
	zip -r release/godown_0.1.0_windows_amd64.zip bin/win


release-all: release-linux release-darwin release-win

test:
	go test ./...