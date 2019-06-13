package collector

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/snowdiceX/metrics-forwarder/log"
)

// MetricsSource is a source of a prometheus metrics.
// It models a source of metrics running in an
// application. Thus, we implement a custom Collector called
// ForwarderCollector, which collects information from a MetricsSource
// parsing its provided metrics and push them into Prometheus Pushgateway.
type MetricsSource struct {
	Zone string
	Host string
	URL  string
}

// Pull is a method for the metrics gathering.
// Since it may actually be really
// expensive, it must only be called once per collection.
func (c *MetricsSource) Pull() (
	metricFamilies map[string]*dto.MetricFamily, err error) {

	resp, err := http.Get(c.URL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	text := string(body)
	log.Debugf("pull response: %s", text)

	// referenced github.com/prometheus/pushgateway/handler/push.go
	var parser expfmt.TextParser
	metricFamilies, err = parser.TextToMetricFamilies(
		strings.NewReader(text))
	if err != nil {
		return nil, err
	}
	log.Infof("metrics: %d", len(metricFamilies))
	for k, m := range metricFamilies {
		log.Trace("metric family: ", k, m.GetName(), "; ", m.Help, "; ", m.Type)
	}
	return
}

// ForwarderCollector implements the Collector interface.
type ForwarderCollector struct {
	Source *MetricsSource
}

// Describe is implemented with DescribeByCollect. That's possible because the
// Collect method will always return the same two metrics with the same two
// descriptors.
func (cc ForwarderCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

// Collect first triggers the ReallyExpensiveAssessmentOfTheSystemState. Then it
// creates constant metrics for each host on the fly based on the returned data.
//
// Note that Collect could be called concurrently, so we depend on
// ReallyExpensiveAssessmentOfTheSystemState to be concurrency-safe.
func (cc ForwarderCollector) Collect(ch chan<- prometheus.Metric) {
	families, err := cc.Source.Pull()
	if err != nil {
		log.Errorf("pull metrics error: %v", err)
		return
	}
	for k, f := range families {
		log.Trace("collect metric family: ", k)
		switch f.GetType() {
		case dto.MetricType_COUNTER:
			{
				buildCounter(cc.Source.Host, f, ch)
			}
		case dto.MetricType_GAUGE:
			{
				buildGauge(cc.Source.Host, f, ch)
			}
		}
	}
}

// NewForwarderCollector first creates a MetricsSource
// instance. Then, it creates a ForwarderCollector for the just created
// MetricsSource. Finally, it registers the ForwarderCollector with a
// wrapping Registerer that adds the zone as a label. In this way, the metrics
// collected by different ForwarderCollectors do not collide.
func NewForwarderCollector(zone, host, url string,
	reg prometheus.Registerer) *ForwarderCollector {
	c := &MetricsSource{
		Zone: zone,
		Host: host,
		URL:  url,
	}
	cc := ForwarderCollector{Source: c}
	prometheus.WrapRegistererWith(
		prometheus.Labels{"zone": zone}, reg).MustRegister(cc)
	return &cc
}

func buildCounter(host string, family *dto.MetricFamily,
	ch chan<- prometheus.Metric) {
	metrics := family.GetMetric()
	for _, m := range metrics {
		g := m.GetCounter()
		lnames := []string{"host"}
		lvalues := []string{host}
		labels := m.GetLabel()
		for _, l := range labels {
			lnames = append(lnames, l.GetName())
			lvalues = append(lvalues, l.GetValue())
			log.Trace("label pair: %s; %s", l.GetName(), l.GetValue())
		}
		desc := prometheus.NewDesc(
			family.GetName(),
			family.GetHelp(),
			lnames, nil,
		)
		ch <- prometheus.MustNewConstMetric(
			desc,
			prometheus.CounterValue,
			g.GetValue(),
			lvalues...,
		)
	}
}

func buildGauge(host string, family *dto.MetricFamily,
	ch chan<- prometheus.Metric) {
	metrics := family.GetMetric()
	for _, m := range metrics {
		g := m.GetGauge()
		lnames := []string{"host"}
		lvalues := []string{host}
		labels := m.GetLabel()
		for _, l := range labels {
			lnames = append(lnames, l.GetName())
			lvalues = append(lvalues, l.GetValue())
			log.Trace("label pair: %s; %s", l.GetName(), l.GetValue())
		}
		desc := prometheus.NewDesc(
			family.GetName(),
			family.GetHelp(),
			lnames, nil,
		)
		ch <- prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			g.GetValue(),
			lvalues...,
		)
	}
}
