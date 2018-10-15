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

import "time"

type IaOption interface {
	WithIaInterceptor(*IaInterceptor)
}

type spanBufferPeriod struct {
	period time.Duration
}

var _ IaOption = (*spanBufferPeriod)(nil)

func (sfd *spanBufferPeriod) WithIaInterceptor(oci *IaInterceptor) {
	oci.spanBufferPeriod = sfd.period
}

// WithSpanBufferPeriod is an option that allows one to configure
// the period that spans are buffered for before the IaInterceptor
// sends them to its SpanReceiver.
func WithSpanBufferPeriod(period time.Duration) IaOption {
	return &spanBufferPeriod{period: period}
}

type spanBufferCount int

var _ IaOption = (*spanBufferCount)(nil)

func (spc spanBufferCount) WithIaInterceptor(oci *IaInterceptor) {
	oci.spanBufferCount = int(spc)
}

// WithSpanBufferCount is an option that allows one to configure
// the number of spans that are buffered before the IaInterceptor
// send them to its SpanReceiver.
func WithSpanBufferCount(count int) IaOption {
	return spanBufferCount(count)
}
