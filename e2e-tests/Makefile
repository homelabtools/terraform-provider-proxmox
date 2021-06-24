include ../tools/tools.mk

PROVIDER ?= virtualbox
BOX_DIR  := vagrant-box

default: test

up:
	cd $(BOX_DIR) && make
	vagrant up --provider=$(PROVIDER)

down:
	vagrant down $(LOGGED)

test: up

shell: up
	vagrant ssh

clean:
	vagrant destroy -f || true
	cd $(BOX_DIR) && make clean
	rm -rf .vagrant

clean-logs:
	rm -rf $(BUILD_LOG_DIR)
	cd $(BOX_DIR) && make clean-logs
