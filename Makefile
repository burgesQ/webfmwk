NAME	= webfmwk

# can be go, go1.15.15, docker run -v `pwd`:/app -w app go:1.18 ...
GO_CC	= go

.PHONY: all
all: test lint ## Run the test and the linter

include ./make_src/align.mk
include ./make_src/lint.mk
include ./make_src/test.mk
include ./make_src/license.mk
include ./make_src/help.mk
