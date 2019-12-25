package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	whhttp "github.com/slok/kubewebhook/pkg/http"
	"github.com/slok/kubewebhook/pkg/log"
	validatingwh "github.com/slok/kubewebhook/pkg/webhook/validating"
)

type podSpecValidator struct {
	logger log.Logger
}

func (v *podSpecValidator) Validate(_ context.Context, obj metav1.Object) (bool, validatingwh.ValidatorResult, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return false, validatingwh.ValidatorResult{}, fmt.Errorf("not a pod")
	}

	for _, vol := range pod.Spec.Volumes {
		if vol.HostPath != nil && vol.HostPath.Path != "" {
			v.logger.Infof("denied")
			res := validatingwh.ValidatorResult{
				Valid:   false,
				Message: fmt.Sprintf(`pod "%s" creation denied, hostPath is not allowed`, pod.Name),
			}
			return false, res, nil
		}
	}

	v.logger.Infof("accepted")
	res := validatingwh.ValidatorResult{
		Valid:   true,
		Message: "accepted",
	}
	return false, res, nil
}

type config struct {
	certFile  string
	keyFile   string
	hostRegex string
	addr      string
}

func initFlags() *config {
	cfg := &config{}

	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.StringVar(&cfg.certFile, "tls-cert-file", "", "TLS certificate file")
	fl.StringVar(&cfg.keyFile, "tls-key-file", "", "TLS key file")
	fl.StringVar(&cfg.addr, "listen-addr", ":8080", "The address to start the server")

	fl.Parse(os.Args[1:])
	return cfg
}

func main() {
	logger := &log.Std{Debug: true}

	cfg := initFlags()

	vl := &podSpecValidator{
		logger: logger,
	}

	vcfg := validatingwh.WebhookConfig{
		Name: "podSpecValidator",
		Obj:  &corev1.Pod{},
	}
	wh, err := validatingwh.NewWebhook(vcfg, vl, nil, nil, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating webhook: %s", err)
		os.Exit(1)
	}

	// Serve the webhook.
	logger.Infof("Listening on %s", cfg.addr)
	err = http.ListenAndServeTLS(cfg.addr, cfg.certFile, cfg.keyFile, whhttp.MustHandlerFor(wh))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error serving webhook: %s", err)
		os.Exit(1)
	}
}
