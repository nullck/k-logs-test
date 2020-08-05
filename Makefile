.DEFAULT_GOAL := none


test: ci-setup
	go test ./pkg/...

ci-setup:
	cd scripts; bash test.sh start; cd ..
