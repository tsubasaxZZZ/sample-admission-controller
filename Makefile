build:
	@go build

image/build:
	docker build -t tokibi/sample-admission-controller:latest .

image/push:
	docker push tokibi/sample-admission-controller

image/update: image/build image/push
