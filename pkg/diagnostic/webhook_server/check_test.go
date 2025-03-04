// SPDX-FileCopyrightText: Copyright (c) 2016-2024, CloudZero, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package webhook_server_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cloudzero/cloudzero-agent-validator/pkg/config"
	"github.com/cloudzero/cloudzero-agent-validator/pkg/diagnostic/webhook_server"
	"github.com/cloudzero/cloudzero-agent-validator/pkg/status"
	"github.com/cloudzero/cloudzero-agent-validator/test"
)

const (
	mockURL = "https://mock.webhook-server.svc.cluster.local"
)

func makeReport() status.Accessor {
	return status.NewAccessor(&status.ClusterStatus{})
}

func TestChecker_CheckOK(t *testing.T) {
	cfg := &config.Settings{
		Deployment: config.Deployment{
			WebhookServerAddress: mockURL,
		},
	}

	provider := webhook_server.NewProvider(context.Background(), cfg)

	mock := test.NewHTTPMock()
	mock.Expect(http.MethodGet, "Hello World", http.StatusOK, nil)
	client := mock.HTTPClient()

	accessor := makeReport()

	err := provider.Check(context.Background(), client, accessor)
	assert.NoError(t, err)

	accessor.ReadFromReport(func(s *status.ClusterStatus) {
		assert.Len(t, s.Checks, 1)
		assert.True(t, s.Checks[0].Passing)
		assert.Empty(t, s.Checks[0].Error)
	})
}

func TestChecker_CheckBadKey(t *testing.T) {
	cfg := &config.Settings{
		Cloudzero: config.Cloudzero{
			Host:       mockURL,
			Credential: "your-api-key",
		},
	}

	provider := webhook_server.NewProvider(context.Background(), cfg)

	mock := test.NewHTTPMock()
	mock.Expect(http.MethodGet, "", http.StatusUnauthorized, nil)
	client := mock.HTTPClient()

	accessor := makeReport()
	err := provider.Check(context.Background(), client, accessor)
	assert.NoError(t, err)

	accessor.ReadFromReport(func(s *status.ClusterStatus) {
		assert.Len(t, s.Checks, 1)
		assert.False(t, s.Checks[0].Passing)
		assert.NotEmpty(t, s.Checks[0].Error)
	})
}

func TestChecker_CheckErrorCondition(t *testing.T) {
	cfg := &config.Settings{
		Cloudzero: config.Cloudzero{
			Host:       mockURL,
			Credential: "your-api-key",
		},
	}

	provider := webhook_server.NewProvider(context.Background(), cfg)

	mock := test.NewHTTPMock()
	mock.Expect(http.MethodGet, "", http.StatusBadGateway, nil)
	client := mock.HTTPClient()

	accessor := makeReport()
	err := provider.Check(context.Background(), client, accessor)
	assert.NoError(t, err)
	accessor.ReadFromReport(func(s *status.ClusterStatus) {
		assert.Len(t, s.Checks, 1)
		assert.Equal(t, webhook_server.DiagnosticWebhookServerAccess, s.Checks[0].Name)
		assert.False(t, s.Checks[0].Passing)
		assert.NotEmpty(t, s.Checks[0].Error)
	})
}
