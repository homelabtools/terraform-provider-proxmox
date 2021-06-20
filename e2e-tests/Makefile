PROVIDER ?= virtualbox

up:
	vagrant up --provider=$(PROVIDER)

down:
	vagrant down

test: up

clean:
	vagrant destroy -f