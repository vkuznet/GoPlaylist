GOARCHS=amd64 arm64
GOOSES=linux darwin windows

all: build

build:
	go build -o DJserver main.go

build_all:
	@for GOARCH in $(GOARCHS); do \
		for GOOS in $(GOOSES); do \
			echo "Building for $$GOOS/$$GOARCH"; \
			GOOS=$$GOOS GOARCH=$$GOARCH go build -o bin/DJserver-$$GOOS-$$GOARCH main.go; \
		done \
	done

clean:
	rm -rf bin/

