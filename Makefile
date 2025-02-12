include .env

dbString := postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)

up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(dbString) goose -dir=$(MIGRATION_PATH) up

reset:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(dbString) goose -dir=$(SQL_FILES_PATH) reset

db-setup:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(dbString) goose -dir=$(SQL_FILES_PATH) up

db-status:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(dbString) goose -dir=$(MIGRATION_PATH) status

db-connect:
	psql $(dbString)
