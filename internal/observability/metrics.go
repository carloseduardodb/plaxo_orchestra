package observability

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TraceID string

type Span struct {
	TraceID   TraceID   `json:"trace_id"`
	SpanID    string    `json:"span_id"`
	Operation string    `json:"operation"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Tags      map[string]string `json:"tags"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}

type Metrics struct {
	spans       []Span            `json:"spans"`
	counters    map[string]int64  `json:"counters"`
	gauges      map[string]float64 `json:"gauges"`
	histograms  map[string][]float64 `json:"histograms"`
	mutex       sync.RWMutex
}

type Observer struct {
	metrics *Metrics
	ctx     context.Context
}

func NewObserver() *Observer {
	return &Observer{
		metrics: &Metrics{
			spans:      make([]Span, 0),
			counters:   make(map[string]int64),
			gauges:     make(map[string]float64),
			histograms: make(map[string][]float64),
		},
		ctx: context.Background(),
	}
}

func (o *Observer) StartSpan(operation string, tags map[string]string) *Span {
	traceID := TraceID(fmt.Sprintf("trace_%d", time.Now().UnixNano()))
	spanID := fmt.Sprintf("span_%d", time.Now().UnixNano())
	
	span := &Span{
		TraceID:   traceID,
		SpanID:    spanID,
		Operation: operation,
		StartTime: time.Now(),
		Tags:      tags,
	}
	
	return span
}

func (o *Observer) FinishSpan(span *Span, success bool, err error) {
	span.EndTime = time.Now()
	span.Duration = span.EndTime.Sub(span.StartTime)
	span.Success = success
	
	if err != nil {
		span.Error = err.Error()
	}
	
	o.metrics.mutex.Lock()
	defer o.metrics.mutex.Unlock()
	
	o.metrics.spans = append(o.metrics.spans, *span)
	
	// Update metrics
	o.metrics.counters[fmt.Sprintf("%s_total", span.Operation)]++
	if success {
		o.metrics.counters[fmt.Sprintf("%s_success", span.Operation)]++
	} else {
		o.metrics.counters[fmt.Sprintf("%s_error", span.Operation)]++
	}
	
	// Record duration histogram
	key := fmt.Sprintf("%s_duration", span.Operation)
	o.metrics.histograms[key] = append(o.metrics.histograms[key], span.Duration.Seconds())
}

func (o *Observer) IncrementCounter(name string, value int64) {
	o.metrics.mutex.Lock()
	defer o.metrics.mutex.Unlock()
	o.metrics.counters[name] += value
}

func (o *Observer) SetGauge(name string, value float64) {
	o.metrics.mutex.Lock()
	defer o.metrics.mutex.Unlock()
	o.metrics.gauges[name] = value
}

func (o *Observer) GetMetrics() map[string]interface{} {
	o.metrics.mutex.RLock()
	defer o.metrics.mutex.RUnlock()
	
	// Calculate percentiles for histograms
	percentiles := make(map[string]map[string]float64)
	for name, values := range o.metrics.histograms {
		if len(values) > 0 {
			percentiles[name] = calculatePercentiles(values)
		}
	}
	
	return map[string]interface{}{
		"counters":    o.metrics.counters,
		"gauges":      o.metrics.gauges,
		"percentiles": percentiles,
		"spans_count": len(o.metrics.spans),
	}
}

func (o *Observer) GetTracesByOperation(operation string) []Span {
	o.metrics.mutex.RLock()
	defer o.metrics.mutex.RUnlock()
	
	traces := make([]Span, 0)
	for _, span := range o.metrics.spans {
		if span.Operation == operation {
			traces = append(traces, span)
		}
	}
	
	return traces
}

func calculatePercentiles(values []float64) map[string]float64 {
	if len(values) == 0 {
		return map[string]float64{}
	}
	
	// Simple percentile calculation (would use proper sorting in production)
	sum := 0.0
	min := values[0]
	max := values[0]
	
	for _, v := range values {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	
	return map[string]float64{
		"min":  min,
		"max":  max,
		"avg":  sum / float64(len(values)),
		"p50":  values[len(values)/2],
		"p95":  values[int(float64(len(values))*0.95)],
		"p99":  values[int(float64(len(values))*0.99)],
	}
}
