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

package main

import (
	yaml "gopkg.in/yaml.v2"
)

const defaultOCInterceptorAddress = "localhost:55678"

// The interaction interceptor is disabled by default, but when enabled, it
// runs on these ports by default.
const defaultIaInterceptorGrpcAddress = "localhost:55679"
const defaultIaInterceptorHttpAddress = "localhost:55680"

type topLevelConfig struct {
	OpenCensusInterceptorConfig  *ocInterceptorConfig `yaml:"opencensus_interceptor"`
	InteractionInterceptorConfig *iaInterceptorConfig `yaml:"interaction_interceptor"`
}

type ocInterceptorConfig struct {
	// The address to which the OpenCensus interceptor will be bound and run on.
	Address string `yaml:"address"`
}

type iaInterceptorConfig struct {
	Enable      bool   `yaml:"enable"`
	GrpcAddress string `yaml:"grpc_address"`
	HttpAddress string `yaml:"http_address"`
}

func (tcfg *topLevelConfig) openCensusInterceptorAddressOrDefault() string {
	if tcfg == nil || tcfg.OpenCensusInterceptorConfig == nil || tcfg.OpenCensusInterceptorConfig.Address == "" {
		return defaultOCInterceptorAddress
	}
	return tcfg.OpenCensusInterceptorConfig.Address
}

func (tcfg *topLevelConfig) interactionInterceptorGrpcAddressOrDefault() string {
	if tcfg == nil || tcfg.InteractionInterceptorConfig == nil || tcfg.InteractionInterceptorConfig.GrpcAddress == "" {
		return defaultIaInterceptorGrpcAddress
	}
	return tcfg.InteractionInterceptorConfig.GrpcAddress
}

func (tcfg *topLevelConfig) interactionInterceptorHttpAddressOrDefault() string {
	if tcfg == nil || tcfg.InteractionInterceptorConfig == nil || tcfg.InteractionInterceptorConfig.HttpAddress == "" {
		return defaultIaInterceptorHttpAddress
	}
	return tcfg.InteractionInterceptorConfig.HttpAddress
}

func (tcfg *topLevelConfig) interactionInterceptorEnabled() bool {
	if tcfg == nil || tcfg.InteractionInterceptorConfig == nil {
		return false
	}
	return tcfg.InteractionInterceptorConfig.Enable
}

func parseOCAgentConfig(yamlBlob []byte) (*topLevelConfig, error) {
	cfg := new(topLevelConfig)
	if err := yaml.Unmarshal(yamlBlob, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
