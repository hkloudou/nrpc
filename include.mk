.PHONY: default
.DEFAULT_GOAL := default
ifneq ($(shell pwd),$(shell git rev-parse --show-toplevel))
	GIT_SUBPATH=$(subst $(shell git rev-parse --show-toplevel)/,,$(shell pwd))
	GIT_SUB_PARAME = -s ${GIT_SUBPATH}
	GIT_CLOSEDVERSION = $(shell git describe --abbrev=0  --match ${GIT_SUBPATH}/v[0-9]*\.[0-9]*\.[0-9]*)
else
	GIT_CLOSEDVERSION = $(shell git describe --abbrev=0  --match v[0-9]*\.[0-9]*\.[0-9]*)
endif
print:
	@echo sub: ${GIT_SUBPATH}
	@echo close: ${GIT_CLOSEDVERSION}
git:
	- git autotag -commit 'auto commit ${GIT_SUBPATH}' -t -f -i -p ${GIT_SUB_PARAME}
	@echo current version:`git describe`