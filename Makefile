all:
	go build -o goext4 bin/*.go


windows:
	GOOS=windows GOARCH=amd64 \
            go build -ldflags="-s -w" \
	    -o goext4.exe ./bin/*.go

generate:
	cd parser/ && binparsegen conversion.spec.yaml > ext4_gen.go

test:
	go test -v ./...
