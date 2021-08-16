PROVIDER         ?= virtualbox
BOX_NAME         ?= proxmox
BOX_PATH         ?= dist/$(BOX_NAME).box
BOX_SRC          := Vagrantfile $(shell find src/)
export HOST_PORT ?= 51823

default: box

box: check $(BOX_PATH)

$(BOX_PATH): $(BOX_SRC)
	mkdir -p $(dir $(BOX_PATH))
	vagrant plugin install --local
	vagrant destroy -f || true
	time vagrant up --provider=$(PROVIDER)
	rm -f $(BOX_PATH)
	time vagrant package --output $(BOX_PATH)
	time vagrant box add $(BOX_PATH) --name $(BOX_NAME) --force
	# Clean up VM
	time vagrant destroy -f

check:
	shellcheck --severity=style $(wildcard src/*.sh)

clean:
	vagrant box remove $(BOX_NAME) --force --all || true
	rm -rf .vagrant dist
