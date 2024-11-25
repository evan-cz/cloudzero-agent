package kms_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/cloudzero/cloudzero-agent-validator/pkg/config"
	"github.com/cloudzero/cloudzero-agent-validator/pkg/diagnostic/kms"
	"github.com/cloudzero/cloudzero-agent-validator/pkg/status"
	"github.com/cloudzero/cloudzero-agent-validator/test"
)

const (
	mockURL = "http://example.com"
)

func makeReport() status.Accessor {
	return status.NewAccessor(&status.ClusterStatus{})
}

func TestChecker_CheckOK(t *testing.T) {
	cfg := &config.Settings{
		Prometheus: config.Prometheus{
			KubeStateMetricsServiceEndpoint: mockURL,
		},
	}
	provider := kms.NewProvider(context.Background(), cfg)

	mock := test.NewHTTPMock()
	mock.Expect(http.MethodGet, "kube_pod_info\nkube_node_info\n", http.StatusOK, nil)
	client := mock.HTTPClient()

	accessor := makeReport()

	err := provider.Check(context.Background(), client, accessor, cfg)
	assert.NoError(t, err)

	accessor.ReadFromReport(func(s *status.ClusterStatus) {
		assert.Len(t, s.Checks, 2)
		for _, c := range s.Checks {
			assert.True(t, c.Passing)
		}
	})
}

func TestChecker_CheckRetry(t *testing.T) {
	cfg := &config.Settings{
		Prometheus: config.Prometheus{
			KubeStateMetricsServiceEndpoint: mockURL,
		},
	}
	provider := kms.NewProvider(context.Background(), cfg)

	// Update the test sleep interval to accelerate the test
	kms.RetryInterval = 10 * time.Millisecond
	kms.MaxRetry = 3

	mock := test.NewHTTPMock()
	for i := 0; i < kms.MaxRetry-1; i++ {
		mock.Expect(http.MethodGet, "", http.StatusNotFound, nil)
	}
	mock.Expect(http.MethodGet, "kube_pod_info\nkube_node_info\n", http.StatusOK, nil)
	client := mock.HTTPClient()

	accessor := makeReport()

	err := provider.Check(context.Background(), client, accessor, cfg)
	assert.NoError(t, err)

	accessor.ReadFromReport(func(s *status.ClusterStatus) {
		assert.Len(t, s.Checks, 2)
		for _, c := range s.Checks {
			assert.True(t, c.Passing)
		}
	})
}

func TestChecker_CheckRetryFailure(t *testing.T) {
	cfg := &config.Settings{
		Prometheus: config.Prometheus{
			KubeStateMetricsServiceEndpoint: mockURL,
		},
	}
	provider := kms.NewProvider(context.Background(), cfg)

	// Update the test sleep interval to accelerate the test
	kms.RetryInterval = 10 * time.Millisecond
	kms.MaxRetry = 3

	mock := test.NewHTTPMock()
	for i := 0; i < kms.MaxRetry; i++ {
		mock.Expect(http.MethodGet, "", http.StatusNotFound, nil)
	}
	client := mock.HTTPClient()

	accessor := makeReport()

	err := provider.Check(context.Background(), client, accessor, cfg)
	assert.NoError(t, err)

	accessor.ReadFromReport(func(s *status.ClusterStatus) {
		assert.Len(t, s.Checks, 1)
		for _, c := range s.Checks {
			assert.False(t, c.Passing)
		}
	})
}

func TestChecker_CheckMetricsValidation(t *testing.T) {
	cfg := &config.Settings{
		Prometheus: config.Prometheus{
			KubeStateMetricsServiceEndpoint: mockURL,
		},
	}
	provider := kms.NewProvider(context.Background(), cfg)

	mock := test.NewHTTPMock()
	mock.Expect(http.MethodGet, "kube_pod_info\nkube_node_info\n", http.StatusOK, nil)
	client := mock.HTTPClient()

	accessor := makeReport()

	err := provider.Check(context.Background(), client, accessor, cfg)
	assert.NoError(t, err)

	accessor.ReadFromReport(func(s *status.ClusterStatus) {
		assert.Len(t, s.Checks, 2)
		for _, c := range s.Checks {
			assert.True(t, c.Passing)
		}
	})
}

func TestChecker_CheckHandles500Error(t *testing.T) {
	cfg := &config.Settings{
		Prometheus: config.Prometheus{
			KubeStateMetricsServiceEndpoint: mockURL,
		},
	}
	provider := kms.NewProvider(context.Background(), cfg)

	mock := test.NewHTTPMock()
	mock.Expect(http.MethodGet, "", http.StatusInternalServerError, nil)
	client := mock.HTTPClient()

	accessor := makeReport()

	err := provider.Check(context.Background(), client, accessor, cfg)
	assert.NoError(t, err)

	accessor.ReadFromReport(func(s *status.ClusterStatus) {
		assert.Len(t, s.Checks, 1)
		for _, c := range s.Checks {
			assert.False(t, c.Passing)
		}
	})
}

func TestChecker_CheckHandlesTimeout(t *testing.T) {
	cfg := &config.Settings{
		Prometheus: config.Prometheus{
			KubeStateMetricsServiceEndpoint: mockURL,
		},
	}
	provider := kms.NewProvider(context.Background(), cfg)

	mock := test.NewHTTPMock()
	mock.Expect(http.MethodGet, "", http.StatusOK, nil)
	client := mock.HTTPClient()

	accessor := makeReport()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := provider.Check(ctx, client, accessor, cfg)
	assert.NoError(t, err)

	accessor.ReadFromReport(func(s *status.ClusterStatus) {
		assert.Len(t, s.Checks, 1)
		for _, c := range s.Checks {
			assert.False(t, c.Passing)
		}
	})
}
