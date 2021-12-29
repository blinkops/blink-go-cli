package utils

import (
	"github.com/blinkops/blink-go-cli/gen/client"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/spf13/viper"
)

func NewTransport() *httptransport.Runtime {
	hostname := viper.GetString(consts.HostnameEntry)
	scheme := viper.GetString(consts.SchemeEntry)
	apiKey := viper.GetString(consts.ApiKeyEntry)

	if scheme == "" {
		scheme = "https"
	}

	r := httptransport.New(hostname, client.DefaultBasePath, []string{scheme})
	r.DefaultAuthentication = httptransport.Compose(
		httptransport.APIKeyAuth(consts.ApiKeyHeader, "header", apiKey),
	)
	return r
}
