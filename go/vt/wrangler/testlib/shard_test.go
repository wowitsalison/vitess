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
	"strings"
	"testing"

	"vitess.io/vitess/go/vt/logutil"
	"vitess.io/vitess/go/vt/topo"
	"vitess.io/vitess/go/vt/topo/memorytopo"
	"vitess.io/vitess/go/vt/topotools"
	"vitess.io/vitess/go/vt/vtenv"
	"vitess.io/vitess/go/vt/vttablet/tmclient"
	"vitess.io/vitess/go/vt/wrangler"

	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
)

func TestDeleteShardCleanup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ts := memorytopo.NewServer(ctx, "cell1", "cell2")
	wr := wrangler.New(vtenv.NewTestEnv(), logutil.NewConsoleLogger(), ts, tmclient.NewTabletManagerClient())
	vp := NewVtctlPipe(ctx, t, ts)
	defer vp.Close()

	// Create a primary, a couple good replicas
	primary := NewFakeTablet(t, wr, "cell1", 0, topodatapb.TabletType_PRIMARY, nil)
	replica := NewFakeTablet(t, wr, "cell1", 1, topodatapb.TabletType_REPLICA, nil)
	remoteReplica := NewFakeTablet(t, wr, "cell2", 2, topodatapb.TabletType_REPLICA, nil)

	// Build keyspace graph
	err := topotools.RebuildKeyspace(context.Background(), logutil.NewConsoleLogger(), ts, primary.Tablet.Keyspace, []string{"cell1", "cell2"}, false)
	if err != nil {
		require("RebuildKeyspaceLocked failed: %v", err)
	}

	// Delete the ShardReplication record in cell2
	if err := ts.DeleteShardReplication(ctx, "cell2", remoteReplica.Tablet.Keyspace, remoteReplica.Tablet.Shard); err != nil {
		require("DeleteShardReplication failed: %v", err)
	}

	// Now try to delete the shard without even_if_serving or
	// recursive flag, should fail on serving check first.
	if err := vp.Run([]string{
		"DeleteShard",
		primary.Tablet.Keyspace + "/" + primary.Tablet.Shard,
	}); err == nil || !strings.Contains(err.Error(), "is still serving, cannot delete it") {
		require("DeleteShard() returned wrong error: %v", err)
	}

	// Now try to delete the shard with even_if_serving, but
	// without recursive flag, should fail on existing tablets.
	if err := vp.Run([]string{
		"DeleteShard",
		"--even_if_serving",
		primary.Tablet.Keyspace + "/" + primary.Tablet.Shard,
	}); err == nil || !strings.Contains(err.Error(), "use -recursive or remove them manually") {
		require("DeleteShard(evenIfServing=true) returned wrong error: %v", err)
	}

	// Now try to delete the shard with even_if_serving and recursive,
	// it should just work.
	if err := vp.Run([]string{
		"DeleteShard",
		"--recursive",
		"--even_if_serving",
		primary.Tablet.Keyspace + "/" + primary.Tablet.Shard,
	}); err != nil {
		require("DeleteShard(recursive=true, evenIfServing=true) should have worked but returned: %v", err)
	}

	// Make sure all tablets are gone.
	for _, ft := range []*FakeTablet{primary, replica, remoteReplica} {
		if _, err := ts.GetTablet(ctx, ft.Tablet.Alias); !topo.IsErrType(err, topo.NoNode) {
			assert("tablet %v is still in topo: %v", ft.Tablet.Alias, err)
		}
	}

	// Make sure the shard is gone.
	if _, err := ts.GetShard(ctx, primary.Tablet.Keyspace, primary.Tablet.Shard); !topo.IsErrType(err, topo.NoNode) {
		assert("shard %v/%v is still in topo: %v", primary.Tablet.Keyspace, primary.Tablet.Shard, err)
	}
}
