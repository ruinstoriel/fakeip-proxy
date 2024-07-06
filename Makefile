SOURCES := $(wildcard *.go)

.PHONY: clean execute

default: build push clean execute
TEST_DIR := ~
TEST_DST := root@demo.lan
OP_DIR := /root/dns
OP_DST := root@fakeip.lan
DIR := $(OP_DIR)
DST := $(OP_DST)
SCP := $(DST):$(DIR)
build: $(SOURCES)
	go env -w GOOS=linux
	go build -buildvcs=false .
	go env -w GOOS=windows
push: build
	ssh $(DST) "pkill fakeip-proxy && echo 'kill progress success'"
	scp fakeip-proxy $(SCP)
	dos2unix iptables.sh
	scp -r config.yaml iptables.sh geosite.dat $(SCP)
execute:
	#ssh $(DST) "chmod 755 $(OP_DIR)/fakeip-proxy"
	ssh $(DST) "rm $(OP_DIR)/dns.log && echo 'delete dns.log success'"
	ssh $(DST) "(cd $(OP_DIR) && ./fakeip-proxy dns --config config.yaml --geosite geosite.dat > dns.log 2>&1 &)"
clean:
	rm fakeip-proxy