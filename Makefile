.DEFAULT_GOAL := none


test-ci: ci-setup
	go get -v -t -d ./pkg/..
	go test ./pkg/...

test:
	go test ./pkg/...

ci-setup:
	bash scripts/test.sh start
