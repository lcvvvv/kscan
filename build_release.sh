export LDFLAGS='-s -w '

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -trimpath -o kscan_linux_amd64 kscan.go
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="$LDFLAGS" -trimpath -o kscan_windows_386.exe  kscan.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -trimpath -o kscan_windows_amd64.exe  kscan.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="$LDFLAGS" -trimpath -o kscan_windows_arm64.exe  kscan.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$LDFLAGS" -trimpath -o kscan_darwin_amd64 kscan.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$LDFLAGS" -trimpath -o kscan_darwin_arm64 kscan.go

upx -9 kscan_linux_amd64
upx -9 kscan_windows_386.exe
upx -9 kscan_windows_amd64.exe
upx -9 kscan_windows_arm64.exe
upx -9 kscan_darwin_amd64
upx -9 kscan_darwin_arm64

zip kscan_linux_amd64.zip kscan_linux_amd64
zip kscan_windows_386.zip kscan_windows_386.exe
zip kscan_windows_amd64.zip kscan_windows_amd64.exe
zip kscan_windows_arm64.zip kscan_windows_arm64.exe
zip kscan_darwin_amd64.zip kscan_darwin_amd64
zip kscan_darwin_arm64.zip kscan_darwin_arm64

rm -f kscan_linux_amd64
rm -f kscan_windows_386.exe
rm -f kscan_windows_amd64.exe
rm -f kscan_windows_arm64.exe
rm -f kscan_darwin_amd64
rm -f kscan_darwin_arm64