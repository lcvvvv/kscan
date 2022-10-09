export LDFLAGS='-s -w '

go build -ldflags="$LDFLAGS" -trimpath kscan.go

upx -9 kscan