package sign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Bilibili
func TestSignMd5(t *testing.T) {
	body := map[string]any{
		"money":        600,
		"out_trade_no": "example-1234567890",
		"product_id":   "com.example.test.item01",
		"notify_url":   "http://demo.com/notify/example",
	}
	r, e := Generate(body, WithSecret("pHcNbNb4vtvn8HXT"), WithSignAlg(Md5))
	assert.NoError(t, e)
	assert.Equal(t, "c736b13ba58ff87d9e8ce5d6190fc154", r)
}

// Douyin
func TestSignHmacSha1WithBase64(t *testing.T) {
	var body = map[string]any{
		"app_id":       1234,
		"access_token": "q3fafa33sHFU+V9h32h0v8weVEH/04hgsrHFHOHNNQOBC9fnwejasubw==",
		"ts":           1555912969,
	}
	r, e := Generate(body,
		WithSecret("1ACCgaRXbazbaVzkd1NzwVc2G7wg1d6G"),
		WithJoiner("&"),
		WithUseKv(true),
		WithSignAlg(HmacSha1WithBase64),
	)
	assert.NoError(t, e)
	assert.Equal(t, "sTFH83DV+Vamgr6SRsC/NNjw0+Q=", r)
}
