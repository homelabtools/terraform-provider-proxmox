PROVIDER ?= virtualbox
BOX_DIR  := vagrant-box

default: test

up:
	#cd $(BOX_DIR) && make
	vagrant up --provider=$(PROVIDER)

down:
	vagrant down

destroy:
	vagrant destroy -f

test: up

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