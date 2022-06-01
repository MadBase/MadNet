SHELL=/bin/bash

BINARY_NAME=madnet
RACE_DETECTOR=madrace

init:
	./scripts/base-scripts/init-githooks.sh

build: init
	go build -o $(BINARY_NAME) ./cmd/main.go;

race:
	go build -o $(RACE_DETECTOR) -race ./cmd/main.go;

generate: generate-bridge generate-go

generate-bridge: init
	find . -iname \*.capnp.go \
	       -o -iname bridge/bindings \
		   -exec rm -rf {} \;
	cd bridge && npm install && npm run build

generate-go: init
	find . -iname \*.pb.go \
	    -o -iname \*.pb.gw.go \
	    -o -iname \*_mngen.go \
		-o -iname \*_mngen_test.go \
		-o -iname \*.swagger.json \
		-o -iname \*.mockgen.go \
		-exec rm -rf {} \;
	go generate ./...

clean:
	go clean
	rm -f $(BINARY_NAME) $(RACE_DETECTOR)