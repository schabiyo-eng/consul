// Copyright IBM Corp. 2024, 2026
// SPDX-License-Identifier: BUSL-1.1

package oidcauth

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractStringMetadata(t *testing.T) {
	emptyMap := make(map[string]string)

	tests := map[string]struct {
		allClaims     map[string]any
		claimMappings map[string]string
		expected      map[string]string
		errExpected   bool
	}{
		"empty": {nil, nil, emptyMap, false},
		"all": {
			map[string]any{
				"data1": "foo",
				"data2": "bar",
			},
			map[string]string{
				"data1": "val1",
				"data2": "val2",
			},
			map[string]string{
				"val1": "foo",
				"val2": "bar",
			},
			false,
		},
		"some": {
			map[string]any{
				"data1": "foo",
				"data2": "bar",
			},
			map[string]string{
				"data1": "val1",
				"data3": "val2",
			},
			map[string]string{
				"val1": "foo",
			},
			false,
		},
		"none": {
			map[string]any{
				"data1": "foo",
				"data2": "bar",
			},
			map[string]string{
				"data8": "val1",
				"data9": "val2",
			},
			emptyMap,
			false,
		},

		"nested data": {
			map[string]any{
				"data1": "foo",
				"data2": map[string]any{
					"child": "bar",
				},
				"data3": true,
				"data4": false,
				"data5": float64(7.9),
				"data6": json.Number("-12345"),
				"data7": int(42),
			},
			map[string]string{
				"data1":        "val1",
				"/data2/child": "val2",
				"data3":        "val3",
				"data4":        "val4",
				"data5":        "val5",
				"data6":        "val6",
				"data7":        "val7",
			},
			map[string]string{
				"val1": "foo",
				"val2": "bar",
				"val3": "true",
				"val4": "false",
				"val5": "7",
				"val6": "-12345",
				"val7": "42",
			},
			false,
		},

		"error: a struct isn't stringifiable": {
			map[string]any{
				"data1": map[string]any{
					"child": "bar",
				},
			},
			map[string]string{
				"data1": "val1",
			},
			nil,
			true,
		},
		"error: a slice isn't stringifiable": {
			map[string]any{
				"data1": []any{
					"child", "bar",
				},
			},
			map[string]string{
				"data1": "val1",
			},
			nil,
			true,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			actual, err := extractStringMetadata(nil, test.allClaims, test.claimMappings)
			if test.errExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}
}

func TestExtractListMetadata(t *testing.T) {
	emptyMap := make(map[string][]string)

	tests := map[string]struct {
		allClaims     map[string]any
		claimMappings map[string]string
		expected      map[string][]string
		errExpected   bool
	}{
		"empty": {nil, nil, emptyMap, false},
		"all - singular": {
			map[string]any{
				"data1": "foo",
				"data2": "bar",
			},
			map[string]string{
				"data1": "val1",
				"data2": "val2",
			},
			map[string][]string{
				"val1": {"foo"},
				"val2": {"bar"},
			},
			false,
		},
		"some - singular": {
			map[string]any{
				"data1": "foo",
				"data2": "bar",
			},
			map[string]string{
				"data1": "val1",
				"data3": "val2",
			},
			map[string][]string{
				"val1": {"foo"},
			},
			false,
		},
		"none - singular": {
			map[string]any{
				"data1": "foo",
				"data2": "bar",
			},
			map[string]string{
				"data8": "val1",
				"data9": "val2",
			},
			emptyMap,
			false,
		},

		"nested data - singular": {
			map[string]any{
				"data1": "foo",
				"data2": map[string]any{
					"child": "bar",
				},
				"data3": true,
				"data4": false,
				"data5": float64(7.9),
				"data6": json.Number("-12345"),
				"data7": int(42),
				"data8": []any{ // mixed
					"foo", true, float64(7.9), json.Number("-12345"), int(42),
				},
			},
			map[string]string{
				"data1":        "val1",
				"/data2/child": "val2",
				"data3":        "val3",
				"data4":        "val4",
				"data5":        "val5",
				"data6":        "val6",
				"data7":        "val7",
				"data8":        "val8",
			},
			map[string][]string{
				"val1": {"foo"},
				"val2": {"bar"},
				"val3": {"true"},
				"val4": {"false"},
				"val5": {"7"},
				"val6": {"-12345"},
				"val7": {"42"},
				"val8": {
					"foo", "true", "7", "-12345", "42",
				},
			},
			false,
		},

		"error: a struct isn't stringifiable (singular)": {
			map[string]any{
				"data1": map[string]any{
					"child": map[string]any{
						"inner": "bar",
					},
				},
			},
			map[string]string{
				"data1": "val1",
			},
			nil,
			true,
		},
		"error: a slice isn't stringifiable (singular)": {
			map[string]any{
				"data1": []any{
					"child", []any{"bar"},
				},
			},
			map[string]string{
				"data1": "val1",
			},
			nil,
			true,
		},

		"non-string-slice data (string)": {
			map[string]any{
				"data1": "foo",
			},
			map[string]string{
				"data1": "val1",
			},
			map[string][]string{
				"val1": {"foo"}, // singular values become lists
			},
			false,
		},

		"all - list": {
			map[string]any{
				"data1": []any{"foo", "otherFoo"},
				"data2": []any{"bar", "otherBar"},
			},
			map[string]string{
				"data1": "val1",
				"data2": "val2",
			},
			map[string][]string{
				"val1": {"foo", "otherFoo"},
				"val2": {"bar", "otherBar"},
			},
			false,
		},
		"some - list": {
			map[string]any{
				"data1": []any{"foo", "otherFoo"},
				"data2": map[string]any{
					"child": []any{"bar", "otherBar"},
				},
			},
			map[string]string{
				"data1":        "val1",
				"/data2/child": "val2",
			},
			map[string][]string{
				"val1": {"foo", "otherFoo"},
				"val2": {"bar", "otherBar"},
			},
			false,
		},
		"none - list": {
			map[string]any{
				"data1": []any{"foo"},
				"data2": []any{"bar"},
			},
			map[string]string{
				"data8": "val1",
				"data9": "val2",
			},
			emptyMap,
			false,
		},
		"list omits empty strings": {
			map[string]any{
				"data1": []any{"foo", "", "otherFoo", ""},
				"data2": "",
			},
			map[string]string{
				"data1": "val1",
				"data2": "val2",
			},
			map[string][]string{
				"val1": {"foo", "otherFoo"},
				"val2": {},
			},
			false,
		},

		"nested data - list": {
			map[string]any{
				"data1": []any{"foo"},
				"data2": map[string]any{
					"child": []any{"bar"},
				},
				"data3": []any{true},
				"data4": []any{false},
				"data5": []any{float64(7.9)},
				"data6": []any{json.Number("-12345")},
				"data7": []any{int(42)},
				"data8": []any{ // mixed
					"foo", true, float64(7.9), json.Number("-12345"), int(42),
				},
			},
			map[string]string{
				"data1":        "val1",
				"/data2/child": "val2",
				"data3":        "val3",
				"data4":        "val4",
				"data5":        "val5",
				"data6":        "val6",
				"data7":        "val7",
				"data8":        "val8",
			},
			map[string][]string{
				"val1": {"foo"},
				"val2": {"bar"},
				"val3": {"true"},
				"val4": {"false"},
				"val5": {"7"},
				"val6": {"-12345"},
				"val7": {"42"},
				"val8": {
					"foo", "true", "7", "-12345", "42",
				},
			},
			false,
		},

		"JSONPointer": {
			map[string]any{
				"foo": "a",
				"bar": map[string]any{
					"baz": []string{"x", "y", "z"},
				},
			},
			map[string]string{
				"foo":        "val1",
				"/bar/baz/1": "val2",
			},
			map[string][]string{
				"val1": {"a"},
				"val2": {"y"},
			},
			false,
		},
		"JSONPointer not found": {
			map[string]any{
				"foo": "a",
				"bar": map[string]any{
					"baz": []string{"x", "y", "z"},
				},
			},
			map[string]string{
				"foo":           "val1",
				"/bar/XXX/1243": "val2",
			},
			map[string][]string{
				"val1": {"a"},
			},
			false,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			actual, err := extractListMetadata(nil, test.allClaims, test.claimMappings)
			if test.errExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}
}

func TestGetClaim(t *testing.T) {
	data := `{
		"a": 42,
		"b": "bar",
		"c": {
			"d": 95,
			"e": [
				"dog",
				"cat",
				"bird"
			],
			"f": {
				"g": "zebra"
			}
		},
		"h": true,
		"i": false
	}`
	var claims map[string]any
	require.NoError(t, json.Unmarshal([]byte(data), &claims))

	tests := []struct {
		claim string
		value any
	}{
		{"a", float64(42)},
		{"/a", float64(42)},
		{"b", "bar"},
		{"/c/d", float64(95)},
		{"/c/e/1", "cat"},
		{"/c/f/g", "zebra"},
		{"nope", nil},
		{"/c/f/h", nil},
		{"", nil},
		{"\\", nil},
		{"h", true},
		{"i", false},
		{"/c/e", []any{"dog", "cat", "bird"}},
		{"/c/f", map[string]any{"g": "zebra"}},
	}

	for _, test := range tests {
		t.Run(test.claim, func(t *testing.T) {
			v := getClaim(nil, claims, test.claim)
			require.Equal(t, test.value, v)
		})
	}
}

func TestNormalizeList(t *testing.T) {
	tests := []struct {
		raw        any
		normalized []any
		ok         bool
	}{
		{
			raw:        []any{"green", 42},
			normalized: []any{"green", 42},
			ok:         true,
		},
		{
			raw:        []any{"green"},
			normalized: []any{"green"},
			ok:         true,
		},
		{
			raw:        []any{},
			normalized: []any{},
			ok:         true,
		},
		{
			raw:        "green",
			normalized: []any{"green"},
			ok:         true,
		},
		{
			raw:        "",
			normalized: []any{""},
			ok:         true,
		},
		{
			raw:        42,
			normalized: []any{42},
			ok:         true,
		},
		{
			raw:        struct{ A int }{A: 5},
			normalized: nil,
			ok:         false,
		},
		{
			raw:        nil,
			normalized: nil,
			ok:         false,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%#v", tc.raw), func(t *testing.T) {
			normalized, ok := normalizeList(tc.raw)
			assert.Equal(t, tc.normalized, normalized)
			assert.Equal(t, tc.ok, ok)
		})
	}
}

func TestStringifyMetadataValue(t *testing.T) {
	cases := map[string]struct {
		value         any
		expect        string
		expectFailure bool
	}{
		"empty string": {"", "", false},
		"string":       {"foo", "foo", false},
		"true":         {true, "true", false},
		"false":        {false, "false", false},
		"json number":  {json.Number("-12345"), "-12345", false},
		"float64":      {float64(7.9), "7", false},
		//
		"float32": {float32(7.9), "7", false},
		"int8":    {int8(42), "42", false},
		"int16":   {int16(42), "42", false},
		"int32":   {int32(42), "42", false},
		"int64":   {int64(42), "42", false},
		"int":     {int(42), "42", false},
		"uint8":   {uint8(42), "42", false},
		"uint16":  {uint16(42), "42", false},
		"uint32":  {uint32(42), "42", false},
		"uint64":  {uint64(42), "42", false},
		"uint":    {uint(42), "42", false},
		// fail
		"string slice": {[]string{"a"}, "", true},
		"int slice":    {[]int64{99}, "", true},
		"map":          {map[string]int{"a": 99}, "", true},
		"nil":          {nil, "", true},
		"struct":       {struct{ A int }{A: 5}, "", true},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got, ok := stringifyMetadataValue(tc.value)
			if tc.expectFailure {
				require.False(t, ok)
			} else {
				require.True(t, ok)
				require.Equal(t, tc.expect, got)
			}
		})
	}
}

func TestValidateAudience(t *testing.T) {
	tests := []struct {
		boundAudiences    []string
		audience          []string
		errExpectedLax    bool
		errExpectedStrict bool
	}{
		{[]string{"a"}, []string{"a"}, false, false},
		{[]string{"a"}, []string{"b"}, true, true},
		{[]string{"a"}, []string{""}, true, true},
		{[]string{}, []string{"a"}, false, true},
		{[]string{"a", "b"}, []string{"a"}, false, false},
		{[]string{"a", "b"}, []string{"b"}, false, false},
		{[]string{"a", "b"}, []string{"a", "b", "c"}, false, false},
		{[]string{"a", "b"}, []string{"c", "d"}, true, true},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(fmt.Sprintf(
			"boundAudiences=%#v audience=%#v strict=false",
			tc.boundAudiences, tc.audience,
		), func(t *testing.T) {
			err := validateAudience(tc.boundAudiences, tc.audience, false)
			if tc.errExpectedLax {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})

		t.Run(fmt.Sprintf(
			"boundAudiences=%#v audience=%#v strict=true",
			tc.boundAudiences, tc.audience,
		), func(t *testing.T) {
			err := validateAudience(tc.boundAudiences, tc.audience, true)
			if tc.errExpectedStrict {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
