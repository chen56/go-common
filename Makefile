#
# Makefile
# test

.PHONY: test
test:
	go test ./...

.PHONY: update
update:
	rm -rf vendor
	go mod tidy -v
	go get -u=patch -v