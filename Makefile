darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/fabric-gen-darwin-amd64 .

linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/fabric-gen-linux-amd64 .

