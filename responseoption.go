package aduket

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"time"
)

type ResponseRuleOption func(*responseRule)

func StatusCode(statusCode int) ResponseRuleOption {
	return func(r *responseRule) {
		r.statusCode = statusCode
	}
}

func JSONBody(body interface{}) ResponseRuleOption {
	return func(r *responseRule) {
		r.body = jsonToResponseBody(body)
	}
}

func XMLBody(body interface{}) ResponseRuleOption {
	return func(r *responseRule) {
		r.body = xmlToResponseBody(body)
	}
}

func StringBody(str string) ResponseRuleOption {
	return func(r *responseRule) {
		r.body = stringToResponseBody(str)
	}
}

func ByteBody(b []byte) ResponseRuleOption {
	return func(r *responseRule) {
		r.body = b
	}
}

func Header(header http.Header) ResponseRuleOption {
	return func(r *responseRule) {
		r.header = header
	}
}

func Timeout(duration time.Duration) ResponseRuleOption {
	return func(r *responseRule) {
		r.timeout = duration
	}
}

func jsonToResponseBody(j interface{}) responseBody {
	jsonBytes, _ := json.Marshal(j)
	return jsonBytes
}

func xmlToResponseBody(x interface{}) responseBody {
	xmlBytes, _ := xml.Marshal(x)
	return xmlBytes
}

func stringToResponseBody(s string) responseBody {
	return []byte(s)
}
