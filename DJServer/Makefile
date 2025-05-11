GOARCHS=amd64 arm64 ppc64le
GOOSES=linux darwin windows freebsd

build:
	@for GOARCH in $(GOARCHS); do \
		for GOOS in $(GOOSES); do \
			echo "Building for $$GOOS/$$GOARCH"; \
			GOOS=$$GOOS GOARCH=$$GOARCH go build -o bin/music-server-$$GOOS-$$GOARCH main.go; \
		done \
	done

clean:
	rm -rf bin/

