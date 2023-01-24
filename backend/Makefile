.PHONY: oas-code-gen-echo
oas-code-gen-echo:
	oapi-codegen --old-config-style -generate types -package echo presentation/oas.yml > presentation/echo/types.gen.go
	oapi-codegen --old-config-style -generate server -package echo presentation/oas.yml > presentation/echo/server.gen.go

.PHONY: build
build:
	go build

.PHONY: install
SEED_FELE	?= /usr/local/lib/eisucon-backend/seed.sql
MODULE	?= prc_hub_back
BIN	?= /usr/local/bin/eisucon-backend
ARG_LOG_LEVEL	?= debug
ARG_MYSQL_DB	?= prc_hub
ARG_MYSQL_USER	?= root
ARG_MYSQL_PASSWORD	?= secret
ARGS	?= --log.level=$(ARG_LOG_LEVEL) --jwt.issuer=prc_hub --jwt.secret=e0VzhtkQ --mysql.db=$(ARG_MYSQL_DB) --mysql.user=$(ARG_MYSQL_USER) --mysql.password=$(ARG_MYSQL_PASSWORD) --migrate-sql-file=$(SEED_FELE)
define UNITFILE
[Unit]
Description=ECC-ISUCON backend
After=network.target

[Service]
Restart=on-failure
RestartSec=10
ExecStart=$(BIN) $(ARGS)

[Install]
WantedBy=multi-user.target
endef
export UNITFILE
UNITFILE_PATH	?= /usr/local/lib/systemd/system/eisucon-backend.service
install: build
	mkdir -p /usr/local/lib/eisucon-backend/
	mysql -u$(ARG_MYSQL_USER) -p$(ARG_MYSQL_PASSWORD) < ./.mysql/init.sql
	@echo "$$UNITFILE" > ./.tmp.service
	cp -f ./.tmp.service $(UNITFILE_PATH)
	@rm -f ./.tmp.service
	cp -n ./domain/model/eisucon/migrate.sql $(SEED_FELE)
	cp -n ./$(MODULE) $(BIN)
	systemctl daemon-reload
	systemctl enable --now eisucon-backend.service


.PHONY: purge
SEED_FELE	?= /usr/local/lib/eisucon-backend/seed.sql
BIN	?= /usr/local/bin/eisucon-backend
UNITFILE_PATH	?= /usr/local/lib/systemd/system/eisucon-backend.service
ARG_MYSQL_USER	?= root
ARG_MYSQL_PASSWORD	?= secret
ARG_MYSQL_DB	?= prc_hub
define SQL
DROP TABLE IF EXISTS $(ARG_MYSQL_DB).documents;
DROP TABLE IF EXISTS $(ARG_MYSQL_DB).event_datetimes;
DROP TABLE IF EXISTS $(ARG_MYSQL_DB).events;
DROP TABLE IF EXISTS $(ARG_MYSQL_DB).users;
DROP DATABASE IF EXISTS $(ARG_MYSQL_DB);
endef
export SQL
purge:
	systemctl disable --now eisucon-backend.service
	@echo "$$SQL" > ./.tmp.sql
	mysql -u$(ARG_MYSQL_USER) -p$(ARG_MYSQL_PASSWORD) < ./.tmp.sql
	@rm -f ./.tmp.sql
	rm -f $(SEED_FELE) $(BIN) $(UNITFILE_PATH)
	-rmdir /usr/local/lib/eisucon-backend/