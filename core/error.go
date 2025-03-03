package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type (
	PlatError struct {
		StatusCode int
		Header     http.Header
		Body       string
		Code       int         `json:"code"`
		Message    string      `json:"message"`
		Details    interface{} `json:"details"`
	}
)

func (e *PlatError) Error() string {
	var buf bytes.Buffer
	_, _ = fmt.Fprintf(&buf, "error http response:[StatusCode: %d Code: \"%d\"", e.StatusCode, e.Code)
	if e.Message != "" {
		_, _ = fmt.Fprintf(&buf, "\nMessage: %s", e.Message)
	}
	if e.Details != nil {
		var detailBuf bytes.Buffer
		enc := json.NewEncoder(&detailBuf)
		enc.SetIndent("", "  ")
		if err := enc.Encode(e.Details); err == nil {
			_, _ = fmt.Fprint(&buf, "\nDetail:")
			_, _ = fmt.Fprintf(&buf, "\n%s", strings.TrimSpace(detailBuf.String()))
		}
	}
	if len(e.Header) > 0 {
		_, _ = fmt.Fprint(&buf, "\nHeader:")
		for key, value := range e.Header {
			_, _ = fmt.Fprintf(&buf, "\n - %v=%v", key, value)
		}
	}
	_, _ = fmt.Fprintf(&buf, "]")
	return buf.String()
}

//func ResponseError(response *resty.Response) (platErr *PlatError, err error) {
//
//	if response.Header().Get(consts.HeaderRequestResult) != "1" {
//		return
//	}
//
//	platErr = new(PlatError)
//
//	if err = json.Unmarshal(response.Body(), platErr); err != nil {
//		return
//	}
//
//	return
//}
