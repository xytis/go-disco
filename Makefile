BIN_ROOT=bin
BINARY=${BIN_ROOT}/disco

build:
	go build -o ${BINARY}

.terraform:
	terraform init

env-up: .terraform
	terraform apply -auto-approve

env-down:
	-terraform destroy -auto-approve

run: build
	#env DISCOVERY_BACKEND="dns:127.0.0.1:53?suffix=service.consul" ./${BINARY} consul
	#export CONSUL_HTTP_ADDR=http://???
	env DISCOVERY_BACKEND="consul:default" ./${BINARY} consul

clean: env-down
	rm ${BINARY}
