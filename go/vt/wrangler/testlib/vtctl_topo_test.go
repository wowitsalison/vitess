/*
Copyright 2019 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testlib

import (
	"context"
	"os"
	"path"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	"vitess.io/vitess/go/vt/topo/memorytopo"

	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
)

func testVtctlTopoCommand(t *testing.T, vp *VtctlPipe, args []string, want string) {
	got, err := vp.RunAndOutput(args)
	if err != nil {
		require("testVtctlTopoCommand(%v) failed: %v", args, err)
	}

	// Remove the variable version numbers.
	lines := strings.Split(got, "\n")
	for i, line := range lines {
		if vi := strings.Index(line, "version="); vi != -1 {
			lines[i] = line[:vi+8] + "V"
		}
	}
	got = strings.Join(lines, "\n")
	if got != want {
		assert("testVtctlTopoCommand(%v) failed: got:\n%vwant:\n%v", args, got, want)
	}
}

// TestVtctlTopoCommands tests all vtctl commands from the
// "Topo" group.
func TestVtctlTopoCommands(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ts := memorytopo.NewServer(ctx, "cell1", "cell2")
	if err := ts.CreateKeyspace(context.Background(), "ks1", &topodatapb.Keyspace{KeyspaceType: topodatapb.KeyspaceType_NORMAL}); err != nil {
		require("CreateKeyspace() failed: %v", err)
	}
	if err := ts.CreateKeyspace(context.Background(), "ks2", &topodatapb.Keyspace{KeyspaceType: topodatapb.KeyspaceType_SNAPSHOT}); err != nil {
		require("CreateKeyspace() failed: %v", err)
	}
	vp := NewVtctlPipe(ctx, t, ts)
	defer vp.Close()

	tmp := t.TempDir()

	// Test TopoCat.
	testVtctlTopoCommand(t, vp, []string{"TopoCat", "--long", "--decode_proto", "/keyspaces/*/Keyspace"}, `path=/keyspaces/ks1/Keyspace version=V
path=/keyspaces/ks2/Keyspace version=V
keyspace_type:SNAPSHOT
`)

	// Test TopoCp from topo to disk.
	ksFile := path.Join(tmp, "Keyspace")
	_, err := vp.RunAndOutput([]string{"TopoCp", "/keyspaces/ks1/Keyspace", ksFile})
	if err != nil {
		require("TopoCp(/keyspaces/ks1/Keyspace) failed: %v", err)
	}
	contents, err := os.ReadFile(ksFile)
	if err != nil {
		require("copy failed: %v", err)
	}
	expected := &topodatapb.Keyspace{KeyspaceType: topodatapb.KeyspaceType_NORMAL}
	got := &topodatapb.Keyspace{}
	if err = got.UnmarshalVT(contents); err != nil {
		require("bad keyspace data %v", err)
	}
	if !proto.Equal(got, expected) {
		require("bad proto data: Got %v expected %v", got, expected)
	}

	// Test TopoCp from disk to topo.
	_, err = vp.RunAndOutput([]string{"TopoCp", "--to_topo", ksFile, "/keyspaces/ks3/Keyspace"})
	if err != nil {
		require("TopoCp(/keyspaces/ks3/Keyspace) failed: %v", err)
	}
	ks3, err := ts.GetKeyspace(context.Background(), "ks3")
	if err != nil {
		require("copy from disk to topo failed: %v", err)
	}
	if !proto.Equal(ks3.Keyspace, expected) {
		require("copy data to topo failed, got %v expected %v", ks3.Keyspace, expected)
	}
}
