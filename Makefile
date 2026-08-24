EXERCISE_NAME := $(strip $(if $(N),$(N),$(NAME)))

.PHONY: run r test t new n build b

run r:
	@if [ -n "$(EXERCISE_NAME)" ]; then \
		bash scripts/run-exercise.sh "$(EXERCISE_NAME)" $(if $(INPUT),-in "$(INPUT)") $(if $(OUTPUT),-out "$(OUTPUT)") $(NUMBERS); \
	else \
		go run . $(if $(INPUT),-in "$(INPUT)") $(if $(OUTPUT),-out "$(OUTPUT)") $(NUMBERS); \
	fi

test t:
	go test ./...

new n:
	@bash scripts/new-exercise.sh "$(EXERCISE_NAME)"

build b:
	@bash scripts/build-exercise.sh "$(EXERCISE_NAME)"
