package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	"github.com/snowdiceX/metrics_forwarder/collector"
	"github.com/snowdiceX/metrics_forwarder/log"
)

func main() {
	// Since we are dealing with custom Collector implementations, it might
	// be a good idea to try it out with a pedantic registry.
	reg := prometheus.NewPedanticRegistry()

	// // Construct forwarder collector. In real code, we would assign them to
	// // variables to then do something with them.
	collector := collector.NewForwarderCollector(
		"mainnet", "78", "http://127.0.0.1:26660/metrics", reg)

	// // // Add the standard process and Go metrics to the custom registry.
	// // reg.MustRegister(
	// // 	prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
	// // 	prometheus.NewGoCollector(),
	// // )

	// http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	// log.Fatal(http.ListenAndServe(":8080", nil))

	// families, err := reg.Gather()
	// if err != nil {
	// 	fmt.Println("gather error: ", err)
	// 	return
	// }
	// for _, f := range families {
	// 	fmt.Println("families: ", f.GetName(), "; ", f.GetHelp(), "; ", f.GetType())
	// }
	// fmt.Println("done")

	if err := push.New("http://127.0.0.1:9091", "irishub").
		Collector(collector).
		Grouping("service", "blockchain").
		Push(); err != nil {
		log.Errorf("Could not push metrics to pushgateway: %d", err)
	}
}
