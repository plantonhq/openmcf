.PHONY: build
build: lint fmt
	rm -rf generated/stubs ../app/frontend/src/gen
	mkdir -p generated/stubs ../app/frontend/src/gen
	buf generate
	cp -R generated/stubs/go/github.com/plantonhq/openmcf/apis/. .
	cp -R generated/stubs/ts/. ../app/frontend/src/gen/
	# Keep Java stubs in generated/ for early detection of reserved keyword issues
	# Go and TS stubs are copied to their destinations, so we only clean those
	rm -rf generated/stubs/go generated/stubs/ts

.PHONY: protos
protos: build

.PHONY: fmt
fmt:
	buf format -w

.PHONY: lint
lint: buf-lint

.PHONY: build-optional-linter-plugin
build-optional-linter-plugin:
	@echo "Building buf lint plugin..."
	@cd ../buf/lint/optional-linter && $(MAKE) build

.PHONY: buf-lint
buf-lint: build-optional-linter-plugin
	buf lint
