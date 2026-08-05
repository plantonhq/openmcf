# --------------------------------------------------------------------------- #
#  Pinned proto toolchain (.tools/) -- local plugin execution
#
#  buf orchestrates codegen (buf.gen.yaml, managed mode), but every generator
#  binary runs locally from the git-ignored .tools/ directory, provisioned on
#  demand by the stamp targets below. Remote BSR plugins are deliberately not
#  used: buf.build's code-generation service enforces a server-side ~120s
#  execution limit that this proto surface exceeds (the hosted Java generator
#  fails with deadline_exceeded), and remote execution ties stub generation to
#  a third-party service being reachable.
#
#  Every pin is coupled to a runtime dependency -- bump both together:
#    PROTOC_VERSION                <-> com.google.protobuf:protobuf-java (MODULE.bazel)
#    PROTOC_GEN_GRPC_JAVA_VERSION  <-> io.grpc:*                         (MODULE.bazel)
#    PROTOC_GEN_GO_VERSION         <-> google.golang.org/protobuf        (go.mod)
#    PROTOC_GEN_GO_GRPC_VERSION    <-> google.golang.org/grpc            (go.mod)
# --------------------------------------------------------------------------- #
PROTO_TOOLS_DIR := .tools

PROTOC_VERSION               := 30.1
PROTOC_GEN_GRPC_JAVA_VERSION := 1.65.0
PROTOC_GEN_GO_VERSION        := v1.36.6
PROTOC_GEN_GO_GRPC_VERSION   := v1.5.1

# protoc and grpc-java release artifacts share this platform naming scheme:
# osx-aarch_64 / osx-x86_64 / linux-aarch_64 / linux-x86_64
PROTO_TOOLS_OS   := $(if $(filter Darwin,$(shell uname -s)),osx,linux)
PROTO_TOOLS_ARCH := $(if $(filter arm64 aarch64,$(shell uname -m)),aarch_64,x86_64)

# Each tool is guarded by a version-stamped marker file: bumping a pin above
# reinstalls that tool on the next build; up-to-date tools cost nothing.
PROTOC_STAMP             := $(PROTO_TOOLS_DIR)/.stamp.protoc.$(PROTOC_VERSION)
GRPC_JAVA_STAMP          := $(PROTO_TOOLS_DIR)/.stamp.protoc-gen-grpc-java.$(PROTOC_GEN_GRPC_JAVA_VERSION)
PROTOC_GEN_GO_STAMP      := $(PROTO_TOOLS_DIR)/.stamp.protoc-gen-go.$(PROTOC_GEN_GO_VERSION)
PROTOC_GEN_GO_GRPC_STAMP := $(PROTO_TOOLS_DIR)/.stamp.protoc-gen-go-grpc.$(PROTOC_GEN_GO_GRPC_VERSION)

$(PROTOC_STAMP):
	@echo "installing protoc $(PROTOC_VERSION) into $(PROTO_TOOLS_DIR)/protoc/"
	@rm -rf $(PROTO_TOOLS_DIR)/protoc $(PROTO_TOOLS_DIR)/.stamp.protoc.*
	@mkdir -p $(PROTO_TOOLS_DIR)/protoc
	@curl -fsSL -o $(PROTO_TOOLS_DIR)/protoc.zip \
		https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(PROTO_TOOLS_OS)-$(PROTO_TOOLS_ARCH).zip
	@unzip -q $(PROTO_TOOLS_DIR)/protoc.zip -d $(PROTO_TOOLS_DIR)/protoc
	@rm -f $(PROTO_TOOLS_DIR)/protoc.zip
	@touch $@

$(GRPC_JAVA_STAMP):
	@echo "installing protoc-gen-grpc-java $(PROTOC_GEN_GRPC_JAVA_VERSION) into $(PROTO_TOOLS_DIR)/"
	@rm -f $(PROTO_TOOLS_DIR)/protoc-gen-grpc-java $(PROTO_TOOLS_DIR)/.stamp.protoc-gen-grpc-java.*
	@mkdir -p $(PROTO_TOOLS_DIR)
	@curl -fsSL -o $(PROTO_TOOLS_DIR)/protoc-gen-grpc-java \
		https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/$(PROTOC_GEN_GRPC_JAVA_VERSION)/protoc-gen-grpc-java-$(PROTOC_GEN_GRPC_JAVA_VERSION)-$(PROTO_TOOLS_OS)-$(PROTO_TOOLS_ARCH).exe
	@chmod +x $(PROTO_TOOLS_DIR)/protoc-gen-grpc-java
	@touch $@

$(PROTOC_GEN_GO_STAMP):
	@echo "installing protoc-gen-go $(PROTOC_GEN_GO_VERSION) into $(PROTO_TOOLS_DIR)/"
	@rm -f $(PROTO_TOOLS_DIR)/.stamp.protoc-gen-go.*
	@mkdir -p $(PROTO_TOOLS_DIR)
	@GOBIN=$(abspath $(PROTO_TOOLS_DIR)) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	@touch $@

$(PROTOC_GEN_GO_GRPC_STAMP):
	@echo "installing protoc-gen-go-grpc $(PROTOC_GEN_GO_GRPC_VERSION) into $(PROTO_TOOLS_DIR)/"
	@rm -f $(PROTO_TOOLS_DIR)/.stamp.protoc-gen-go-grpc.*
	@mkdir -p $(PROTO_TOOLS_DIR)
	@GOBIN=$(abspath $(PROTO_TOOLS_DIR)) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	@touch $@

.PHONY: proto-tools
proto-tools: $(PROTOC_STAMP) $(GRPC_JAVA_STAMP) $(PROTOC_GEN_GO_STAMP) $(PROTOC_GEN_GO_GRPC_STAMP)

.PHONY: build
build: lint fmt proto-tools
	rm -rf generated/stubs
	mkdir -p generated/stubs
	buf generate
	cp -R generated/stubs/go/github.com/plantonhq/planton/apis/. .
	rm -rf generated/stubs/go
	# Generate a BUILD.bazel for the Java stubs so Bazel can compile them.
	# This serves as a build-time gate: if any proto generates invalid Java
	# (e.g., reserved keyword conflicts), the build fails here instead of in
	# downstream consumers.
	@printf '# gazelle:ignore\njava_library(\n    name = "java",\n    srcs = glob(["**/*.java"]),\n    visibility = ["//visibility:public"],\n    deps = [\n        "@maven//:build_buf_protovalidate",\n        "@maven//:com_google_guava_guava",\n        "@maven//:com_google_protobuf_protobuf_java",\n        "@maven//:io_grpc_grpc_api",\n        "@maven//:io_grpc_grpc_protobuf",\n        "@maven//:io_grpc_grpc_stub",\n        "@maven//:javax_annotation_javax_annotation_api",\n    ],\n)\n' > generated/stubs/java/BUILD.bazel

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

.PHONY: buf-breaking
# Compares the module's protos against the main branch. The rule set lives in
# buf.yaml (breaking:). While every kind sits in the pre-launch alpha posture
# this check is advisory — CI surfaces findings as warnings; it becomes
# channel-aware blocking once versions declare beta/stable channels.
buf-breaking:
	buf breaking --against '../.git#branch=main,subdir=apis'
