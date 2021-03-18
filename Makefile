.PHONY: build

VETPACKAGES=`go list ./... | grep -v /vendor/ | grep -v /examples/`
GOFILES=`find . -name "*.go" -type f -not -path "./vendor/*"`

gofmt:
		echo "正在使用gofmt格式化文件..."
		gofmt -s -w ${GOFILES}
		echo "格式化完成"
govet:
		echo "正在进行静态检测..."
		go vet $(VETPACKAGES)