package metrics_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/metrics"
)

func TestNewCounter(t *testing.T) {

	// create expected counters
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		comp  metrics.ICounter
		match bool
	}{
		{"same-Name-same-counter", args{"req.total.num"}, metrics.Counter("req.total.num"), true},
		{"diff-Name-diff-counter", args{"req.total.num"}, metrics.Counter("req.total.fail"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := metrics.Counter(tt.args.name); reflect.DeepEqual(got, tt.comp) != tt.match {
				t.Errorf("Counter() = %v, comp %v, match should be %v", got, tt.comp, tt.match)
			}
		})
	}
}

func TestNewGauge(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		comp  metrics.IGauge
		match bool
	}{
		{"same-Name-same-gauge", args{"cpu.load.average"}, metrics.Gauge("cpu.load.average"), true},
		{"diff-Name-diff-gauge", args{"cpu.load.average"}, metrics.Gauge("cpu.load.max"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := metrics.Gauge(tt.args.name); reflect.DeepEqual(got, tt.comp) != tt.match {
				t.Errorf("Gauge() = %v, comp %v, match should be %v", got, tt.comp, tt.match)
			}
		})
	}
}

func TestNewTimer(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		comp  metrics.ITimer
		match bool
	}{
		{"same-Name-same-timer", args{"req.1.timecost"}, metrics.Timer("req.1.timecost"), true},
		{"diff-Name-diff-timer", args{"req.1.timecost"}, metrics.Timer("req.2.timecost"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := metrics.Timer(tt.args.name); reflect.DeepEqual(got, tt.comp) != tt.match {
				t.Errorf("Timer() = %v, compared with %v, match should be %v", got, tt.comp, tt.match)
			}
		})
	}
}

func TestNewHistogram(t *testing.T) {
	buckets := metrics.NewDurationBounds(0*time.Millisecond,
		100*time.Millisecond, 500*time.Millisecond, 1000*time.Millisecond)

	type args struct {
		name    string
		buckets metrics.BucketBounds
	}
	tests := []struct {
		name  string
		args  args
		comp  metrics.IHistogram
		match bool
	}{
		{"same-Name-same-histogram", args{"cmd.1.timecost", buckets},
			metrics.Histogram("cmd.1.timecost", buckets), true},
		{"diff-Name-diff-histogram", args{"cmd.1.timecost", buckets},
			metrics.Histogram("cmd.2.timecost", buckets), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := metrics.Histogram(tt.args.name, tt.args.buckets); reflect.
				DeepEqual(got, tt.comp) != tt.match {

				t.Errorf("Histogram() = %v, comp %v, match should be %v", got, tt.comp, tt.match)
			}
		})
	}
}

func TestRegisterMetricsSink(t *testing.T) {
	type args struct {
		sink metrics.Sink
	}
	tests := []struct {
		name string
		args args
	}{
		{"noop", args{&metrics.NoopSink{}}},
		{"console", args{metrics.NewConsoleSink()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics.RegisterMetricsSink(tt.args.sink)
		})
	}
}

func TestIncrCounter(t *testing.T) {
	type args struct {
		key   string
		value float64
	}
	tests := []struct {
		name string
		args args
	}{
		{"counter-1", args{"req.total", 100}},
		{"counter-2", args{"req.fail", 1}},
		{"counter-3", args{"req.succ", 99}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics.IncrCounter(tt.args.key, tt.args.value)
		})
	}
}

func TestSetGauge(t *testing.T) {
	type args struct {
		key   string
		value float64
	}
	tests := []struct {
		name string
		args args
	}{
		{"gauge-1", args{"cpu.avgload", 70.1}},
		{"gauge-2", args{"mem.avgload", 80.0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics.SetGauge(tt.args.key, tt.args.value)
		})
	}
}

func TestRecordTimer(t *testing.T) {
	type args struct {
		key      string
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		{"timer-1", args{"timer.cmd.1", time.Second}},
		{"timer-2", args{"timer.cmd.2", time.Second * 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics.RecordTimer(tt.args.key, tt.args.duration)
		})
	}
}

func TestAddSample(t *testing.T) {
	metrics.Histogram("timecost.dist", metrics.NewDurationBounds(time.Second,
		time.Second*2, time.Second*3, time.Second*4))
	metrics.RegisterMetricsSink(metrics.NewConsoleSink())
	type args struct {
		key     string
		buckets metrics.BucketBounds
		value   float64
	}
	buckets := metrics.NewDurationBounds(time.Second, time.Second*2, time.Second*5)
	tests := []struct {
		name string
		args args
	}{
		{"histogram-1", args{"timecost.dist", buckets, float64(time.Second)}},
		{"histogram-2", args{"timecost.dist", buckets, float64(time.Second * 2)}},
		{"histogram-2", args{"timecost.dist", buckets, float64(time.Second * 3)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics.AddSample(tt.args.key, tt.args.buckets, tt.args.value)
		})
	}
}

type unhealthySink struct{}

func (u *unhealthySink) Name() string {
	return "unhealthy"
}

func (u *unhealthySink) Report(rec metrics.Record, opts ...metrics.Option) error {
	time.Sleep(time.Millisecond * 100)
	return errors.New("timeout")
}

type unstableSink struct{}

func (p *unstableSink) Name() string {
	return "unstable"
}

func (p *unstableSink) Report(rec metrics.Record, opts ...metrics.Option) error {
	time.Sleep(time.Millisecond * 100)
	return errors.New("backend error")
}

func TestReport(t *testing.T) {

	metrics.RegisterMetricsSink(metrics.NewConsoleSink())
	metrics.RegisterMetricsSink(&metrics.NoopSink{})
	metrics.RegisterMetricsSink(&unhealthySink{})
	metrics.RegisterMetricsSink(&unstableSink{})

	rec := metrics.NewSingleDimensionMetrics("total.req", float64(100), metrics.PolicySUM)
	tests := []struct {
		name    string
		rec     metrics.Record
		wantErr bool
	}{
		{"reportHasError", rec, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := metrics.Report(rec); (err != nil) != tt.wantErr {
				t.Errorf("Report() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
