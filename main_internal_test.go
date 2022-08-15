package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelmToHcl(t *testing.T) {
	type test struct {
		value     string
		want      string
		wantError string
	}

	tests := []test{
		{
			value:     "",
			want:      "resource \"helm_release\" \"default\" {\n  values = [\n    yamlencode(\n      {\n      }\n    )\n  ]\n}\n",
			wantError: "",
		},
		{
			value:     "image:\n  {}{}",
			want:      "",
			wantError: "yaml: line 1: did not find expected key",
		},
		{
			value:     "image:\n  registry: docker.io\n  repository: bitnami/postgresql\n  tag: 14.4.0-debian-11-r13",
			want:      "resource \"helm_release\" \"default\" {\n  values = [\n    yamlencode(\n      {\n        image = {\n          registry = \"docker.io\"\n          repository = \"bitnami/postgresql\"\n          tag = \"14.4.0-debian-11-r13\"\n        }\n      }\n    )\n  ]\n}\n",
			wantError: "",
		},
		{
			value:     "key1:\n  - value1\n  - value2\n  - value3",
			want:      "resource \"helm_release\" \"default\" {\n  values = [\n    yamlencode(\n      {\n        key1 = [\n          \"value1\"\n          \"value2\"\n          \"value3\"\n        ]\n      }\n    )\n  ]\n}\n",
			wantError: "",
		},
		{
			value:     "key1: [value1,value2,value3]",
			want:      "resource \"helm_release\" \"default\" {\n  values = [\n    yamlencode(\n      {\n        key1 = [\n          \"value1\"\n          \"value2\"\n          \"value3\"\n        ]\n      }\n    )\n  ]\n}\n",
			wantError: "",
		},
		{
			value:     "one:\n  - id: 1\n    name: franc\n  - id: 11\n    name: Tom",
			want:      "resource \"helm_release\" \"default\" {\n  values = [\n    yamlencode(\n      {\n        one = [\n          {\n            id = 1\n            name = \"franc\"\n          },\n          {\n            id = 11\n            name = \"Tom\"\n          },\n        ]\n      }\n    )\n  ]\n}\n",
			wantError: "",
		},
		{
			value:     "auth:\n  secretKeys:\n    adminPasswordKey: postgres-password",
			want:      "resource \"helm_release\" \"default\" {\n  values = [\n    yamlencode(\n      {\n        auth = {\n          secretKeys = {\n            adminPasswordKey = \"postgres-password\"\n          }\n        }\n      }\n    )\n  ]\n}\n",
			wantError: "",
		},
		{
			value:     "tolerations:\n  - key: \"dedicated\"\n    operator: \"Equal\"\n    value: \"services\"\n    effect: \"NoSchedule\"\n  - key: \"dedicated\"\n    operator: \"Equal\"\n    value: \"service-2\"\n    effect: \"NoSchedule\"",
			want:      "resource \"helm_release\" \"default\" {\n  values = [\n    yamlencode(\n      {\n        tolerations = [\n          {\n            key = \"dedicated\"\n            operator = \"Equal\"\n            value = \"services\"\n            effect = \"NoSchedule\"\n          },\n          {\n            key = \"dedicated\"\n            operator = \"Equal\"\n            value = \"service-2\"\n            effect = \"NoSchedule\"\n          },\n        ]\n      }\n    )\n  ]\n}\n",
			wantError: "",
		},
	}

	for _, tc := range tests {
		fmt.Println(tc.value)
		fmt.Println()

		value, err := helmToHcl([]byte(tc.value))
		fmt.Println(value)
		fmt.Println()

		assert.Equal(t, tc.want, value)
		if err != nil {
			assert.Equal(t, tc.wantError, err.Error())
		}

	}
}

func TestConvertValue(t *testing.T) {
	type test struct {
		value interface{}
		want  string
	}

	tests := []test{
		{
			value: "value",
			want:  `"value"`,
		},
		{
			value: 1,
			want:  "1",
		},
		{
			value: true,
			want:  "true",
		},
		{
			value: nil,
			want:  "",
		},
	}

	for _, tc := range tests {
		value := convertValue(tc.value)
		assert.Equal(t, tc.want, value)
	}
}

func TestStringWithIndent(t *testing.T) {
	type test struct {
		value  string
		indent int
		want   string
	}

	tests := []test{
		{
			value:  "key",
			indent: 0,
			want:   "key",
		},
		{
			value:  "key",
			indent: 2,
			want:   "  key",
		},
	}

	for _, tc := range tests {
		value := stringWithIndent(tc.value, tc.indent)
		assert.Equal(t, tc.want, value)
	}
}

func TestFormatKey(t *testing.T) {
	type test struct {
		key  string
		want string
	}

	tests := []test{
		{
			key:  "value",
			want: "value",
		},
		{
			key:  "warman.io/value",
			want: `"warman.io/value"`,
		},
	}

	for _, tc := range tests {
		key := formatKey(tc.key)
		assert.Equal(t, tc.want, key)
	}
}
