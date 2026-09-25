// Copyright 2026 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"testing"
)

// TestInitMatchersCompatVerbose ensures --verbose does not panic while
// initializing the promslog level (regression test: promslogConfig.Level
// is nil until explicitly allocated, and Set() dereferences it).
func TestInitMatchersCompatVerbose(t *testing.T) {
	prevVerbose := verbose
	defer func() { verbose = prevVerbose }()

	verbose = true
	if err := initMatchersCompat(nil); err != nil {
		t.Fatalf("initMatchersCompat with verbose=true failed: %v", err)
	}
}
