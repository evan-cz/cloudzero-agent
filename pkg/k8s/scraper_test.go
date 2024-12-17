// SPDX-FileCopyrightText: Copyright (c) 2016-2024, CloudZero, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/util/homedir"

	"github.com/cloudzero/cloudzero-insights-controller/pkg/config"
	"github.com/cloudzero/cloudzero-insights-controller/pkg/storage"
	"github.com/cloudzero/cloudzero-insights-controller/pkg/types"
)

type MockClusterScraper struct {
	k8sClient kubernetes.Interface
}

func (m *MockClusterScraper) Start() {
	fmt.Printf("got /hello request\n")
	// List all namespaces
	namespaces, err := m.k8sClient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error().Msgf("Error listing namespaces: %v", err)
	}
	log.Info().Msgf("There are %d namespaces in the cluster\n", len(namespaces.Items))
}

func TestScraper_Start(t *testing.T) {
	db := storage.SetupDatabase()
	settings := &config.Settings{}
	writer := storage.NewWriter(db, settings)

	t.Run("with fake client", func(t *testing.T) {
		clientset := fake.NewSimpleClientset()
		scraper := NewScraper(clientset, writer, settings)
		scraper.Start()
		var results []types.ResourceTags
		db.Find(&results)

	})

	t.Run("with real client; integration test", func(t *testing.T) {
		if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
			t.Skip("Skipping integration test as RUN_INTEGRATION_TESTS is not set to true")
		}
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		settings.K8sClient.PaginationLimit = 3

		k8sClient, err := BuildKubeClient(kubeconfig)

		assert.NoError(t, err)
		scraper := NewScraper(k8sClient, writer, settings)
		scraper.Start()
		time.Sleep(5 * time.Second)
		var results []types.ResourceTags
		db.Find(&results)
		assert.NotEmpty(t, results)
	})

}
