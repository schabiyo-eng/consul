// Copyright IBM Corp. 2024, 2026
// SPDX-License-Identifier: BUSL-1.1

package protohcl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/consul/internal/protohcl/testproto"
)

func TestSortAnyFieldsLast(t *testing.T) {
	plain := (&testproto.Primitives{}).ProtoReflect().Descriptor().Fields().ByName("string_val")
	anyDesc := (&testproto.DynamicWellKnown{}).ProtoReflect().Descriptor().Fields().ByName("any_val")
	require.False(t, isAnyField(plain))
	require.True(t, isAnyField(anyDesc))

	// Large enough that the sort does not fall back to insertion sort, which
	// previously masked an inconsistent comparator.
	var fields []*IterField
	var wantPlain, wantAny []string
	for i := 0; i < 40; i++ {
		name := fmt.Sprintf("f%02d", i)
		f := &IterField{Name: name, Desc: plain}
		if i%3 == 0 {
			f.Desc = anyDesc
			wantAny = append(wantAny, name)
		} else {
			wantPlain = append(wantPlain, name)
		}
		fields = append(fields, f)
	}

	sortAnyFieldsLast(fields)

	var got []string
	for _, f := range fields {
		got = append(got, f.Name)
	}
	require.Equal(t, append(wantPlain, wantAny...), got)
}
