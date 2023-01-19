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

.PHONY: ssh
SSH_FILE ?= ~/Downloads/eisucon/eisucon.pem
IP	?=	13.208.215.218
ssh:
	@ssh -i ${SSH_FILE} ec2-user@${IP}

.PYONY: log-save
log-save: /home/ec2-user/benchmark_logs
	cd ./backend ; \
		docker-compose logs backend | \
			sed 's/backend-backend-1  | //g' | \
			grep 'time:' \
			> /home/ec2-user/benchmark_logs/$$(date +%s).log

.PYONY: log-dl
SSH_FILE ?= ~/Downloads/eisucon/eisucon.pem
IP	?=	13.208.215.218
log-dl:
	@ssh -i ${SSH_FILE} ec2-user@${IP} "ls -t /home/ec2-user/benchmark_logs/ | head -1" | \
		xargs -I SomeString scp -i ${SSH_FILE} ec2-user@${IP}:/home/ec2-user/benchmark_logs/SomeString ./benchmark_ltsv.log

.PYONY: alp
alp:
	cat benchmark_ltsv.log | \
		alp ltsv \
			-m '/users/[0-9]+$$,/users/[0-9]+/star,/events/[0-9]+$$,/events/[0-9]+/documents$$,/events/[0-9]+/documents/[0-9]+' \
			-r --sort avg
