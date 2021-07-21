PROVIDER ?= virtualbox
BOX_DIR  := vagrant-box
TESTBED  := testbed

default: test

vagrant-box:
	cd $(BOX_DIR) && time make

up: vagrant-box
	time vagrant up --provider=$(PROVIDER)

down:
	time vagrant down

destroy:
	time vagrant destroy -f

test:
	go test -v .

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
