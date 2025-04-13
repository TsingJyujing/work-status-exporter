package workstatus

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"net/http"
	"time"
	"work-status-exporter/logging"
)

var (
	MetricNamespace        = "work_status"
	zoomMeetingTimeSeconds = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: MetricNamespace,
		Subsystem: "zoom",
		Name:      "meeting_time_seconds",
		Help:      "Total time spent in zoom meetings in seconds",
	})
	cameraStreamingTimeSeconds = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: MetricNamespace,
		Subsystem: "camera",
		Name:      "streaming_time_seconds",
		Help:      "Total time camera is on in seconds",
	})
	cameraActivatedTimeSeconds = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: MetricNamespace,
		Subsystem: "camera",
		Name:      "activated_time_seconds",
		Help:      "Total time camera is on in seconds",
	})
)

func RegisterPrometheusMetricsUpdatingJob(monitoringInterval time.Duration) {
	c := cron.New(cron.WithSeconds())
	spec := fmt.Sprintf("@every %v", monitoringInterval) // every monitoringInterval seconds
	logging.Logger.WithField("spec", spec).Debug("Adding zoom monitoring job")
	_, zoomJobErr := c.AddFunc(spec, func() {
		zoomStatus, zoomErr := GetZoomMeetingStatus()
		if zoomErr != nil {
			logging.Logger.WithError(zoomErr).Error("error getting zoom meeting status")
		} else {
			logging.Logger.WithField("zoomStatus", zoomStatus).Debug("Got zoom meeting status")
		}
		if zoomStatus {
			zoomMeetingTimeSeconds.Add(monitoringInterval.Seconds())
		}
	})
	if zoomJobErr != nil {
		logging.Logger.WithError(zoomJobErr).Fatal("Error while adding zoom monitoring job")
	}
	_, cameraJobErr := c.AddFunc(spec, func() {
		cameraActivated, cameraIsStreaming, cameraStatusErr := GetMacOSCameraStatus()
		if cameraStatusErr != nil {
			logging.Logger.WithError(cameraStatusErr).Error("error getting camera status")
		}
		if cameraActivated {
			cameraActivatedTimeSeconds.Add(monitoringInterval.Seconds())
		}
		if cameraIsStreaming {
			cameraStreamingTimeSeconds.Add(monitoringInterval.Seconds())
		}
	})
	if cameraJobErr != nil {
		logging.Logger.WithError(cameraJobErr).Fatal("Error while adding camera job")
	}
	go c.Run()
}

func StartPrometheusMetricsServer(httpAddr string, monitoringInterval time.Duration) {
	RegisterPrometheusMetricsUpdatingJob(monitoringInterval)
	logging.Logger.Infof("Starting prometheus metrics server on %s", httpAddr)
	httpSrv := &http.Server{Addr: httpAddr}
	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.Handler())
	httpSrv.Handler = m
	err := httpSrv.ListenAndServe()
	if err != nil {
		logging.Logger.WithError(err).Fatal("Failed to start metrics server")
	}
}
