package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBlockStorageVolumeTypeV3_basic(t *testing.T) {
	var volumetype volumetypes.VolumeType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageVolumeTypeV3Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageVolumeTypeV3Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageVolumeTypeV3Exists(t.Context(), "openstack_blockstorage_volume_type_v3.volume_type_1", &volumetype),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "name", "foo"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "description", "foo"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "is_public", "true"),
				),
			},
			{
				Config: testAccBlockStorageVolumeTypeV3Update1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageVolumeTypeV3Exists(t.Context(), "openstack_blockstorage_volume_type_v3.volume_type_1", &volumetype),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "name", "bar-baz"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "description", "bar-baz"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "is_public", "false"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "extra_specs.%", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "extra_specs.bar", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "extra_specs.baz", "baz"),
				),
			},
			{
				Config: testAccBlockStorageVolumeTypeV3Update2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageVolumeTypeV3Exists(t.Context(), "openstack_blockstorage_volume_type_v3.volume_type_1", &volumetype),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "name", "foo-foo"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "description", "bar-bar"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "is_public", "false"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "extra_specs.%", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "extra_specs.bar", "baz"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "extra_specs.foo", "foo"),
				),
			},
		},
	})
}

func testAccCheckBlockStorageVolumeTypeV3Destroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		blockStorageClient, err := config.BlockStorageV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack block storage client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_blockstorage_volume_type_v3" {
				continue
			}

			_, err := volumetypes.Get(ctx, blockStorageClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("VolumeType still exists")
			}
		}

		return nil
	}
}

func testAccCheckBlockStorageVolumeTypeV3Exists(ctx context.Context, n string, volumetype *volumetypes.VolumeType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		blockStorageClient, err := config.BlockStorageV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack block storage client: %w", err)
		}

		found, err := volumetypes.Get(ctx, blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("VolumeType not found")
		}

		*volumetype = *found

		return nil
	}
}

const testAccBlockStorageVolumeTypeV3Basic = `
resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
	name = "foo"
	description = "foo"
	is_public = true

}
`

const testAccBlockStorageVolumeTypeV3Update1 = `
resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
	name = "bar-baz"
	description = "bar-baz"
	is_public = false
	extra_specs = {
	  bar = "bar"
	  baz = "baz"
	}

}
`

const testAccBlockStorageVolumeTypeV3Update2 = `
resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
	name = "foo-foo"
	description = "bar-bar"
	is_public = false
	extra_specs = {
      bar = "baz"
	  foo = "foo"
	}
}
`
