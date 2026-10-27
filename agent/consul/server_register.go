// Copyright IBM Corp. 2024, 2026
// SPDX-License-Identifier: BUSL-1.1

package consul

import "github.com/hashicorp/consul/logging"

func init() {
	registerEndpoint(func(s *Server) any { return &ACL{s, s.loggers.Named(logging.ACL)} })
	registerEndpoint(func(s *Server) any { return &Catalog{s, s.loggers.Named(logging.Catalog)} })
	registerEndpoint(func(s *Server) any { return NewCoordinate(s, s.logger) })
	registerEndpoint(func(s *Server) any { return &ConfigEntry{s, s.loggers.Named(logging.ConfigEntry)} })
	registerEndpoint(func(s *Server) any { return &ConnectCA{srv: s, logger: s.loggers.Named(logging.Connect)} })
	registerEndpoint(func(s *Server) any { return &FederationState{s} })
	registerEndpoint(func(s *Server) any { return &DiscoveryChain{s} })
	registerEndpoint(func(s *Server) any { return &Health{s, s.loggers.Named(logging.Health)} })
	registerEndpoint(func(s *Server) any { return &Intention{s, s.loggers.Named(logging.Intentions)} })
	registerEndpoint(func(s *Server) any { return &Internal{s, s.loggers.Named(logging.Internal)} })
	registerEndpoint(func(s *Server) any { return &KVS{s, s.loggers.Named(logging.KV)} })
	registerEndpoint(func(s *Server) any { return &Operator{s, s.loggers.Named(logging.Operator)} })
	registerEndpoint(func(s *Server) any { return &PreparedQuery{s, s.loggers.Named(logging.PreparedQuery)} })
	registerEndpoint(func(s *Server) any { return &Session{s, s.loggers.Named(logging.Session)} })
	registerEndpoint(func(s *Server) any { return &Status{s} })
	registerEndpoint(func(s *Server) any { return &Txn{s, s.loggers.Named(logging.Transaction)} })
}
