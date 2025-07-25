# Obtiene la lista de todos los directorios dentro de la carpeta "api/serverless"
CMD_DIRS := $(shell find ./api -mindepth 2 -type d)

# Obtiene la lista de nombres de los programas a partir de los nombres de las carpetas
PROGRAMS := $(notdir $(CMD_DIRS))

# Flags de compilación (puedes personalizarlos según tus necesidades)
LDFLAGS := -ldflags="-s -w"

# Comando predeterminado: compila todos los programas
all: build

.PHONY: build
build:
	@for dir in $(CMD_DIRS); do \
		echo "Building $$dir..."; \
		GOOS=linux GOARCH=amd64 go build -o bootstrap $$dir/main.go && \
		zip -j function.zip bootstrap && \
		mv function.zip $${dir}/function.zip && \
		rm bootstrap; \
	done

deploy:
	sam validate --lint
	sam deploy --config-env dev 

build-deploy:
	make build
	make deploy
	make clean

# Regla para limpiar los binarios compilados
clean:
	@echo "Cleaning compiled files..."
	@for dir in $(CMD_DIRS); do \
		rm -f $$dir/main; \
		rm -f $$dir/prefill/main; \
		rm -f $$dir/function.zip; \
		rm -f $$dir/bootstrap; \
	done

# Indicar que "clean" no es un archivo de salida
.PHONY: clean

# TODO migrate to taskfile