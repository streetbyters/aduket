package aduket

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http/httptest"

	"github.com/labstack/echo"
)

type responseBody []byte

func JSONResponse(j interface{}) responseBody {
	jsonBytes, _ := json.Marshal(j)
	return jsonBytes
}

func XMLResponse(x interface{}) responseBody {
	xmlBytes, _ := xml.Marshal(x)
	return xmlBytes
}

func StringResponse(s string) responseBody {
	return []byte(s)
}

func ByteResponse(b []byte) responseBody {
	return b
}

func NoResponse() responseBody {
	return nil
}

func NewServer(httpMethod, path string, statusCode int, response responseBody) (*httptest.Server, *RequestRecorder) {
	requestRecorder := NewRequestRecorder()
	e := createEcho(requestRecorder, httpMethod, path, statusCode, response)
	return httptest.NewServer(e), requestRecorder
}

func createEcho(requestRecorder *RequestRecorder, httpMethod, path string, statusCode int, body responseBody) *echo.Echo {
	e := echo.New()

	e.Add(httpMethod, path, func(ctx echo.Context) error {

		if ctx.Request().Header.Get(echo.HeaderContentType) == echo.MIMEApplicationXML {
			requestRecorder.bindXML(ctx.Request().Body)
		} else if err := ctx.Bind(&requestRecorder.Body); err != nil {
			data, err := ioutil.ReadAll(ctx.Request().Body)
			if err != nil {
				return err
			}
			requestRecorder.setData(data)
		}

		requestRecorder.setParams(ctx.ParamNames(), ctx.ParamValues())
		requestRecorder.setQueryParams(ctx.QueryParams())
		requestRecorder.setFormParams(ctx.Request().Form)

		if body == nil {
			return ctx.NoContent(statusCode)
		}

		ctx.Response().WriteHeader(statusCode)
		_, err := ctx.Response().Write(body)
		if err != nil {
			return err
		}

		return nil
	})

	return e
}
