.PHONY: build clean tool lint help

run:
	go run .
prod:
	GOOS=linux go build -o bin/admin .
update:
	hz update -I idl --idl=idl/$(p).proto --customize_package=template/package.yaml: --snake_tag --no_recurse
	protoc-go-inject-tag -input=biz/model/$(p)/*.pb.go
