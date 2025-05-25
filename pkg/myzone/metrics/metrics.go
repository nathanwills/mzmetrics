package metrics

import (
	"github.com/nathanwills/mzmetrics/pkg/myzone"
	"github.com/prometheus/client_golang/prometheus"
)

// Define the structure for the metrics
type metrics struct {
	// Aircon information metrics
	acState          *prometheus.GaugeVec
	acTemperatureSet *prometheus.GaugeVec

	// Zone information metrics
	zoneTemperatureMeasured *prometheus.GaugeVec
	zoneTemperatureSet      *prometheus.GaugeVec
	zoneState               *prometheus.GaugeVec
	zoneRssi                *prometheus.GaugeVec
	zoneDamper              *prometheus.GaugeVec

	// Other general metrics for system state
	systemErrorCode prometheus.Gauge
}

// New function to instantiate the metrics struct and register the metrics
func New(reg prometheus.Registerer) *metrics {
	// Instantiate the metrics struct
	m := &metrics{
		// Aircon information metrics
		acState: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "aircon_state",
				Help: "State of the aircon (on=1, off=0)",
			},
			[]string{"aircon_name"},
		),
		acTemperatureSet: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "aircon_set_temperature",
				Help: "Set temperature for the aircon",
			},
			[]string{"aircon_name"},
		),

		// Zone information metrics
		zoneTemperatureMeasured: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "zone_measured_temperature",
				Help: "Measured temperature in the zone",
			},
			[]string{"zone_name", "aircon_name"},
		),
		zoneTemperatureSet: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "zone_set_temperature",
				Help: "Set temperature in the zone",
			},
			[]string{"zone_name", "aircon_name"},
		),
		zoneState: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "zone_state",
				Help: "State of the zone (open=1, closed=0)",
			},
			[]string{"zone_name", "aircon_name"},
		),
		zoneRssi: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "zone_rssi",
				Help: "RSSI signal strength for the zone",
			},
			[]string{"zone_name", "aircon_name"},
		),
		zoneDamper: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "zone_damper",
				Help: "Zone damper open percentage",
			},
			[]string{"zone_name", "aircon_name"},
		),

		// Other general metrics for system state
		systemErrorCode: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "system_error_code",
				Help: "System error code (0 for no error)",
			},
		),
	}

	// Register the metrics with Prometheus
	reg.MustRegister(m.acState)
	reg.MustRegister(m.acTemperatureSet)
	reg.MustRegister(m.zoneTemperatureMeasured)
	reg.MustRegister(m.zoneTemperatureSet)
	reg.MustRegister(m.zoneState)
	reg.MustRegister(m.zoneRssi)
	reg.MustRegister(m.systemErrorCode)

	return m
}

func (m *metrics) SetMetrics(ac *myzone.Aircon) {
	// Iterate over the aircon data and update metrics for each aircon
	for _, acData := range ac.Aircons {
		// Set the aircon state and temperature
		if acData.Info.State == "on" {
			m.acState.WithLabelValues(acData.Info.Name).Set(1)
		} else {
			m.acState.WithLabelValues(acData.Info.Name).Set(0)
		}
		m.acTemperatureSet.WithLabelValues(acData.Info.Name).Set(acData.Info.SetTemp)

		// Iterate over zones and set metrics for each zone
		for _, zone := range acData.Zones {
			m.zoneTemperatureMeasured.WithLabelValues(zone.Name, acData.Info.Name).Set(zone.MeasuredTemp)
			m.zoneTemperatureSet.WithLabelValues(zone.Name, acData.Info.Name).Set(zone.SetTemp)

			// Set the zone state (open=1, closed=0)
			if zone.State == "open" {
				m.zoneState.WithLabelValues(zone.Name, acData.Info.Name).Set(1)
			} else {
				m.zoneState.WithLabelValues(zone.Name, acData.Info.Name).Set(0)
			}

			// Set the RSSI for the zone
			m.zoneRssi.WithLabelValues(zone.Name, acData.Info.Name).Set(float64(zone.Rssi))
		}

		// Set a system error code, assuming no error in this case
		m.systemErrorCode.Set(0) // No error
	}
}
