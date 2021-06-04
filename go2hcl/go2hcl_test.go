package go2hcl

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"testing"
)

func TestHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTP Client Tests")
}

var _ = Describe("Help Functions", func() {
	Describe("ResourceConfigCreate", func() {
		It("should create hcl resource correctly", func() {
			hcl := ResourceConfigCreate("env0_x", "myresource", map[string]interface{}{
				"aaa": "test",
				"bbb": []int{0, 1, 2},
				"ccc": map[string]interface{}{
					"field1": 123,
					"field2": true,
					"field3": "hello",
				},
			})

			expected :=
				`resource "env0_x" "myresource" {
	aaa = "test"
	bbb = [0, 1, 2]
	ccc {
		field1 = 123
		field2 = true
		field3 = "hello"
	}
}`
			Expect(hcl).To(Equal(expected))
		})
	})

	Describe("DataSourceConfigCreate", func() {
		It("should create hcl resource correctly", func() {
			hcl := DataSourceConfigCreate("env0_x", "mydata", map[string]interface{}{
				"aaa": "test",
				"bbb": []int{0, 1, 2},
				"ccc": map[string]interface{}{
					"field1": 123,
					"field2": true,
					"field3": "hello",
				},
			})

			expected :=
				`data "env0_x" "mydata" {
	aaa = "test"
	bbb = [0, 1, 2]
	ccc {
		field1 = 123
		field2 = true
		field3 = "hello"
	}
}`
			Expect(hcl).To(Equal(expected))
		})
	})

	DescribeTable("toHclValue",
		func(value interface{}, expected types.GomegaMatcher, expectedError types.GomegaMatcher) {
			result, err := toHclValue(value, 0)
			Expect(result).To(expected)
			Expect(err).To(expectedError)
		},
		Entry("int", 123, Equal("123"), BeNil()),
		Entry("float", 123.456, Equal("123.456"), BeNil()),
		Entry("float", 123.456789012345678, Equal("123.45678901234568"), BeNil()),
		Entry("string", "hello", Equal("\"hello\""), BeNil()),
		Entry("boolean true", true, Equal("true"), BeNil()),
		Entry("boolean false", false, Equal("false"), BeNil()),
		Entry("array of strings", []string{"hello", "world"}, Equal("[\"hello\", \"world\"]"), BeNil()),
		Entry("array of int", []int{123, 456}, Equal("[123, 456]"), BeNil()),
		Entry("array of bool", []bool{true, false}, Equal("[true, false]"), BeNil()),
		Entry("map fields", map[string]interface{}{
			"field1": 123,
			"field2": true,
			"field3": "hello",
		}, Equal("{\n\tfield1 = 123\n\tfield2 = true\n\tfield3 = \"hello\"\n}"), BeNil()),
		Entry("map of map fields", map[string]map[string]interface{}{
			"statement": {
				"field1": 123,
				"field2": true,
				"field3": "hello",
			},
		}, Equal("{\n\tstatement {\n\t\tfield1 = 123\n\t\tfield2 = true\n\t\tfield3 = \"hello\"\n\t}\n}"), BeNil()),
	)
})
