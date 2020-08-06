.DEFAULT_GOAL := none


test-ci: ci-setup
	go get -v -t -d ./..
	go test ./pkg/...

test: ci-setup
	go test ./pkg/...

ci-setup:
	bash scripts/test.sh start
