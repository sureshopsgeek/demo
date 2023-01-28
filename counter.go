package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var mutex = &sync.Mutex{}

var metric = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_count",
		Help: "response count partitioned by status code",
	},
	[]string{"response"},
)

var (
	counterPort int
	promPort    int
	verbose     bool
)

type Count struct {
	Code int
}

func CountHandler(w http.ResponseWriter, r *http.Request) {
	if verbose {
		reqDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("REQUEST:\n%s", string(reqDump))
	}
	if r.Method == http.MethodPut {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "error reading request body\n")
			return
		}
		var ct Count
		err = json.Unmarshal(body, &ct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "invalid request body\n")
			return
		}
		mutex.Lock()
		metric.With(prometheus.Labels{"response": fmt.Sprintf("%d", ct.Code)}).Inc()
		mutex.Unlock()
		w.WriteHeader(http.StatusOK)
	} else {
		io.WriteString(w, "method not allowed\n")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func startServer(port int, mux *http.ServeMux) {
	log.Printf("Starting HTTP server at port: %d\n", port)
	server := http.Server{Addr: fmt.Sprintf(":%d", port),
		Handler: mux}
	log.Println(server.ListenAndServe())
}

func main() {
	flag.IntVar(&counterPort, "cp", 8080, "counter server port")
	flag.IntVar(&promPort, "mp", 9201, "metrics server port")
	flag.BoolVar(&verbose, "v", false, "set to enable verbose logging")
	flag.Parse()
	reg := prometheus.NewRegistry()
	reg.MustRegister(metric)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	counterMux := http.NewServeMux()
	counterMux.HandleFunc("/counter", CountHandler)
	metricMux := http.NewServeMux()
	metricMux.Handle("/metrics",promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	go func() {
		startServer(counterPort, counterMux)
		wg.Done()
	}()
	go func() {
		startServer(promPort, metricMux)
		wg.Done()
	}()
	wg.Wait()
}
