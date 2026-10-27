// Copyright IBM Corp. 2024, 2026
// SPDX-License-Identifier: BUSL-1.1

package configentry

import (
	"cmp"
	"slices"

	"github.com/hashicorp/consul/agent/structs"
)

func SortSlice(configs []structs.ConfigEntry) {
	slices.SortStableFunc(configs, Compare)
}

func Less(first structs.ConfigEntry, second structs.ConfigEntry) bool {
	return Compare(first, second) < 0
}

func Compare(first structs.ConfigEntry, second structs.ConfigEntry) int {
	if c := cmp.Compare(first.GetKind(), second.GetKind()); c != 0 {
		return c
	}

	if first.GetEnterpriseMeta().LessThan(second.GetEnterpriseMeta()) {
		return -1
	}
	if second.GetEnterpriseMeta().LessThan(first.GetEnterpriseMeta()) {
		return 1
	}

	return cmp.Compare(first.GetName(), second.GetName())
}

func EqualID(e1, e2 structs.ConfigEntry) bool {
	return e1.GetKind() == e2.GetKind() &&
		e1.GetEnterpriseMeta().IsSame(e2.GetEnterpriseMeta()) &&
		e1.GetName() == e2.GetName()
}

func SameHash(e1, e2 structs.ConfigEntry) bool {
	if e1.GetHash() == 0 || e2.GetHash() == 0 {
		return false
	}
	return e1.GetHash() == e2.GetHash()
}
