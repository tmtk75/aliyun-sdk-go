aly: clean
	go build -o aly main.go

generate: gen/gen ./node_modules/aliyun-sdk
	go generate

gen/gen: gen/*.go
	cd gen; go build

./node_modules/aliyun-sdk:
	npm install aliyun-sdk

version=`./aly -v | sed 's/aly version //g'`

release: ./aly
	rm -f pkg/*_amd64 pkg/*.exe
	ghr -u tmtk75 v$(version) pkg

compress:
	gzip -fk pkg/*_amd64
	zip pkg/aly_windows_amd64.zip pkg/aly_windows_amd64.exe

build:
	gox \
	  -os "linux darwin windows" \
	  -arch "amd64" \
	  -output "pkg/aly_{{.OS}}_{{.Arch}}" \
	  .

clean:
	rm -f aly

distclean: clean
	rm -rf pkg gen/gen
