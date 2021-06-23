PROVIDER      ?= virtualbox
BUILD_LOG_DIR := $(shell mkdir -p .buildlogs && printf .buildlogs)
LOGGED = 2>&1 | tee $(BUILD_LOG_DIR)/target_$@_$(shell date | tr ' ' '_')

default: test

up:
	vagrant up --provider=$(PROVIDER) $(LOGGED)

down:
	vagrant down

test: up

shell: up
	vagrant ssh

clean:
	vagrant destroy -f

clean-logs:
	rm -rf $(BUILD_LOG_DIR)