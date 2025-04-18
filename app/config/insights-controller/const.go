// SPDX-FileCopyrightText: Copyright (c) 2016-2024, CloudZero, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

type ResourceType int

const (
	Unknown     ResourceType = 0
	Deployment  ResourceType = 1
	StatefulSet ResourceType = 2
	Pod         ResourceType = 3
	Node        ResourceType = 4
	Namespace   ResourceType = 5
	Job         ResourceType = 6
	CronJob     ResourceType = 7
	DaemonSet   ResourceType = 8
)

var ResourceTypeToMetricName = map[ResourceType]string{
	Unknown:     "unknown",
	Deployment:  "deployment",
	StatefulSet: "statefulset",
	Pod:         "pod",
	Node:        "node",
	Namespace:   "namespace",
	Job:         "job",
	CronJob:     "cronjob",
	DaemonSet:   "daemonset",
}
