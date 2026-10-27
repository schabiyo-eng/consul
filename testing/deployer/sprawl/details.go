// Copyright IBM Corp. 2024, 2026
// SPDX-License-Identifier: BUSL-1.1

package sprawl

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	retry "github.com/avast/retry-go"

	"github.com/hashicorp/consul/api"
)

// PrintDetails will dump relevant addressing and naming data to the logger for
// human interaction purposes.
func (s *Sprawl) PrintDetails() error {
	det := logDetails{
		TopologyID: s.topology.ID,
	}

	for _, cluster := range s.topology.Clusters {
		client := s.clients[cluster.Name]

		var cfg *api.RaftConfiguration
		var err error
		err = retry.Do(
			func() error {
				cfg, err = client.Operator().RaftGetConfiguration(nil)
				if err != nil {
					return fmt.Errorf("error get raft config: %w", err)
				}
				return nil
			},
			retry.MaxDelay(5*time.Second),
			retry.Attempts(15),
		)
		if err != nil {
			return fmt.Errorf("could not get raft config for cluster %q: %w", cluster.Name, err)
		}

		var leaderNode string
		for _, svr := range cfg.Servers {
			if svr.Leader {
				leaderNode = strings.TrimSuffix(svr.Node, "-pod")
			}
		}

		cd := clusterDetails{
			Name:   cluster.Name,
			Leader: leaderNode,
		}

		for _, node := range cluster.Nodes {
			if node.Disabled {
				continue
			}

			var addrs []string
			for _, addr := range node.Addresses {
				addrs = append(addrs, addr.Network+"="+addr.IPAddress)
			}
			sort.Strings(addrs)

			if node.IsServer() {
				cd.Apps = append(cd.Apps, appDetail{
					Type:        "server",
					Container:   node.DockerName(),
					Addresses:   addrs,
					ExposedPort: node.ExposedPort(8500),
				})
			}

			for _, wrk := range node.Workloads {
				if wrk.IsMeshGateway {
					cd.Apps = append(cd.Apps, appDetail{
						Type:                  "mesh-gateway",
						Container:             node.DockerName(),
						ExposedPort:           node.ExposedPort(wrk.Port),
						ExposedEnvoyAdminPort: node.ExposedPort(wrk.EnvoyAdminPort),
						Addresses:             addrs,
						Service:               wrk.ID.String(),
					})
				} else {
					cd.Apps = append(cd.Apps, appDetail{
						Type:                  "app",
						Container:             node.DockerName(),
						ExposedPort:           node.ExposedPort(wrk.Port),
						ExposedEnvoyAdminPort: node.ExposedPort(wrk.EnvoyAdminPort),
						Addresses:             addrs,
						Service:               wrk.ID.String(),
					})
				}
			}
		}

		det.Clusters = append(det.Clusters, cd)
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', tabwriter.Debug)

	score := map[string]int{
		"server":       0,
		"mesh-gateway": 1,
		"app":          2,
	}

	for _, cluster := range det.Clusters {
		fmt.Fprintf(w, "CLUSTER\tTYPE\tCONTAINER\tNAME\tADDRS\tPORTS\t\n")
		slices.SortFunc(cluster.Apps, func(a, b appDetail) int {
			if c := cmp.Compare(score[a.Type], score[b.Type]); c != 0 {
				return c
			}
			if c := cmp.Compare(a.Container, b.Container); c != 0 {
				return c
			}
			return cmp.Compare(a.Service, b.Service)
		})
		for _, d := range cluster.Apps {
			if d.Type == "server" && d.Container == cluster.Leader {
				d.Type = "leader"
			}
			portStr := "app=" + strconv.Itoa(d.ExposedPort)
			if d.ExposedEnvoyAdminPort > 0 {
				portStr += " envoy=" + strconv.Itoa(d.ExposedEnvoyAdminPort)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n",
				cluster.Name,
				d.Type,
				d.Container,
				d.Service,
				strings.Join(d.Addresses, ", "),
				portStr,
			)
		}
		fmt.Fprintf(w, "\t\t\t\t\t\n")
	}

	w.Flush()

	s.logger.Debug("CURRENT SPRAWL DETAILS", "details", buf.String())

	return nil
}

type logDetails struct {
	TopologyID string
	Clusters   []clusterDetails
}

type clusterDetails struct {
	Name string

	Leader string
	Apps   []appDetail
}

type appDetail struct {
	Type                  string // server|mesh-gateway|app
	Container             string
	Addresses             []string
	ExposedPort           int `json:",omitempty"`
	ExposedEnvoyAdminPort int `json:",omitempty"`
	// just services
	Service string `json:",omitempty"`
}
