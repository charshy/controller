// Copyright 2016 IBM Corporation
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package nginx

import (
	"io"
)

// MockGenerator implements interface
type MockGenerator struct {
	GenerateString string
	GenerateError  error
}

// Generate mocks method
func (m *MockGenerator) Generate(w io.Writer, id string) error {
	w.Write([]byte(m.GenerateString))
	return m.GenerateError
}
