build:
	@go build

image/build:
	docker build -t tokibi/sample-admission-controller:latest .

image/push:
	docker push tokibi/sample-admission-controller

image/update: image/build image/push

cert/generate:
	@cd deploy && \
	cfssl gencert -initca ca-csr.json | cfssljson -bare ca - && \
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=server server-csr.json | cfssljson -bare server

deploy:

