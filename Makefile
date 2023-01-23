.PHONY: ps
SERVICE_NAME	?= eisucon-backend.service
ps:
	systemctl status $(SERVICE_NAME)

.PHONY: build
build:
	cd ./backend ; \
		make build

.PHONY: up
SERVICE_NAME	?= eisucon-backend.service
up:
	systemctl stop $(SERVICE_NAME)

.PHONY: down
SERVICE_NAME	?= eisucon-backend.service
down:
	systemctl stop $(SERVICE_NAME)

.PHONY: purge
purge:
	cd ./backend ; \
		make purge

.PHONY: upgrade
upgrade: purge build up

.PHONY: ssh
SSH_FILE ?= ~/Downloads/eisucon/eisucon.pem
IP	?=	13.208.215.218
ssh:
	@ssh -i ${SSH_FILE} ec2-user@${IP}

.PYONY: log-save
SERVICE_NAME	?= eisucon-backend.service
log-save: /home/ec2-user/benchmark_logs
	journalctl -u $(SERVICE_NAME) --no-pager | \
		cut -b 54- | \
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
			-m '/users/[0-9a-f\-]+$$,/users/[0-9a-f\-]+/star,/events/[0-9a-f\-]+$$,/events/[0-9a-f\-]+/documents$$,/events/[0-9a-f\-]+/documents/[0-9a-f\-]+' \
			-r --sort avg
