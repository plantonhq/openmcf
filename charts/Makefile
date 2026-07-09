# =============================================================================
# InfraCharts Validate, Build & Publish
# =============================================================================
# Discovers all charts (provider/chart-name directories containing values.yaml)
# and validates/publishes them.
#
# Usage (run from this charts/ directory):
#   make validate       Validate all charts OFFLINE (renders with default
#                       values and validates every manifest and reference
#                       against the protos in this working tree — no control
#                       plane needed)
#   make build          Build (validate) all charts through the platform
#   make build-verbose  Build all charts with full output on failure
#   make publish        Publish all charts to the Planton platform
#
# Two `planton` binaries are in play:
#   - `make validate` uses the OPEN-SOURCE CLI's `chart validate` command and
#     builds it from this working tree (so validation always matches the protos
#     you are editing). Override with VALIDATE_BIN=/path/to/planton to skip the
#     build.
#   - `make build` / `make publish` use the Planton *Platform* CLI: the
#     `InfraChart` kind is a platform API, and server-side chart build renders
#     and validates through the live control plane — the authoritative gate
#     before publishing.
# =============================================================================

CHARTS := $(shell find . -mindepth 2 -name values.yaml -not -path '*/build/*' \
            | sed 's|/values.yaml||' | sed 's|^\./||' | sort)

.DEFAULT_GOAL := build

VALIDATE_BIN ?= /tmp/planton-chart-validate

.PHONY: validate
validate:
	@if [ "$(VALIDATE_BIN)" = "/tmp/planton-chart-validate" ]; then \
	  echo "building planton CLI from the working tree..."; \
	  (cd .. && go build -o $(VALIDATE_BIN) .) || exit 1; \
	fi; \
	passed=0; failed=0; failures=""; i=0; \
	total=$(words $(CHARTS)); \
	for chart in $(CHARTS); do \
	  i=$$((i+1)); \
	  if $(VALIDATE_BIN) chart validate "$$chart" >/dev/null 2>&1; then \
	    passed=$$((passed+1)); \
	    printf "[%d/%d] ✔ %s\n" "$$i" "$$total" "$$chart"; \
	  else \
	    failed=$$((failed+1)); \
	    failures="$$failures  $$chart\n"; \
	    printf "[%d/%d] ✘ %s\n" "$$i" "$$total" "$$chart"; \
	  fi; \
	done; \
	echo ""; \
	echo "Validate complete: $$passed passed, $$failed failed out of $$total charts"; \
	if [ $$failed -gt 0 ]; then \
	  echo ""; echo "Failed charts (re-run '$(VALIDATE_BIN) chart validate <chart>' for details):"; printf "$$failures"; exit 1; \
	fi

.PHONY: build
build:
	@passed=0; failed=0; failures=""; i=0; \
	total=$(words $(CHARTS)); \
	for chart in $(CHARTS); do \
	  i=$$((i+1)); \
	  if planton chart build "$$chart" --no-browser >/dev/null 2>&1; then \
	    passed=$$((passed+1)); \
	    printf "[%d/%d] ✔ %s\n" "$$i" "$$total" "$$chart"; \
	  else \
	    failed=$$((failed+1)); \
	    failures="$$failures  $$chart\n"; \
	    printf "[%d/%d] ✘ %s\n" "$$i" "$$total" "$$chart"; \
	  fi; \
	done; \
	echo ""; \
	echo "Build complete: $$passed passed, $$failed failed out of $$total charts"; \
	if [ $$failed -gt 0 ]; then \
	  echo ""; echo "Failed charts:"; printf "$$failures"; exit 1; \
	fi

.PHONY: build-verbose
build-verbose:
	@passed=0; failed=0; failures=""; \
	total=$(words $(CHARTS)); \
	for chart in $(CHARTS); do \
	  echo "--- $$chart ---"; \
	  if planton chart build "$$chart" --no-browser 2>&1; then \
	    passed=$$((passed+1)); \
	  else \
	    failed=$$((failed+1)); \
	    failures="$$failures  $$chart\n"; \
	  fi; \
	  echo ""; \
	done; \
	echo "Build complete: $$passed passed, $$failed failed out of $$total charts"; \
	if [ $$failed -gt 0 ]; then \
	  echo ""; echo "Failed charts:"; printf "$$failures"; exit 1; \
	fi

.PHONY: publish
publish:
	@passed=0; failed=0; failures=""; i=0; \
	total=$(words $(CHARTS)); \
	for chart in $(CHARTS); do \
	  i=$$((i+1)); \
	  if planton chart publish "$$chart" --platform 2>&1 | grep -q "applied"; then \
	    passed=$$((passed+1)); \
	    printf "[%d/%d] ✔ %s\n" "$$i" "$$total" "$$chart"; \
	  else \
	    failed=$$((failed+1)); \
	    failures="$$failures  $$chart\n"; \
	    printf "[%d/%d] ✘ %s\n" "$$i" "$$total" "$$chart"; \
	  fi; \
	done; \
	echo ""; \
	echo "Publish complete: $$passed passed, $$failed failed out of $$total charts"; \
	if [ $$failed -gt 0 ]; then \
	  echo ""; echo "Failed charts:"; printf "$$failures"; exit 1; \
	fi
