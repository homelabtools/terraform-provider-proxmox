BOX_DIR  := vagrant-box
TESTBED  := testbed
GIT_ROOT := $(shell git rev-parse --show-toplevel)
MODULE   := $(shell awk 'NR==1{print $$2}' $(GIT_ROOT)/go.mod)
LDFLAGS  := "-X $(MODULE)/proxmox.disableHTTPSCheck=true"

# NOTE: the endpoint here is just 127.0.0.1:80, which is the port *inside* the guest machine.
# This is because traffic will go through mitmproxy inside the VM, which will then pass it along
# to that endpoint internally within the VM.
export PROXMOX_ENDPOINT ?= http://127.0.0.1:80
export HTTP_PROXY       ?= http://127.0.0.1:58080
export PROVIDER         ?= virtualbox

default: test

vagrant-box:
	cd $(BOX_DIR) && time make

up: vagrant-box
	time vagrant up --provider=$(PROVIDER)

down:
	time vagrant down

destroy:
	time vagrant destroy -f

B:
	echo $(ABC)

test:
	go test -ldflags $(LDFLAGS) -v .

# Same as test but skip the cleanup for debugging purposes
test-debug: SKIP_CLEANUP = 1
test-debug: test

shell: up
	vagrant ssh

clean:
	rm -rf $(TESTBED)
	vagrant destroy -f || true
	rm -rf .vagrant

clean-box:
	cd $(BOX_DIR) && make clean

clean-all: clean clean-box
