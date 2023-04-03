package clients

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

func CreateClient(baseURL string) *resty.Client {
	return resty.New().SetBaseURL(fmt.Sprintf("http://%s", baseURL)).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(2 * time.Second)
}
