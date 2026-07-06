<h1 align="center">Aduket</h1>
<p align="center">
    <img height="200px" src="assets/aduket.gif">
</p>
<p align="center">
    <i>Straight-forward HTTP client testing, assertions included!</i>
</p>

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/streetbyters/aduket/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/streetbyters/aduket/tree/master)
![GitHub License](https://img.shields.io/github/license/streetbyters/aduket)
[![codecov](https://codecov.io/github/streetbyters/aduket/graph/badge.svg?token=6YY26UIWSU)](https://codecov.io/github/streetbyters/aduket)

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
