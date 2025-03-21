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
	"time"

	"vitess.io/vitess/go/vt/discovery"
	"vitess.io/vitess/go/vt/topo/topoproto"
	"vitess.io/vitess/go/vt/vtenv"
	"vitess.io/vitess/go/vt/vttablet/tmclient"

	"vitess.io/vitess/go/sqltypes"
	"vitess.io/vitess/go/vt/logutil"
	"vitess.io/vitess/go/vt/topo"
	"vitess.io/vitess/go/vt/topo/memorytopo"
	"vitess.io/vitess/go/vt/wrangler"

	querypb "vitess.io/vitess/go/vt/proto/query"
	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
)

func TestPermissions(t *testing.T) {
	delay := discovery.GetTabletPickerRetryDelay()
	defer func() {
		discovery.SetTabletPickerRetryDelay(delay)
	}()
	discovery.SetTabletPickerRetryDelay(5 * time.Millisecond)

	// Initialize our environment
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	ts := memorytopo.NewServer(ctx, "cell1", "cell2")
	wr := wrangler.New(vtenv.NewTestEnv(), logutil.NewConsoleLogger(), ts, tmclient.NewTabletManagerClient())
	vp := NewVtctlPipe(ctx, t, ts)
	defer vp.Close()

	primary := NewFakeTablet(t, wr, "cell1", 0, topodatapb.TabletType_PRIMARY, nil)
	replica := NewFakeTablet(t, wr, "cell1", 1, topodatapb.TabletType_REPLICA, nil)

	// mark the primary inside the shard
	_, err := ts.UpdateShardFields(ctx, primary.Tablet.Keyspace, primary.Tablet.Shard, func(si *topo.ShardInfo) error {
		si.PrimaryAlias = primary.Tablet.Alias
		return nil
	})
	if err != nil {
		require("UpdateShardFields failed: %v", err)
	}

	// primary will be asked for permissions
	primary.FakeMysqlDaemon.FetchSuperQueryMap = map[string]*sqltypes.Result{
		"SELECT * FROM mysql.user ORDER BY host, user": {
			Fields: []*querypb.Field{
				{
					Name: "Host",
					Type: sqltypes.Char,
				},
				{
					Name: "User",
					Type: sqltypes.Char,
				},
				{
					Name: "Password",
					Type: sqltypes.Char,
				},
				{
					Name: "Select_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Insert_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Update_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Delete_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Create_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Drop_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Reload_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Shutdown_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Process_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "File_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Grant_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "References_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Index_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Alter_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Show_db_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Super_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Create_tmp_table_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Lock_tables_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Execute_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Repl_slave_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Repl_client_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Create_view_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Show_view_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Create_routine_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Alter_routine_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Create_user_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Event_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Trigger_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Create_tablespace_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "ssl_type",
					Type: sqltypes.Char,
				},
				{
					Name: "ssl_cipher",
					Type: 252,
				},
				{
					Name: "x509_issuer",
					Type: 252,
				},
				{
					Name: "x509_subject",
					Type: 252,
				},
				{
					Name: "max_questions",
					Type: 3,
				},
				{
					Name: "max_updates",
					Type: 3,
				},
				{
					Name: "max_connections",
					Type: 3,
				},
				{
					Name: "max_user_connections",
					Type: 3,
				},
				{
					Name: "plugin",
					Type: sqltypes.Char,
				},
				{
					Name: "authentication_string",
					Type: 252,
				},
				{
					Name: "password_expired",
					Type: sqltypes.Char,
				},
				{
					Name: "is_role",
					Type: sqltypes.Char,
				}},
			RowsAffected: 0x6,
			InsertID:     0x0,
			Rows: [][]sqltypes.Value{
				{
					sqltypes.NewVarBinary("test_host1"),
					sqltypes.NewVarBinary("test_user1"),
					sqltypes.NewVarBinary("test_password1"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N")},
				{
					sqltypes.NewVarBinary("test_host2"),
					sqltypes.NewVarBinary("test_user2"),
					sqltypes.NewVarBinary("test_password2"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N")},
				{
					sqltypes.NewVarBinary("test_host3"),
					sqltypes.NewVarBinary("test_user3"),
					sqltypes.NewVarBinary("test_password3"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N")},
				{
					sqltypes.NewVarBinary("test_host4"),
					sqltypes.NewVarBinary("test_user4"),
					sqltypes.NewVarBinary("test_password4"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary("0"),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary(""),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("N"),
				},
			},
		},
		"SELECT * FROM mysql.db ORDER BY host, db, user": {
			Fields: []*querypb.Field{
				{
					Name: "Host",
					Type: sqltypes.Char,
				},
				{
					Name: "Db",
					Type: sqltypes.Char,
				},
				{
					Name: "User",
					Type: sqltypes.Char,
				},
				{
					Name: "Select_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Insert_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Update_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Delete_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Create_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Drop_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Grant_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "References_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Index_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Alter_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Create_tmp_table_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Lock_tables_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Create_view_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Show_view_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Create_routine_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Alter_routine_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Execute_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Event_priv",
					Type: sqltypes.Char,
				},
				{
					Name: "Trigger_priv",
					Type: sqltypes.Char,
				},
			},
			RowsAffected: 0,
			InsertID:     0,
			Rows: [][]sqltypes.Value{
				{
					sqltypes.NewVarBinary("test_host"),
					sqltypes.NewVarBinary("test_db"),
					sqltypes.NewVarBinary("test_user"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("N"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
					sqltypes.NewVarBinary("Y"),
				},
			},
		},
	}
	primary.StartActionLoop(t, wr)
	defer primary.StopActionLoop(t)

	// Make a two-level-deep copy, so we can make them diverge later.
	user := *primary.FakeMysqlDaemon.FetchSuperQueryMap["SELECT * FROM mysql.user ORDER BY host, user"]
	user.Fields = append([]*querypb.Field{}, user.Fields...)

	// replica will be asked for permissions
	replica.FakeMysqlDaemon.FetchSuperQueryMap = map[string]*sqltypes.Result{
		"SELECT * FROM mysql.user ORDER BY host, user":   &user,
		"SELECT * FROM mysql.db ORDER BY host, db, user": primary.FakeMysqlDaemon.FetchSuperQueryMap["SELECT * FROM mysql.db ORDER BY host, db, user"],
	}
	replica.FakeMysqlDaemon.SetReplicationSourceInputs = append(replica.FakeMysqlDaemon.SetReplicationSourceInputs, topoproto.MysqlAddr(primary.Tablet))
	replica.FakeMysqlDaemon.ExpectedExecuteSuperQueryList = []string{
		// These 3 statements come from tablet startup
		"STOP REPLICA",
		"FAKE SET SOURCE",
		"START REPLICA",
	}
	replica.StartActionLoop(t, wr)
	defer replica.StopActionLoop(t)

	// Overwrite with the correct value to make sure it passes.
	replica.FakeMysqlDaemon.FetchSuperQueryMap["SELECT * FROM mysql.user ORDER BY host, user"].Fields[0] = &querypb.Field{
		Name: "Host",
		Type: sqltypes.Char,
	}

	// run ValidatePermissionsKeyspace, this should work
	if err := vp.Run([]string{"ValidatePermissionsKeyspace", primary.Tablet.Keyspace}); err != nil {
		require("ValidatePermissionsKeyspace failed: %v", err)
	}

	// modify one field, this should fail
	replica.FakeMysqlDaemon.FetchSuperQueryMap["SELECT * FROM mysql.user ORDER BY host, user"].Fields[0] = &querypb.Field{
		Name: "Wrong",
		Type: sqltypes.Char,
	}

	// run ValidatePermissionsKeyspace again, this should now fail
	if err := vp.Run([]string{"ValidatePermissionsKeyspace", primary.Tablet.Keyspace}); err == nil || !strings.Contains(err.Error(), "has an extra user") {
		require("ValidatePermissionsKeyspace has unexpected err: %v", err)
	}

}
