SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o windows_inotify_consume.exe
7z a bin/inotify_consume.zip ./

