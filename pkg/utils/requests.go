package utils

import (
	"fmt"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

func NewRequest(method string, endpoint string, body io.Reader, headers map[string]string) (*http.Request, error) {

	request, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	for header, value := range headers {
		request.Header.Set(header, value)
	}

	var authKey string
	if viper.IsSet(consts.ApiKeyHeader) {
		authKey = viper.GetString(consts.ApiKeyHeader)
	}

	request.Header.Set(consts.ApiKeyHeader, authKey)
	return request, nil
}

func GetBaseURL() string {
	hostname := viper.GetString("hostname")
	scheme := viper.GetString("scheme")
	return fmt.Sprintf("%s://%s", scheme, hostname)
}
