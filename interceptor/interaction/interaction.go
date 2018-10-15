// Copyright 2018, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package iainterceptor

import (
	"errors"
	"time"

	"github.com/census-instrumentation/opencensus-service/spanreceiver"
)

type IaInterceptor struct {
	spanSink         spanreceiver.SpanReceiver
	spanBufferPeriod time.Duration
	spanBufferCount  int
}

func New(sr spanreceiver.SpanReceiver, opts ...IaOption) (*IaInterceptor, error) {
	if sr == nil {
		return nil, errors.New("needs a non-nil spanReceiver")
	}
	iai := &IaInterceptor{spanSink: sr}
	for _, opt := range opts {
		opt.WithIaInterceptor(iai)
	}
	return iai, nil
}
