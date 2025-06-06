// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dms_test

import (
	"context"
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfdms "github.com/hashicorp/terraform-provider-aws/internal/service/dms"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccDMSReplicationSubnetGroup_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_dms_replication_subnet_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DMSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckReplicationSubnetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccReplicationSubnetGroupConfig_basic(rName, "desc1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckReplicationSubnetGroupExists(ctx, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_subnet_group_arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_subnet_group_description", "desc1"),
					resource.TestCheckResourceAttr(resourceName, "replication_subnet_group_id", rName),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "3"),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, "0"),
					resource.TestCheckResourceAttrSet(resourceName, names.AttrVPCID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccReplicationSubnetGroupConfig_basic(rName, "desc2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckReplicationSubnetGroupExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "replication_subnet_group_description", "desc2"),
				),
			},
		},
	})
}

func TestAccDMSReplicationSubnetGroup_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_dms_replication_subnet_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DMSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckReplicationSubnetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccReplicationSubnetGroupConfig_basic(rName, "desc1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckReplicationSubnetGroupExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdms.ResourceReplicationSubnetGroup(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckReplicationSubnetGroupExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DMSClient(ctx)

		_, err := tfdms.FindReplicationSubnetGroupByID(ctx, conn, rs.Primary.ID)

		return err
	}
}

func testAccCheckReplicationSubnetGroupDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DMSClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_dms_replication_subnet_group" {
				continue
			}

			_, err := tfdms.FindReplicationSubnetGroupByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("DMS Replication Subnet Group %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccReplicationSubnetGroupConfig_basic(rName, description string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 3), fmt.Sprintf(`
resource "aws_dms_replication_subnet_group" "test" {
  replication_subnet_group_id          = %[1]q
  replication_subnet_group_description = %[2]q
  subnet_ids                           = aws_subnet.test[*].id
}
`, rName, description))
}
