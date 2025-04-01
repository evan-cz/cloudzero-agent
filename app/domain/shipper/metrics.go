// SPDX-FileCopyrightText: Copyright (c) 2016-2024, CloudZero, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package shipper

import (
	"github.com/cloudzero/cloudzero-agent-validator/app/instr"
	"github.com/prometheus/client_golang/prometheus"
)

var (

	// Other
	// ----------------------------------------------------------
	metricShutdownTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shipper_shutdown_total",
			Help: "Total number of shutdown requests",
		},
		[]string{},
	)

	// New File Processing
	// ----------------------------------------------------------
	metricNewFilesErrorTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shipper_new_files_error_total",
			Help: "Total number of errors encountered when running segments of the program",
		},
		[]string{},
	)

	metricNewFilesProcessingCurrent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "shipper_new_files_processing_current",
			Help: "The current number of files being processed by the shipper",
		},
		[]string{},
	)

	// Generic Request Handling
	// ----------------------------------------------------------
	metricHandleRequestFileCount = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "shipper_handle_request_file_count",
			Help:    "Number of files requested for a replay request",
			Buckets: prometheus.ExponentialBuckets(10, 2, 15),
		},
	)

	metricHandleRequestSuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shipper_handle_request_success_total",
			Help: "Total number of successful runs of the `HandleRequest` function",
		},
		[]string{},
	)

	// Presigned URLs
	// ----------------------------------------------------------
	metricPresignedURLErrorTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shipper_presigned_url_error_total",
			Help: "Total number of errors seen when creating all presigned urls",
		},
		[]string{},
	)

	// Replay Requests
	// ----------------------------------------------------------
	metricReplayRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shipper_replay_request_total",
			Help: "Total number of replay requests receieved from the remote file receiver.",
		},
		[]string{},
	)

	metricReplayRequestCurrent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "shipper_replay_request_current",
			Help: "The current number of replay requests queued",
		},
		[]string{},
	)

	metricReplayRequestFileCount = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "shipper_replay_request_file_count",
			Help:    "Number of files requested for a replay request",
			Buckets: prometheus.ExponentialBuckets(10, 2, 15),
		},
	)

	metricReplayRequestErrorTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shipper_replay_request_error_total",
			Help: "Number of errors observed while processing replay requests",
		},
		[]string{},
	)

	metricReplayRequestAbandonFilesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shipper_replay_request_abandon_files_total",
			Help: "total number of abandoned files",
		},
		[]string{},
	)

	metricReplayRequestAbandonFilesErrorTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shipper_replay_request_abandon_files_error_total",
			Help: "total number of abandoned files",
		},
		[]string{},
	)

	// Disk Usage
	// ----------------------------------------------------------

	metricDiskTotalSizeBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "shipper_disk_total_size_bytes",
			Help: "Total Size (bytes) of the pv",
		},
		[]string{},
	)

	metricCurrentDiskUsageBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "shipper_disk_current_usage_bytes",
			Help: "Size (bytes) currently used in the pv",
		},
		[]string{},
	)

	metricCurrentDiskUsagePercentage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "shipper_disk_current_usage_percentage",
			Help: "Percentage currently used in the pv",
		},
		[]string{},
	)

	metricCurrentDiskUnsentFile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "shipper_disk_current_unsent_file",
			Help: "Current number of unsent files found in the pv",
		},
		[]string{},
	)

	metricCurrentDiskSentFile = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "shipper_disk_current_sent_file",
			Help: "Current number of sent files found in the pv",
		},
		[]string{},
	)

	metricCurrentDiskReplayRequest = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "shipper_disk_replay_request_current",
			Help: "Current number of replay requests found in the pv",
		},
		[]string{},
	)

	metricDiskCleanupFailureTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shipper_disk_cleanup_failure_total",
			Help: "Number of failures when purging files for disk space",
		},
		[]string{"storage_warning", "error"},
	)

	metricDiskCleanupSuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shipper_disk_cleanup_success_total",
			Help: "Number of successes when purging files for disk space",
		},
		[]string{"storage_warning"},
	)

	metricDiskCleanupPercentage = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "shipper_disk_cleanup_percentage",
			Help:    "Percent removed from the storage volume during a purge operation",
			Buckets: prometheus.LinearBuckets(0, 10, 11), // 0% to 100% in steps of 10%
		},
	)
)

func InitMetrics() (*instr.PrometheusMetrics, error) {
	return instr.NewPrometheusMetrics(
		// instr.WithDefaultRegistry(),
		instr.WithPromMetrics(
			// other
			metricShutdownTotal,

			// new files
			metricNewFilesErrorTotal,
			metricNewFilesProcessingCurrent,

			// generic request handling
			metricHandleRequestFileCount,
			metricHandleRequestSuccessTotal,

			metricPresignedURLErrorTotal,

			// replay requests
			metricReplayRequestTotal,
			metricReplayRequestCurrent,
			metricReplayRequestFileCount,
			metricReplayRequestErrorTotal,
			metricReplayRequestAbandonFilesTotal,
			metricReplayRequestAbandonFilesErrorTotal,

			// disk usage
			metricDiskTotalSizeBytes,
			metricCurrentDiskUsageBytes,
			metricCurrentDiskUsagePercentage,
			metricCurrentDiskUnsentFile,
			metricCurrentDiskSentFile,
			metricCurrentDiskReplayRequest,
			metricDiskCleanupFailureTotal,
			metricDiskCleanupSuccessTotal,
			metricDiskCleanupPercentage,
		),
	)
}
