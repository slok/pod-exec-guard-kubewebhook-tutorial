package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	kwhlog "github.com/slok/kubewebhook/v2/pkg/log"
	kwhlogrus "github.com/slok/kubewebhook/v2/pkg/log/logrus"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhvalidating "github.com/slok/kubewebhook/v2/pkg/webhook/validating"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Logger interface {
	kwhlog.Logger
}

type config struct {
	listenAddress string
	certFile      string
	keyFile       string
	deleteAfter   string
}

func initFlags() *config {
	cfg := &config{}

	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.StringVar(&cfg.certFile, "tls-cert-file", "", "TLS certificate file")
	fl.StringVar(&cfg.keyFile, "tls-key-file", "", "TLS key file")
	fl.StringVar(&cfg.listenAddress, "listen-address", ":8080", "Listen address")
	fl.StringVar(&cfg.deleteAfter, "delete-after", "1h", "The duration that will be used to delete the drifted pods")
	_ = fl.Parse(os.Args[1:])
	return cfg
}

type podExecValidator struct {
	deleteAfter time.Duration
	k8sCli      kubernetes.Interface
	logger      Logger
}

func (p podExecValidator) Validate(ctx context.Context, ar *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhvalidating.ValidatorResult, error) {
	name := ar.Name
	ns := ar.Namespace
	if ns == "" {
		ns = "default"
	}

	if ar.DryRun {
		p.logger.Debugf("Ignoring %s/%s because is dry-run", ns, name)
		return &kwhvalidating.ValidatorResult{Valid: true}, nil
	}

	podExecOpts, ok := obj.(*unstructured.Unstructured)
	if !ok {
		p.logger.Warningf("Not the type expected, got %T", obj)
		return &kwhvalidating.ValidatorResult{Valid: true}, nil
	}

	kind, ok := podExecOpts.Object["kind"].(string)
	if !ok || kind != "PodExecOptions" {
		p.logger.Warningf("Not the kind expected, got %q", kind)
		return &kwhvalidating.ValidatorResult{Valid: true}, nil
	}

	// Get the pod
	pod, err := p.k8sCli.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get pod: %w", err)
	}
	pod = pod.DeepCopy()

	// Mark as drifted.
	now := time.Now().UTC()
	if pod.Labels == nil {
		pod.Labels = map[string]string{}
	}
	pod.Labels["pod-exec-guard.slok.dev/drift"] = fmt.Sprintf("%d", now.Unix())

	// Set TTL for deletion with `kube-janitor`.
	deleteAfterTS := now.Add(p.deleteAfter)
	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	pod.Annotations["janitor/expires"] = deleteAfterTS.Format(time.RFC3339)

	// Update the pod
	_, err = p.k8sCli.CoreV1().Pods(pod.Namespace).Update(ctx, pod, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not update pod: %w", err)
	}

	p.logger.Infof("Pod '%s/%s' marked as drifted and will be deleted at '%s'", ns, name, deleteAfterTS)

	return &kwhvalidating.ValidatorResult{Valid: true}, nil
}

func run(ctx context.Context) error {
	logrusLogEntry := logrus.NewEntry(logrus.New())
	logrusLogEntry.Logger.SetLevel(logrus.DebugLevel)
	logger := kwhlogrus.NewLogrus(logrusLogEntry)

	cfg := initFlags()

	deleteAfter, err := time.ParseDuration(cfg.deleteAfter)
	if err != nil {
		return fmt.Errorf("invalid delete after duration: %q", deleteAfter)
	}

	// Create Kubernetes client.
	kubeHome := filepath.Join(homedir.HomeDir(), ".kube", "config")
	k8sCfg, err := clientcmd.BuildConfigFromFlags("", kubeHome)
	if err != nil {
		return fmt.Errorf("could not load Kubernetes configuration: %w", err)
	}
	k8sCli, err := kubernetes.NewForConfig(k8sCfg)
	if err != nil {
		return fmt.Errorf("could not create Kubernetes client: %w", err)
	}

	// Create validator.
	validator := podExecValidator{
		deleteAfter: deleteAfter,
		logger:      logger,
		k8sCli:      k8sCli,
	}

	// Create webhook.
	wcfg := kwhvalidating.WebhookConfig{
		ID:        "podExecGuard",
		Validator: validator,
		Logger:    logger,
	}
	wh, err := kwhvalidating.NewWebhook(wcfg)
	if err != nil {
		return fmt.Errorf("could not create webhook: %w", err)
	}

	// Create HTTP handler.
	whHandler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{Webhook: wh, Logger: logger})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating webhook handler: %s", err)
		os.Exit(1)
	}

	// Serve our webhook.
	logger.Infof("Listening on: %s", cfg.listenAddress)
	err = http.ListenAndServeTLS(cfg.listenAddress, cfg.certFile, cfg.keyFile, whHandler)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error serving webhook: %s", err)
		os.Exit(1)
	}

	logger.Infof("hello webhook world")
	return nil
}

func main() {
	ctx := context.Background()
	err := run(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}
}
