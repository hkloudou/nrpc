.PHONY: default
git:
	-cd $(shell git rev-parse --show-toplevel) && git autotag -commit 'auto commit' -t -i -f -p
	@echo current version:`git describe`