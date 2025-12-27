build_x86:
	export GOOS=linux  GOARCH=amd64;\
	go build -o ./out/cronweb-x86_64

build_arm64:
	export GOOS=linux  GOARCH=arm64;\
	go build -o ./out/cronweb-aarch64

build_webapp:
	bun run build
build: build_webapp build_x86 build_arm64
