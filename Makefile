
.PHONY: build

build:
	go get github.com/yuce/picon/cmd/picon
	go build github.com/yuce/picon/cmd/picon
