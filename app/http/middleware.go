// SPDX-FileCopyrightText: Copyright (c) 2016-2024, CloudZero, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/cloudzero/cloudzero-agent/app/instr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

// Define your Prometheus metrics
var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of all HTTP requests processed, labeled by route, method and status code.",
		},
		[]string{"route", "method", "status_code"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request durations, labeled by route, method and status code.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route", "method", "status_code"},
	)

	httpMiddlewareStatsOnce sync.Once
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func MetricsMiddlewareWrapper(next http.Handler) http.Handler {
	httpMiddlewareStatsOnce.Do(func() {
		prometheus.MustRegister(RequestCount, RequestDuration)
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := instr.StartSpan(r.Context(), "metricsHttpMiddleware")
		defer span.End()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(recorder, r)

		statusCode := strconv.Itoa(recorder.status)
		route := r.URL.Path
		method := r.Method

		// Increment the request count
		RequestCount.WithLabelValues(route, method, statusCode).Inc()

		// Observe the request duration
		RequestDuration.WithLabelValues(route, method, statusCode).Observe(span.GetDuration().Seconds())
	})
}

func LoggingMiddlewareWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(recorder, r)

		duration := time.Since(startTime)
		statusCode := recorder.status
		route := r.URL.Path
		method := r.Method

		// Log the request details
		log.Ctx(r.Context()).Debug().
			Str("method", method).
			Str("route", route).
			Int("status_code", statusCode).
			Str("status", http.StatusText(statusCode)).
			Dur("duration", duration).
			Msg("HTTP request")
	})
}
