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
	return func(next BuilderService) BuilderService {
		return instrumentationMiddleware{requestCount, requestLatency, countResult, next}
	}
}

func (mw instrumentationMiddleware) BuildPlane(s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "buildplane", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.BuilderService.BuildPlane(s)
	return
}

type ServiceMiddleware func(BuilderService) BuilderService

type instrumentationMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	BuilderService
}
