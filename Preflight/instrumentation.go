package main

import (
	"fmt"
	"github.com/go-kit/kit/metrics"
	"time"
)

func instrumentingMiddleware(
	requestCount metrics.Counter,
	requestLatency metrics.Histogram,
	countResult metrics.Histogram,
) ServiceMiddleware {
	return func(next PreflightService) PreflightService {
		return instrumentationMiddleware{requestCount, requestLatency, countResult, next}
	}
}

func (mw instrumentationMiddleware) ProvisionPlane(s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "provisionplane", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.PreflightService.ProvisionPlane(s)
	return
}

type ServiceMiddleware func(PreflightService) PreflightService

type instrumentationMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	PreflightService
}