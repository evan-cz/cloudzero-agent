package kms

import (
	"context"
	"fmt"
	"io/ioutil"
	net "net/http"
	"strings"
	"time"

	"github.com/cloudzero/cloudzero-agent-validator/pkg/config"
	"github.com/cloudzero/cloudzero-agent-validator/pkg/diagnostic"
	"github.com/cloudzero/cloudzero-agent-validator/pkg/logging"
	"github.com/cloudzero/cloudzero-agent-validator/pkg/status"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const DiagnosticKMS = config.DiagnosticKMS

var (
	// Exported so that it can be overridden in tests
	MaxRetry      = 12
	RetryInterval = 10 * time.Second
)

type checker struct {
	cfg       *config.Settings
	logger    *logrus.Entry
	clientset *kubernetes.Clientset
}

func NewProvider(ctx context.Context, cfg *config.Settings) diagnostic.Provider {
	// Use the in-cluster config if running inside a cluster, otherwise use the default kubeconfig
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := clientcmd.RecommendedHomeFile
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &checker{
		cfg: cfg,
		logger: logging.NewLogger().
			WithContext(ctx).WithField(logging.OpField, "ksm"),
		clientset: clientset,
	}
}

func (c *checker) Check(ctx context.Context, client *net.Client, accessor status.Accessor, cfg *config.Settings) error {
	var (
		retriesRemaining = MaxRetry
		namespace        = "prom-agent"
		serviceName      = "cz-prom-agent-kube-state-metrics"
		endpointURL      string
	)

	// Wait for the pod to become ready and find the first available endpoint
	for retriesRemaining > 0 {
		endpoints, err := c.clientset.CoreV1().Endpoints(namespace).Get(ctx, serviceName, metav1.GetOptions{})
		if err != nil {
			c.logger.Errorf("Failed to get service endpoints: %v", err)
			accessor.AddCheck(&status.StatusCheck{Name: DiagnosticKMS, Passing: false, Error: fmt.Sprintf("Failed to get service endpoints: %s", err.Error())})
			return nil
		}

		// Log the endpoints for debugging
		c.logger.Infof("Endpoints: %v", endpoints)

		// Check if there are any ready addresses and find the first available endpoint
		for _, subset := range endpoints.Subsets {
			for _, address := range subset.Addresses {
				c.logger.Infof("Address: %v", address)
				for _, port := range subset.Ports {
					c.logger.Infof("Port: %v", port)
					if port.Port == 8080 {
						endpointURL = fmt.Sprintf("http://%s:%d/metrics", address.IP, port.Port)
						break
					}
				}
				if endpointURL != "" {
					break
				}
			}
			if endpointURL != "" {
				break
			}
		}

		if endpointURL != "" {
			break
		}

		c.logger.Infof("Pod is not ready, waiting...")
		retriesRemaining--
		time.Sleep(RetryInterval)
	}

	if retriesRemaining == 0 {
		c.logger.Errorf("Pod did not become ready in time")
		accessor.AddCheck(&status.StatusCheck{Name: DiagnosticKMS, Passing: false, Error: "Pod did not become ready in time"})
		return nil
	}

	c.logger.Infof("Using endpoint URL: %s", endpointURL)

	// Retry logic to handle transient issues
	retriesRemaining = MaxRetry
	for retriesRemaining > 0 {
		resp, err := client.Get(endpointURL)
		if err == nil && resp.StatusCode == net.StatusOK {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				c.logger.Errorf("Failed to read metrics: %v", err)
				accessor.AddCheck(&status.StatusCheck{Name: DiagnosticKMS, Passing: false, Error: fmt.Sprintf("Failed to read metrics: %s", err.Error())})
				return nil
			}

			metrics := string(body)
			requiredMetrics := []string{"kube_pod_info", "kube_node_info"} // Add the required metrics here
			allMetricsFound := true
			for _, metric := range requiredMetrics {
				if !strings.Contains(metrics, metric) {
					c.logger.Errorf("Required metric %s not found", metric)
					accessor.AddCheck(&status.StatusCheck{Name: DiagnosticKMS, Passing: false, Error: fmt.Sprintf("Required metric %s not found", metric)})
					allMetricsFound = false
				}
			}

			if allMetricsFound {
				c.logger.Infof("All required metrics found: %v", requiredMetrics)
				accessor.AddCheck(&status.StatusCheck{Name: DiagnosticKMS, Passing: true})
			}
			return nil
		}

		c.logger.Errorf("Failed to fetch metrics: %v", err)
		retriesRemaining--
		time.Sleep(RetryInterval)
	}

	accessor.AddCheck(&status.StatusCheck{Name: DiagnosticKMS, Passing: false, Error: fmt.Sprintf("Failed to fetch metrics after %d retries", MaxRetry)})
	return nil
}
