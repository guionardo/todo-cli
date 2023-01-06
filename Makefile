test:
	go test -v ./... -cover

build_dist:
	bash -x ./.github/scripts/build_dist.sh

install:
	bash -x ./.github/scripts/build_dist.sh
