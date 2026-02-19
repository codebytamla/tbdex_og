SHELL := /bin/sh

.PHONY: backend dashboard dev

backend:
	$(MAKE) -C backend run

dashboard:
	npm --prefix dashboard run dev

dev:
	( $(MAKE) backend ) & ( $(MAKE) dashboard ) & wait
