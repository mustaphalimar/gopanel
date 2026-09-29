include .env
MIGRATIONS_PATH = ./internal/database/migrations


.PHONY: migrations\:new
migrations\:new:
	@if [ -z "$(name)" ]; then \
		echo "usage: make migrations:new name=migrations_name"; \
		exit 1; \
	fi
	@tern new -m $(MIGRATIONS_PATH) $(name)

.PHONY: migrations\:up
migrations\:up:
	@tern migrate -m $(MIGRATIONS_PATH) --conn-string $(DATABASE_URL)


.PHONY: migrations\:down
migrations\:down:
	@tern migrate --destination -1 -m $(MIGRATIONS_PATH) --conn-string $(DATABASE_URL)



.PHONY: migrations\:reset
migrations\:reset:
	@tern migrate --destination 0 -m $(MIGRATIONS_PATH) --conn-string $(DATABASE_URL)
