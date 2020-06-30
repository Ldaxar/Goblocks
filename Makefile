##
# Playground
#
# @file
# @version 0.1
run: build
	./goblocks

build: goblocks.go
	go build goblocks.go


# end
