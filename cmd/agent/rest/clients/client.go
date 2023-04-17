package clients

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

func CreateClient(baseURL string) *resty.Client {
	c := resty.New().SetBaseURL(fmt.Sprintf("http://%s", baseURL)).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(2 * time.Second)
	onBeforeRequest(c)
	return c
}

func onBeforeRequest(c *resty.Client) {
	setXRealIPToHeader(c)
}

func setXRealIPToHeader(c *resty.Client) {
	c.OnBeforeRequest(
		func(resty *resty.Client, r *resty.Request) error {
			r.SetHeader("X-Real-IP", "127.0.0.1/24")
			return nil
		})
}
