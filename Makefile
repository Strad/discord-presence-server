build: build_win build_linux build_darwin

build_win:
	@echo "Building Windows targets"
	rsrc -manifest app.manifest -ico icon/icon.ico
	GOOS=windows GO111MODULE=on go build -o bin/discord-presence-server-win-debug.exe
	GOOS=windows GO111MODULE=on go build -ldflags "-H=windowsgui" -o bin/discord-presence-server-win.exe

build_linux:
	@echo "Building Linux targets"
	GOOS=linux GO111MODULE=on go build -o bin/discord-presence-server-linux-amd64

build_darwin:
	@echo "Building macOS targets"
	GOOS=darwin GOARCH=arm64 GO111MODULE=on go build -o bin/discord-presence-server-darwin-arm64