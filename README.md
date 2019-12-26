# sample-admission-controller

## deploy

```shell
$ make deploy
secret/admission-controller-certs created
deployment.apps/admission-controller created
service/admission-controller created
validatingwebhookconfiguration.admissionregistration.k8s.io/validating-webhook created
```

## test

```shell
$ kubectl apply -f example/pod.yaml
pod/valid-pod created
Error from server: error when creating "example/pod.yaml": admission webhook "validating-webhook.example.com" denied the request: pod "invalid-pod" creation denied, hostPath is not allowed
```
