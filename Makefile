SOURCES := $(wildcard *.go)

.PHONY: clean execute

default: build push clean execute
TEST_DIR := ~
TEST_DST := root@demo.lan
OP_DIR := "~/dns"
OP_DST := root@fakeip.lan
DIR := $(OP_DIR)
DST := $(OP_DST)
SCP := $(DST):$(DIR)
build: $(SOURCES)
	go env -w GOOS=linux
	go build -buildvcs=false .
	go env -w GOOS=windows
push: build
	scp fakeip-proxy $(SCP)
	dos2unix iptables.sh
	scp -r config.yaml iptables.sh v2geo/geosite.dat $(SCP)
execute:
	#ssh $(DST) "chmod 755 ~/dns/fakeip-proxy ~/dns/geosite.dat"
	ssh $(DST) "~/dns/fakeip-proxy dns --config ~/dns/config.yaml"
clean:
	rm fakeip-proxy