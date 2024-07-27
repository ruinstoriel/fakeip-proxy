SOURCES := $(wildcard *.go)

.PHONY: clean execute

default: build push clean execute
ETC_DIR := /etc/me
BIN_DIR := /usr/bin
NAME := fakeip-proxy
DST := root@fakeip.lan
ETC := $(DST):$(ETC_DIR)
BIN := $(DST):$(BIN_DIR)
DOS := iptables.sh install.sh  me
ATTACH := config.yaml iptables.sh install.sh geosite.dat
build: $(SOURCES)
	go env -w GOOS=linux
	go build -buildvcs=false -o $(NAME) .
	go env -w GOOS=windows
push: build
	ssh $(DST) "(rc-status | grep  me && rc-service me stop) || true"
	ssh $(DST) "rm $(BIN_DIR)/$(NAME)"
	scp $(NAME) $(BIN)
	ssh $(DST) "chmod  755 $(BIN_DIR)/$(NAME)"
	dos2unix $(DOS)
	scp me $(DST):/etc/init.d
	scp -r $(ATTACH) $(ETC)
	ssh $(DST) "chmod -R 755 $(ETC_DIR)"

execute:
	ssh $(DST) "rc-service me start && rc-update add me"
clean:
	rm $(NAME)