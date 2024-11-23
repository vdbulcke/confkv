
gen-proto: 
	protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative ./src/pb/service.proto

scan: 
	trivy fs .

build: 
	goreleaser build --clean

build-snapshot: 
	goreleaser build --clean --snapshot --single-target



release-snapshot: 
	goreleaser release --clean  --snapshot  --skip=sign


lint: 
	golangci-lint run ./...

changelog:
	# git-chglog -o CHANGELOG.md
	git cliff -o CHANGELOG.md

test:
	go test -v  ./...

