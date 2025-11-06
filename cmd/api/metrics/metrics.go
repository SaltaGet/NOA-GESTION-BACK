package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    RequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    RequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )

    ActiveProducts = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_products_total",
            Help: "Total number of active products",
        },
    )

    StockLevels = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "product_stock_levels",
            Help: "Current stock levels by product",
        },
        []string{"product_id", "product_name", "location"},
    )
)