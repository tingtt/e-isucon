.PHONY: ps
ps:
	cd ./backend ; \
		docker-compose ps

.PHONY: build
build:
	cd ./backend ; \
		docker-compose build

.PHONY: up
up:
	cd ./backend ; \
		docker-compose up -d

.PHONY: down
down:
	cd ./backend ; \
		docker-compose down

.PHONY: purge
purge:
	cd ./backend ; \
		docker-compose down --volumes

.PHONY: upgrade
upgrade: purge build up