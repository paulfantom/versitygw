// Copyright 2026 Versity Software
// This file is licensed under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

//go:build !windows

package meta

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestSideCarStoreAttributeFilePerm(t *testing.T) {
	tests := []struct {
		name string
		perm *fs.FileMode
		want fs.FileMode
	}{
		{name: "default", want: 0644},
		// 0664 is not reachable through a 022 umask, so it proves the mode is
		// applied explicitly
		{name: "group writable", perm: new(fs.FileMode(0664)), want: 0664},
		{name: "owner only", perm: new(fs.FileMode(0600)), want: 0600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			s, err := NewSideCar(dir)
			if err != nil {
				t.Fatal(err)
			}
			if tt.perm != nil {
				s = s.WithFilePerm(*tt.perm)
			}

			for _, object := range []string{"", "obj"} {
				// the second store replaces the attribute file
				for range 2 {
					if err := s.StoreAttribute(nil, "bucket", object, "etag", []byte("etag")); err != nil {
						t.Fatal(err)
					}
				}

				fi, err := os.Stat(filepath.Join(dir, "bucket", object, sidecarmeta, "etag"))
				if err != nil {
					t.Fatal(err)
				}
				if got := fi.Mode().Perm(); got != tt.want {
					t.Errorf("object %q: expected attribute file mode %v, got %v", object, tt.want, got)
				}
			}
		})
	}
}
