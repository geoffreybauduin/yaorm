ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

test:
	go test -v -timeout=1200s `go list ./... | grep -v vendor` $(filter-out $@,$(MAKECMDGOALS))

vendor:
	mv _vendor vendor ; govendor add $(filter-out $@,$(MAKECMDGOALS)) ; mv vendor _vendor

testwithcoverage:
	ROOT_DIR=$(ROOT_DIR) ./test-coverage.sh $(filter-out $@,$(MAKECMDGOALS))