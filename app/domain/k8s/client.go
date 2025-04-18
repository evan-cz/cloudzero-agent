// SPDX-FileCopyrightText: Copyright (c) 2016-2024, CloudZero, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package k8s provides utilities for interacting with Kubernetes clusters.
// It includes functions for creating Kubernetes clients and managing resources within a cluster.
package k8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	EnvVarHostname   = "HOSTNAME"
	PodNamePartial   = "insights-controller-server"
	queriesPerSecond = 50
	maxBurst         = 100
)

// NewClient creates a new Kubernetes client using the provided kubeconfig file path.
// It returns a kubernetes.Interface which can be used to interact with the Kubernetes API.
// The function sets the QPS (Queries Per Second) and Burst rate for the client to ensure efficient communication with the cluster.
// If there is an error building the kubeconfig or creating the clientset, it returns an error.
func NewClient(kubeconfigPath string) (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}
	config.QPS = queriesPerSecond
	config.Burst = maxBurst
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to build clientset: %w", err)
	}
	return clientset, nil
}
