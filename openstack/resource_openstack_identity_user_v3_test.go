package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIdentityV3User_basic(t *testing.T) {
	var project projects.Project

	projectName := "ACCPTTEST-" + acctest.RandString(5)

	var user users.User

	userName := "ACCPTTEST-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3UserDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3UserBasic(projectName, userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserExists(t.Context(), "openstack_identity_user_v3.user_1", &user),
					testAccCheckIdentityV3ProjectExists(t.Context(), "openstack_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_user_v3.user_1", "name", &user.Name),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_user_v3.user_1", "description", &user.Description),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "ignore_change_password_upon_first_use", "true"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_enabled", "true"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.#", "2"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.0.rule.0", "password"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.0.rule.1", "totp"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.1.rule.0", "password"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.1.rule.1", "custom-auth-method"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "extra.email", "jdoe@example.com"),
				),
			},
			{
				Config: testAccIdentityV3UserUpdate(projectName, userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserExists(t.Context(), "openstack_identity_user_v3.user_1", &user),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_user_v3.user_1", "name", &user.Name),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_user_v3.user_1", "description", &user.Description),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "enabled", "false"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "ignore_change_password_upon_first_use", "false"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.0.rule.0", "password"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.0.rule.1", "totp"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "extra.email", "jdoe@foobar.com"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3UserDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		identityClient, err := config.IdentityV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_identity_user_v3" {
				continue
			}

			_, err := users.Get(ctx, identityClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("User still exists")
			}
		}

		return nil
	}
}

func testAccCheckIdentityV3UserExists(ctx context.Context, n string, user *users.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		identityClient, err := config.IdentityV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %w", err)
		}

		found, err := users.Get(ctx, identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("User not found")
		}

		*user = *found

		return nil
	}
}

func testAccIdentityV3UserBasic(projectName, userName string) string {
	return fmt.Sprintf(`
    resource "openstack_identity_project_v3" "project_1" {
      name = "%s"
    }

    resource "openstack_identity_user_v3" "user_1" {
      default_project_id = "${openstack_identity_project_v3.project_1.id}"
      name = "%s"
      description = "A user"
      password = "password123"
      ignore_change_password_upon_first_use = true
      multi_factor_auth_enabled = true

      multi_factor_auth_rule {
        rule = ["password", "totp"]
      }

      multi_factor_auth_rule {
        rule = ["password", "custom-auth-method"]
      }

      extra = {
        email = "jdoe@example.com"
      }
    }
  `, projectName, userName)
}

func testAccIdentityV3UserUpdate(projectName, userName string) string {
	return fmt.Sprintf(`
    resource "openstack_identity_project_v3" "project_1" {
      name = "%s"
    }

    resource "openstack_identity_user_v3" "user_1" {
      default_project_id = "${openstack_identity_project_v3.project_1.id}"
      name = "%s"
      description = "Some user"
      enabled = false
      password = "password123"
      ignore_change_password_upon_first_use = false
      multi_factor_auth_enabled = true

      multi_factor_auth_rule {
        rule = ["password", "totp"]
      }

      extra = {
        email = "jdoe@foobar.com"
      }
    }
  `, projectName, userName)
}
