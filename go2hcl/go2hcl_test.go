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
	DescribeTable("toHclValue",
		func(value interface{}, expected types.GomegaMatcher, expectedError types.GomegaMatcher) {
			result, err := toHclValue(value)
			Expect(result).To(expected)
			Expect(err).To(expectedError)
		},
		Entry("int", 123, Equal("123"), BeNil()),
		Entry("float", 123.456, Equal("123.456"), BeNil()),
		Entry("float", 123.456789012345678, Equal("123.45678901234568"), BeNil()),
		Entry("string", "hello", Equal("\"hello\""), BeNil()),
		Entry("boolean true", true, Equal("true"), BeNil()),
		Entry("boolean false", false, Equal("false"), BeNil()),
		Entry("array of strings", []string{"hello", "world"}, Equal("[\n\t\"hello\",\n\t\"world\"\n]"), BeNil()),
		Entry("array of int", []int{123, 456}, Equal("[\n\t123,\n\t456\n]"), BeNil()),
		Entry("array of bool", []bool{true, false}, Equal("[\n\ttrue,\n\tfalse\n]"), BeNil()),
		Entry("map fields", map[string]interface{}{
			"field1": 123,
			"field2": true,
			"field3": "hello",
		}, Equal("{\n\tfield1 = 123\n\tfield2 = true\n\tfield3 = \"hello\"\n}"), BeNil()),
	)
})
