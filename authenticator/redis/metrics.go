// Copyright 2019 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redis

import "github.com/dfuse-io/dmetrics"

var metricset = dmetrics.NewSet()

var ContextWithCutoffCanceledCounter = metricset.NewCounter("context_with_cutoff_cancel", "number of context canceled because of a cutoff")
var ContextWithCutoffCounter = metricset.NewCounter("context_with_cutoff", "number context with cutoff created")

func init() {
	metricset.Register()
}
