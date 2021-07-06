PROVIDER ?= virtualbox
BOX_DIR  := vagrant-box

TIME := /usr/bin/env time -f "Time elapsed %E from command: %C"

default: test

vagrant-box:
	cd $(BOX_DIR) && $(TIME) make

up: vagrant-box
	$(TIME) vagrant up --provider=$(PROVIDER)

down:
	$(TIME) vagrant down

destroy:
	$(TIME) vagrant destroy -f

test: up
	go test -v .

shell: up
	vagrant ssh

clean:
	vagrant destroy -f || true
	rm -rf .vagrant

clean-box:
	cd $(BOX_DIR) && make clean

clean-logs:
	rm -rf $(BUILD_LOG_DIR)
	cd $(BOX_DIR) && make clean-logs

clean-all: clean clean-box clean-logs