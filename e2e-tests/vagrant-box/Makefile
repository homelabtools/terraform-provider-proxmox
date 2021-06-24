include ../../tools/tools.mk

PROVIDER         ?= virtualbox
BOX_NAME         ?= proxmox
BOX_PATH         ?= dist/$(BOX_NAME).box
export HOST_PORT ?= 51823
PROXMOX_URL      := https://127.0.0.1:$(HOST_PORT)

# Changes to any of these files will cause the box to be rebuilt.
# Any new source files should be added to this list.
BOX_SRC  := Vagrantfile install-proxmox.sh install-guest-additions.sh

default: box

box: prereqs $(BOX_PATH)

$(BOX_PATH): $(BOX_SRC)
	mkdir -p $(dir $(BOX_PATH))
	vagrant destroy -f $(LOGGED) || true
	vagrant up --provider=$(PROVIDER) $(LOGGED)
	# Wait for Proxmox to come up before packaging the box
	timeout 300 sh -cx "while [ $$(curl -kLso /dev/null -w %{http_code} $(PROXMOX_URL)) != 200 ]; do sleep 1; done"
	vagrant package --output $(BOX_PATH) $(LOGGED)
	vagrant box add $(BOX_PATH) --name $(BOX_NAME) --force $(LOGGED)
	# Clean up VM
	vagrant destroy -f $(LOGGED)

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