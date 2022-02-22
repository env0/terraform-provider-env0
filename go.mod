module github.com/env0/terraform-provider-env0

go 1.16 // please change also in `ci.yml`,`update-generated-docs.yml` and `release.yml`

require (
	github.com/adhocore/gronx v0.2.6
	github.com/go-resty/resty/v2 v2.6.0
	github.com/golang/mock v1.4.3
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/terraform-exec v0.15.0 // indirect
	github.com/hashicorp/terraform-plugin-docs v0.5.1 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.4.4
	github.com/jarcoal/httpmock v1.0.8
	github.com/jinzhu/copier v0.3.2
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/onsi/ginkgo v1.16.2
	github.com/onsi/gomega v1.12.0
	github.com/zclconf/go-cty v1.10.0 // indirect
)
