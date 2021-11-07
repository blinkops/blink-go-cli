package utils

import (
	"fmt"
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
	if viper.IsSet("BLINK-API-KEY") {
		authKey = viper.GetString("BLINK-API-KEY")
	}

	request.Header.Set("BLINK-API-KEY", authKey)
	return request, nil
}

func GetBaseURL() string {
	hostname := viper.GetString("hostname")
	scheme := viper.GetString("scheme")
	return fmt.Sprintf("%s://%s", scheme, hostname)
}
