##
#
# @file
# @version 1.0
build: config
	go build goblocks.go

config:
	./checkConfig

install: build
	sudo cp -f goblocks /usr/local/bin/goblocks; \
	sudo chmod 755 /usr/local/bin/goblocks

uninstall:
	sudo rm -f /usr/local/bin/goblocks
	rm -f ${HOME}/.config/goblocks.json

run: build
	./goblocks
