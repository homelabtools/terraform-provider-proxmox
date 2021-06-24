PROVIDER         ?= virtualbox
BOX_NAME         ?= proxmox
BOX_PATH         ?= dist/$(BOX_NAME).box
BOX_SRC          := Vagrantfile $(shell find src/)
export HOST_PORT ?= 51823
PROXMOX_URL      := https://127.0.0.1:$(HOST_PORT)

default: box

box: prereqs $(BOX_PATH)

$(BOX_PATH): $(BOX_SRC)
	mkdir -p $(dir $(BOX_PATH))
	vagrant destroy -f || true
	vagrant up --provider=$(PROVIDER)
	# Wait for Proxmox to come up before packaging the box
	timeout 300 sh -cx "while [ $$(curl -kLso /dev/null -w %{http_code} $(PROXMOX_URL)) != 200 ]; do sleep 1; done"
	vagrant package --output $(BOX_PATH)
	vagrant box add $(BOX_PATH) --name $(BOX_NAME) --force
	# Clean up VM
	vagrant destroy -f

check:
	shellcheck --severity=style *.sh

clean:
	vagrant box remove $(BOX_NAME) --force --all || true
	vagrant destroy -f $(VM_NAME) || true
	rm -rf .vagrant dist $(BUILD_LOG_DIR)

clean-logs:
	rm -rf $(BUILD_LOG_DIR)

prereqs:
	vagrant plugin install vagrant-reload