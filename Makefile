build:
	docker build -t pager .
.PHONY: build

test:
	docker run -it --rm --name test-pager pager
.PHONY: test
