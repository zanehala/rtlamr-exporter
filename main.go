package main

import (
	"bufio"
	"encoding/json"
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
	Time    time.Time        `json:"Time"`
	Type    string           `json:"Type"`
	Message *json.RawMessage `json:"Message"`
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

		switch input.Type {
		case "SCM+":
			var scmplus SCMPlus
			json.Unmarshal(*input.Message, &scmplus)
			metrics.WithLabelValues(
				input.Type,
				strconv.FormatUint(uint64(scmplus.EndpointID), 10),
				strconv.FormatUint(uint64(scmplus.EndpointType), 10),
			).Set(float64(scmplus.Consumption))
		case "SCM":
			var scm SCM
			json.Unmarshal(*input.Message, &scm)
			metrics.WithLabelValues(
				input.Type,
				strconv.FormatUint(uint64(scm.ID), 10),
				strconv.FormatUint(uint64(scm.Type), 10),
			).Set(float64(scm.Consumption))
		case "IDM":
			var idm IDM
			json.Unmarshal(*input.Message, &idm)
			metrics.WithLabelValues(
				input.Type,
				strconv.FormatUint(uint64(idm.ERTSerialNumber), 10),
				strconv.FormatUint(uint64(idm.ERTType), 10),
			).Set(float64(idm.LastConsumptionCount))
		case "NetIDM":
			var netidm IDM
			json.Unmarshal(*input.Message, &netidm)
			metrics.WithLabelValues(
				input.Type,
				strconv.FormatUint(uint64(netidm.ERTSerialNumber), 10),
				strconv.FormatUint(uint64(netidm.ERTType), 10),
			).Set(float64(netidm.LastConsumptionCount))
		case "R900":
			var r900 R900
			json.Unmarshal(*input.Message, &r900)
			metrics.WithLabelValues(
				input.Type,
				strconv.FormatUint(uint64(r900.ID), 10),
				strconv.FormatUint(uint64(r900.Unkn1), 10),
			).Set(float64(r900.Consumption))
		case "R900BCD":
			var r900bcd R900
			json.Unmarshal(*input.Message, &r900bcd)
			metrics.WithLabelValues(
				input.Type,
				strconv.FormatUint(uint64(r900bcd.ID), 10),
				strconv.FormatUint(uint64(r900bcd.Unkn1), 10),
			).Set(float64(r900bcd.Consumption))
		default:
			logger.Info("Unrecognized message type", "message", reader.Text())
		}
	}
}
