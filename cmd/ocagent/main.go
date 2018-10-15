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

// Program ocagent collects OpenCensus stats and traces
// to export to a configured backend.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	agentinteractionpb "github.com/census-instrumentation/opencensus-proto/gen-go/agent/interaction/v1"
	agenttracepb "github.com/census-instrumentation/opencensus-proto/gen-go/agent/trace/v1"
	"github.com/census-instrumentation/opencensus-service/cmd/ocagent/exporterparser"
	"github.com/census-instrumentation/opencensus-service/exporter"
	"github.com/census-instrumentation/opencensus-service/interceptor/interaction"
	"github.com/census-instrumentation/opencensus-service/interceptor/opencensus"
	"github.com/census-instrumentation/opencensus-service/internal"
	"github.com/census-instrumentation/opencensus-service/spanreceiver"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
)

func main() {
	exportersYAMLConfigFile := flag.String("exporters-yaml", "config.yaml", "The YAML file with the configurations for the various exporters")
	flag.Parse()

	yamlBlob, err := ioutil.ReadFile(*exportersYAMLConfigFile)
	if err != nil {
		log.Fatalf("Cannot read the YAML file %v error: %v", exportersYAMLConfigFile, err)
	}
	ownConfig, err := parseOCAgentConfig(yamlBlob)
	if err != nil {
		log.Fatalf("Failed to parse own configuration %v error: %v", exportersYAMLConfigFile, err)
	}

	traceExporters, _, closeFns := exporterparser.ExportersFromYAMLConfig(yamlBlob)

	commonSpanReceiver := exporter.OCExportersToTraceExporter(traceExporters...)

	ocInterceptorAddr := ownConfig.openCensusInterceptorAddressOrDefault()
	ocInterceptorDoneFn, err := runOCInterceptor(ocInterceptorAddr, commonSpanReceiver)
	if err != nil {
		log.Fatal(err)
	}
	closeFns = append(closeFns, ocInterceptorDoneFn)

	if ownConfig.interactionInterceptorEnabled() {
		iaInterceptorGrpcAddr := ownConfig.interactionInterceptorGrpcAddressOrDefault()
		iaInterceptorHttpAddr := ownConfig.interactionInterceptorHttpAddressOrDefault()
		iaInterceptorDoneFn, err := runIaInterceptor(iaInterceptorGrpcAddr, iaInterceptorHttpAddr, commonSpanReceiver)
		if err != nil {
			log.Fatal(err)
		}
		closeFns = append(closeFns, iaInterceptorDoneFn)
	}

	// Add other interceptors here as they are implemented

	// Always cleanup finally
	defer func() {
		for _, closeFn := range closeFns {
			if closeFn != nil {
				closeFn()
			}
		}
	}()

	signalsChan := make(chan os.Signal)
	signal.Notify(signalsChan, os.Interrupt)

	// Wait for the closing signal
	<-signalsChan
}

func runOCInterceptor(addr string, sr spanreceiver.SpanReceiver) (doneFn func(), err error) {
	oci, err := ocinterceptor.New(sr, ocinterceptor.WithSpanBufferPeriod(800*time.Millisecond))
	if err != nil {
		return nil, fmt.Errorf("Failed to create the OpenCensus interceptor: %v", err)
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("Cannot bind to address %q: %v", addr, err)
	}
	srv := internal.GRPCServerWithObservabilityEnabled()
	if err := view.Register(internal.AllViews...); err != nil {
		return nil, fmt.Errorf("Failed to register internal.AllViews: %v", err)
	}
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		return nil, fmt.Errorf("Failed to register ocgrpc.DefaultServerViews: %v", err)
	}

	agenttracepb.RegisterTraceServiceServer(srv, oci)
	go func() {
		log.Printf("Running OpenCensus interceptor as a gRPC service at %q", addr)
		_ = srv.Serve(ln)
	}()
	doneFn = func() { _ = ln.Close() }
	return doneFn, nil
}

func runIaInterceptor(grpcAddr string, httpAddr string, sr spanreceiver.SpanReceiver) (doneFn func(), err error) {
	iai, err := iainterceptor.New(sr, iainterceptor.WithSpanBufferPeriod(800*time.Millisecond))
	if err != nil {
		return nil, fmt.Errorf("Failed to create interaction interceptor: %v", err)
	}

	ln, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return nil, fmt.Errorf("Cannot bind to address %q: %v", grpcAddr, err)
	}
	srv := internal.GRPCServerWithObservabilityEnabled()
	if err := view.Register(internal.AllViews...); err != nil {
		return nil, fmt.Errorf("Failed to register internal.AllViews: %v", err)
	}
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		return nil, fmt.Errorf("Failed to register ocgrpc.DefaultServerViews: %v", err)
	}

	agentinteractionpb.RegisterInteractionServiceServer(srv, iai)
	go func() {
		log.Printf("Running interaction gRPC service at %q", grpcAddr)
		_ = srv.Serve(ln)
	}()
	doneFn = func() {}
	return doneFn, nil
}
