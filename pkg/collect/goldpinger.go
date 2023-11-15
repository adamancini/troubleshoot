package collect

import (
	"context"
	"time"

	troubleshootv1beta2 "github.com/replicatedhq/troubleshoot/pkg/apis/troubleshoot/v1beta2"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	DefaultGoldpingerNamespace = "kurl"
)

type GoldpingerHTTPRequest struct {
	Route          string
	Format         string
	DefaultTimeout string
}

var GoldpingerHTTPRequests = []GoldpingerHTTPRequest{
	// {
	// 	ID:             "cluster_health",
	// 	Path:           "http://localhost:8080/healthz",
	// 	Format:         "json",
	// 	DefaultTimeout: "30s",
	// },
	{
		Route:          "check_all",
		Format:         "json",
		DefaultTimeout: "30s",
	},
}

type GoldpingerHost struct {
	HostIP  string `json:"hostIP"`
	PodIP   string `json:"podIP"`
	PodName string `json:"podName"`
}

type PingResult struct {
	HostIP         string    `json:"HostIP"`
	OK             bool      `json:"OK"`
	PingTime       time.Time `json:"PingTime"`
	PodIP          string    `json:"PodIP"`
	BootTime       time.Time `json:"response.boot_time"`
	ResponseTimeMs int       `json:"response-time-ms,omitempty"`
	StatusCode     int       `json:"status-code"`
}

type GoldpingerResult struct {
	Hosts     []GoldpingerHost      `json:"hosts"`
	Responses map[string]PingResult `json:"responses"`
}

type CollectGoldpinger struct {
	Collector    *troubleshootv1beta2.Goldpinger
	BundlePath   string
	Namespace    string
	ClientConfig *rest.Config
	Client       kubernetes.Interface
	Context      context.Context
	RBACErrors
}
