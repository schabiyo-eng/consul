// Copyright IBM Corp. 2024, 2026
// SPDX-License-Identifier: BUSL-1.1

package topology

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNode_SortedWorkloads(t *testing.T) {
	gw := &Workload{ID: NewID("a-gateway", "", ""), IsMeshGateway: true}
	app1 := &Workload{ID: NewID("b-app", "", "")}
	app2 := &Workload{ID: NewID("c-app", "", "")}
	want := []*Workload{app1, app2, gw}

	for name, in := range map[string][]*Workload{
		"gateway first": {gw, app1, app2},
		"gateway last":  {app2, app1, gw},
		"gateway mid":   {app2, gw, app1},
	} {
		t.Run(name, func(t *testing.T) {
			n := &Node{Workloads: in}
			require.Equal(t, want, n.SortedWorkloads())
		})
	}
}
