package http_clt

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type IHttpClientV2 interface {
	// Возвращает базовый HTTPClient
	Client() *http.Client
	// Если ответ = 200, то распарсит его в response и вернет nil
	// Если ответ != 200, то не будет парсить и вернет badStatusHandler(resp)
	// Если произошла ошибка, вернет err
	JSON(method, urlSuffix string, requestOpt, responseOpt any, requestId string, headersOpt map[string]string, ctx context.Context) error
}

func NewHttpClientV2(c *http.Client, m IHttpClientMetrics, baseUrl string) IHttpClientV2 {
	return &httpClientV2{
		clt:              c,
		m:                m,
		baseUrl:          baseUrl,
		badStatusHandler: InvalidResponseStatusReadBodyToError,
	}
}

type httpClientV2 struct {
	clt              *http.Client
	m                IHttpClientMetrics
	baseUrl          string
	badStatusHandler InvalidResponseStatusHandlerFunc
}

func (c *httpClientV2) Client() *http.Client {
	return c.clt
}

func (c *httpClientV2) JSON(method, urlSuffix string, requestOpt, responseOpt any, requestId string, headersOpt map[string]string, ctx context.Context) error {
	var requestReader io.Reader
	if requestOpt != nil {
		requestBytes, err := json.Marshal(requestOpt)
		if err != nil {
			return err
		}
		requestReader = bytes.NewReader(requestBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseUrl+urlSuffix, requestReader)
	if err != nil {
		return err
	}
	if requestOpt != nil {
		req.Header.Set(HeaderContentType, ContentTypeJSON)
	}
	if requestId != "" {
		req.Header.Set(HeaderRequestId, requestId)
	}
	for k, v := range headersOpt {
		req.Header.Set(k, v)
	}

	start := time.Now()
	resp, err := c.clt.Do(req)
	if resp != nil { // regardless of err value
		defer func() {
			_ = resp.Body.Close()
		}()
	}
	if err != nil {
		c.m.IncError()
		return err
	}

	if resp.StatusCode != http.StatusOK {
		c.m.IncBadStatus(resp.StatusCode)
		return c.badStatusHandler(resp)
	}

	// Время считается только для 200 ответов
	c.m.RequestDuration(time.Since(start))

	c.m.IncSuccess()
	if responseOpt != nil { // нужно ли читать ответ?
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(responseOpt); err != nil {
			return err
		}
	}

	return nil
}
