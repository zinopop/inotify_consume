SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o build/darwin_consume

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o build/windows_consume.exe

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o build/linux_consume

copy config.toml build

7z a bin/inotify_consume.zip build

