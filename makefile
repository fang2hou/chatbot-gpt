EMPTY:=

cmd_dir:=cmd
bin_dir:=bin

define cross_build
	@sh scripts/build.sh ${1} windows arm64
	@sh scripts/build.sh ${1} windows amd64
	@sh scripts/build.sh ${1} linux arm64
	@sh scripts/build.sh ${1} linux arm 6
	@sh scripts/build.sh ${1} linux arm 7
	@sh scripts/build.sh ${1} linux amd64
	@sh scripts/build.sh ${1} darwin arm64
	@sh scripts/build.sh ${1} darwin amd64
endef

build-all:
	${foreach dir, $(shell ls cmd), @sh scripts/build.sh ${dir}}

cross-build-all:
	${foreach dir, $(shell ls cmd), $(call cross_build, ${dir})}