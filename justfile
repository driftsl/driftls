entry := "cmd/driftls/main.go"

run *ARGS:
    @go run {{entry}} {{ARGS}}

build:
    @go build -o dist/driftls {{entry}}

update-driftc:
    @go get github.com/driftsl/driftc
