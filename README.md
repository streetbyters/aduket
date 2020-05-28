<h1 align="center">Aduket</h1>
<p align="center">
    <img height="200px" src="assets/aduket.gif">
</p>
<p align="center">
    <i>Straight-forward HTTP client testing, assertions included!</i>
</p>

<p align="center">
  <a href="https://github.com/streetbyters/aduket/actions">
    <img src="https://img.shields.io/github/workflow/status/streetbyters/aduket/Go" />
  </a>
  <a href="https://codecov.io/gh/streetbyters/aduket">
    <img src="https://codecov.io/gh/streetbyters/aduket/branch/master/graph/badge.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/streetbyters/aduket">
    <img src="https://goreportcard.com/badge/github.com/streetbyters/aduket" />
  </a>
  <a href="https://github.com/streetbyters/aduket/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/streetbyters/aduket.svg">
  </a>
</p>

Simple `httptest.Server` wrapper with a little request recorder spice on it. No special DSL, no complex API to learn. Just create a server and fire your request like an **Hadouken** then assert it.

## Why?

Aduket's raison d'etre is to make you able to test your HTTP client related logic in ease _or for some reason you just need a mock server to make it respond in a way you need_.

## What?

Aduket currently provides following utilities to make your testing process easier and faster:

- Lean way to spin up a mock HTTP server to imitate different responses _(even timeouts!)_.
- Assertion helpers to validate if you're sending the correct request.

## LICENSE

Copyright 2020 StreetByters Community

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
