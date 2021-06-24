PROVIDER ?= virtualbox
BOX_DIR  := vagrant-box

default: test

vagrant-box:
	cd $(BOX_DIR) && time make

up: vagrant-box
	time vagrant up --provider=$(PROVIDER)

down:
	time vagrant down

destroy:
	time vagrant destroy -f

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