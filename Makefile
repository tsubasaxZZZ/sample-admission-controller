build:
	@go build

image/build:
	docker build -t tsubasaxzzz/sample-admission-controller:latest .

image/push:
	docker push tsubasaxzzz/sample-admission-controller

image/update: image/build image/push

cert/generate:
	@cd certs/ && \
	cfssl gencert -initca ca-csr.json | cfssljson -bare ca - && \
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=server server-csr.json | cfssljson -bare server

	@echo "ca: $(shell openssl enc -base64 -in ./certs/ca.pem | tr -d '\n')"
	@echo "cert: $(shell openssl enc -base64 -in ./certs/server.pem | tr -d '\n')"
	@echo "key: $(shell openssl enc -base64 -in ./certs/server-key.pem | tr -d '\n')"

deploy:
	kubectl apply -f deploy/

.PHONY: build deploy
