.PHONY : update updatetools tools
.PHONY : lint test integration coverage 
.PHONY : build
.PHONY : clean cleancoverage cleantools

PROJECT_PATH = $(shell pwd -L)
TOOLSDIR = $(PROJECT_PATH)/tools
BINDIR = $(PROJECT_PATH)/bin
BUILDDIR = $(PROJECT_PATH)/.build
GOFLAGS ::= ${GOFLAGS}
COVERDIR = $(PROJECT_PATH)/.coverage
COVEROUT = $(wildcard $(COVERDIR)/*.out)
COVERINTERCHANGE = $(COVEROUT:.out=.interchange)
COVERHTML = $(COVEROUT:.out=.html)
COVERXML = $(COVEROUT:.out=.xml)
COVERCOMBINED ::= $(COVERDIR)/combined.out
BUILDNAME = remouse

# Tools need to be enumerated here in order to support updating them.
# They will only be used in the context of the /tools directory which
# is a special sub-module that is only engaged when handling the tools.
# We do this to prevent having tools included in the actual project
# dependencies.
TOOLS ::= github.com/golang/mock/mockgen
TOOLS ::= $(TOOLS) golang.org/x/tools/cmd/goimports
TOOLS ::= $(TOOLS) github.com/golangci/golangci-lint/cmd/golangci-lint
TOOLS ::= $(TOOLS) github.com/axw/gocov/gocov
TOOLS ::= $(TOOLS) github.com/matm/gocov-html
TOOLS ::= $(TOOLS) github.com/AlekSi/gocov-xml
TOOLS ::= $(TOOLS) github.com/wadey/gocovmerge

UNIT_PKGS = $(shell go list ./... | sed 1d | paste -sd "," -)

update:
	GO111MODULE=on go get -u

updatetools: cleantools
	# Regenerate the module files for the tools and then
	# reinstall them. This should be done periodically but
	# infrequently.
	cd $(TOOLSDIR) && GO111MODULE=on go get -u $(TOOLS)
	$(MAKE) $(BINDIR)

$(BINDIR):
	cd $(TOOLSDIR) && GOBIN=$(BINDIR) go install $(TOOLS)

tools: $(BINDIR)
	# This is an alias for generating the tools. Unless
	# $(BINDIR) is set elsewhere it will generate a local
	# /bin directory in the repo.

fmt: $(BINDIR)
	# Apply goimports to all code files. Here we intentionally
	# ignore everything in /vendor if it is present.
	GO111MODULE=on \
	GOFLAGS="$(GOFLAGS)" \
	$(BINDIR)/goimports -w -v \
	-local github.com/kevinconway/ \
	$(shell find . -type f -name '*.go' -not -path "./vendor/*")

lint: $(BINDIR)
	GO111MODULE=on \
	GOFLAGS="$(GOFLAGS)" \
	$(BINDIR)/golangci-lint run \
		--config .golangci.yaml \
		--print-resources-usage \
		--verbose

test: $(BINDIR) $(COVERDIR)
	GO111MODULE=on \
	GOFLAGS="$(GOFLAGS)" \
	go test \
		-v \
		-cover \
		-race \
		-coverpkg="$(UNIT_PKGS)" \
		-coverprofile="$(COVERDIR)/unit.out" \
		./...

integration: $(BINDIR) $(COVERDIR)
	# Integration tests leverage docker-compose to manage
	# test depenencies and state isolation. The _runintegration
	# command is always executed from within the container which
	# will use the 'go test' command with build tags 'integration'
	# targeted at the local /tests directory.
	DIR=$(PROJECT_PATH) \
	docker-compose \
		-f docker-compose.integration.yaml \
		up \
			--abort-on-container-exit \
			--build \
			--exit-code-from test

_runintegration:
	# This is a "private" rule that is usually used within
	# the integration testing container. If you run it
	# directly then you must have already orchestrated any
	# test dependencies and set any relevant environment
	# variables.
	GO111MODULE=on \
	GOFLAGS="$(GOFLAGS)" \
	go test \
		-v \
		-tags=integration \
		-cover \
		-race \
		-coverpkg="$(UNIT_PKGS)" \
		-coverprofile="$(COVERDIR)/integration.out" \
		./tests

build: $(BUILDDIR)
	# Optionally build the service if it has an executable
	# present in the project root.
	GO111MODULE=on \
	GOFLAGS="$(GOFLAGS)" \
	go build -o $(BUILDDIR)/$(BUILDNAME) main.go

$(BUILDDIR):
	mkdir -p $(BUILDDIR)

generate: $(BINDIR)
	# Run any code generation steps.
	GO111MODULE=on \
	GOFLAGS="$(GOFLAGS)" \
	PATH="${PATH}:$(BINDIR)" \
	go generate github.com/kevinconway/remouse/pkg github.com/kevinconway/remouse/pkg/internal
	$(MAKE) fmt

coverage: $(BINDIR) $(COVERDIR) $(COVERCOMBINED) $(COVERINTERCHANGE) $(COVERHTML) $(COVERXML)
	# The cover rule is an alias for a number of other rules that each
	# generate part of the full coverage report. First, any coverage reports
	# are combined so that there is a report both for an individual test run
	# and a report that covers all test runs together. Then all coverage
	# files are converted to an interchange format. From there we generate
	# an HTML and XML report. XML reports may be used with jUnit style parsers,
	# the HTML report is for human consumption in order to help identify
	# the location of coverage gaps, and the original reports are available
	# for any purpose.
	GO111MODULE=on \
	GOFLAGS="$(GOFLAGS)" \
	go tool cover -func $(COVERCOMBINED)

$(COVERCOMBINED):
	GO111MODULE=on \
	GOFLAGS="$(GOFLAGS)" \
 	$(BINDIR)/gocovmerge $(COVERDIR)/*.out > $(COVERCOMBINED)

	# NOTE: I couldn't figure out how to automatically include
	# the combined files with the list of other .out files that
	# are processed in bulk. For now, this needs to have specific
	# calls to make for combined coverage.
	$(MAKE) $(COVERCOMBINED:.out=.interchange)
	$(MAKE) $(COVERCOMBINED:.out=.xml)
	$(MAKE) $(COVERCOMBINED:.out=.html)

$(COVERDIR)/%.interchange: $(COVERDIR)/%.out
	GO111MODULE=on \
	GOFLAGS="$(GOFLAGS)" \
	$(BINDIR)/gocov convert $< > $@

$(COVERDIR)/%.xml: $(COVERDIR)/%.interchange
	cat $< | \
	GO111MODULE=on \
	GOFLAGS="$(GOFLAGS)" \
	$(BINDIR)/gocov-xml > $@

$(COVERDIR)/%.html: $(COVERDIR)/%.interchange
	cat $< | \
	GO111MODULE=on \
	GOFLAGS="$(GOFLAGS)" \
	$(BINDIR)/gocov-html > $@

$(COVERDIR): 
	mkdir -p $(COVERDIR)

clean: cleancoverage cleantools ;

cleantools:
	rm -rf $(BINDIR)

cleancoverage:
	rm -rf $(COVERDIR)
