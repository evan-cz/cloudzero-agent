// SPDX-FileCopyrightText: Copyright (c) 2016-2024, CloudZero, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package egress contains code for checking egress access.
package webhook_server

import (
	"context"
	net "net/http"

	"github.com/sirupsen/logrus"

	"github.com/cloudzero/cloudzero-agent-validator/pkg/config"
	"github.com/cloudzero/cloudzero-agent-validator/pkg/diagnostic"
	"github.com/cloudzero/cloudzero-agent-validator/pkg/http"
	"github.com/cloudzero/cloudzero-agent-validator/pkg/logging"
	"github.com/cloudzero/cloudzero-agent-validator/pkg/status"
)

const DiagnosticWebhookServerAccess = config.DiagnosticWebhookServerAccess

type checker struct {
	cfg    *config.Settings
	logger *logrus.Entry
}

func NewProvider(ctx context.Context, cfg *config.Settings) diagnostic.Provider {
	return &checker{
		cfg: cfg,
		logger: logging.NewLogger().
			WithContext(ctx).WithField(logging.OpField, "egress"),
	}
}

func (c *checker) Check(ctx context.Context, client *net.Client, accessor status.Accessor) error {
	// ensure we can reach the webhook server over the k8s network
	url := c.cfg.Deployment.WebhookServerAddress + "/healthz"
	_, err := http.Do(ctx, client, net.MethodGet, nil, nil, url, nil)
	if err == nil {
		accessor.AddCheck(&status.StatusCheck{Name: DiagnosticWebhookServerAccess, Passing: true})
		return nil
	}
	accessor.AddCheck(&status.StatusCheck{Name: DiagnosticWebhookServerAccess, Passing: false, Error: err.Error()})
	return nil
}
