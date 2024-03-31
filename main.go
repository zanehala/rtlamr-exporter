package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type message struct {
	Time    time.Time              `json:"Time"`
	Type    string                 `json:"Type"`
	Message map[string]interface{} `json:"Message"`
}

type SCMPlus struct {
	ProtocolID   uint8  `json:"ProtocolID"`
	EndpointType uint8  `json:"EndpointType"`
	EndpointID   uint32 `json:"EndpointID"`
	Consumption  uint32 `json:"Consumption"`
	Tamper       uint16 `json:"Tamper"`
}

var metrics = promauto.NewGaugeVec(
	prometheus.GaugeOpts{Name: "meters", Help: "Metrics about various consumption meters"},
	[]string{"protocolType", "meterId", "meterType"},
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":9090", nil)
	}()

	reader := bufio.NewScanner(os.Stdin)
	var input message

	for reader.Scan() {
		err := json.Unmarshal([]byte(reader.Text()), &input)
		if err != nil {
			logger.Warn("Failed to parse input.", "message", reader.Text())
			continue
		}

		marshalled, err := json.Marshal(input.Message) // There is probably a better way to do this
		if err != nil {
			logger.Warn("Failed to parse consumption message.", "message", fmt.Sprintf("%s", input.Message))
			continue
		}

		// TODO: add more message types
		switch input.Type {
		case "SCM+":
			var scmplus SCMPlus
			json.Unmarshal(marshalled, &scmplus)
			metrics.WithLabelValues(
				input.Type,
				strconv.FormatUint(uint64(scmplus.EndpointID), 10),
				strconv.FormatUint(uint64(scmplus.EndpointType), 10),
			).Set(float64(scmplus.Consumption))
		}
	}
}
