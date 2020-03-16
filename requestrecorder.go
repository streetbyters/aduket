// Copyright 2020 StreetByters Community
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aduket

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/clbanning/mxj"
	"github.com/labstack/echo"
)

type RequestRecorder struct {
	Body              Body
	Header            http.Header
	Data              []byte
	Params            map[string]string
	QueryParams       url.Values
	FormParams        url.Values
	isRequestReceived bool
}

type Body map[string]interface{}

func NewRequestRecorder() *RequestRecorder {
	requestRecorder := &RequestRecorder{}
	requestRecorder.Body = make(Body)
	requestRecorder.Params = make(map[string]string)
	return requestRecorder
}

func (r *RequestRecorder) saveContext(ctx echo.Context) error {
	r.isRequestReceived = true

	if ctx.Request().Header.Get(echo.HeaderContentType) == echo.MIMEApplicationXML {
		r.bindXML(ctx.Request().Body)
		return nil
	}

	defaultBinder := new(echo.DefaultBinder)

	if err := defaultBinder.Bind(&r.Body, ctx); err != nil {
		data, err := ioutil.ReadAll(ctx.Request().Body)
		if err != nil {
			return err
		}

		r.setData(data)
	}

	r.setParams(ctx.ParamNames(), ctx.ParamValues())
	r.setQueryParams(ctx.QueryParams())
	r.setFormParams(ctx.Request().Form)
	r.setHeader(ctx.Request().Header)

	return nil

}

func (r *RequestRecorder) setQueryParams(queryParams url.Values) {
	r.QueryParams = queryParams
}

func (r *RequestRecorder) setParams(paramNames, paramValues []string) {
	for index, name := range paramNames {
		r.Params[name] = paramValues[index]
	}
}

func (r *RequestRecorder) setFormParams(formParams url.Values) {
	r.FormParams = formParams
}

func (r *RequestRecorder) setData(b []byte) {
	r.Data = b
}

func (r *RequestRecorder) setHeader(header http.Header) {
	r.Header = header.Clone()
}

func (r *RequestRecorder) bindXML(from io.ReadCloser) error {
	body, err := ioutil.ReadAll(from)
	if err != nil {
		return err
	}

	mv, err := mxj.NewMapXml(body)
	if err != nil {
		return err
	}

	r.Body = mv.Old()

	return nil
}
