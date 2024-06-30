// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDnsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsResourceConfigEmpty(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nanoid_dns.test", "length", "10"),
					resource.TestCheckResourceAttrWith("nanoid_dns.test", "id", testCheckLen(10)),
				),
			},
			{
				ResourceName:      "nanoid_dns.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDnsResource_WithLength(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsResourceConfig(9),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nanoid_dns.test", "length", "9"),
					resource.TestCheckResourceAttrWith("nanoid_dns.test", "id", testCheckLen(9)),
				),
			},
			{
				ResourceName:      "nanoid_dns.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDnsResourceConfig(length int) string {
	lengthStr := fmt.Sprintf("length = %d", length)
	return fmt.Sprintf(`
resource "nanoid_dns" "test" {
  %s
}
`, lengthStr)
}

func testAccDnsResourceConfigEmpty() string {
	return `resource "nanoid_dns" "test" {}`
}
