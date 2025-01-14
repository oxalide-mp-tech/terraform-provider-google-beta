// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccComputeRoute_routeBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(10),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRoute_routeBasicExample(context),
			},
			{
				ResourceName:      "google_compute_route.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRoute_routeBasicExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_route" "default" {
  name        = "network-route%{random_suffix}"
  dest_range  = "15.0.0.0/24"
  network     = "${google_compute_network.default.name}"
  next_hop_ip = "10.132.1.5"
  priority    = 100
}

resource "google_compute_network" "default" {
  name = "compute-network%{random_suffix}"
}
`, context)
}

func TestAccComputeRoute_routeIlbBetaExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(10),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProvidersOiCS,
		CheckDestroy: testAccCheckComputeRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRoute_routeIlbBetaExample(context),
			},
		},
	})
}

func testAccComputeRoute_routeIlbBetaExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_network" "default" {
  provider                = "google-beta"
  name                    = "compute-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  provider      = "google-beta"
  name          = "compute-subnet%{random_suffix}"
  ip_cidr_range = "10.0.1.0/24"
  region        = "us-central1"
  network       = "${google_compute_network.default.self_link}"
}

resource "google_compute_health_check" "hc" {
  provider           = "google-beta"
  name               = "proxy-health-check%{random_suffix}"
  check_interval_sec = 1
  timeout_sec        = 1

  tcp_health_check {
    port = "80"
  }
}

resource "google_compute_region_backend_service" "backend" {
  provider              = "google-beta"
  name                  = "compute-backend%{random_suffix}"
  region                = "us-central1"
  health_checks         = ["${google_compute_health_check.hc.self_link}"]
}

resource "google_compute_forwarding_rule" "default" {
  provider              = "google-beta"
  name                  = "compute-forwarding-rule%{random_suffix}"
  region                = "us-central1"

  load_balancing_scheme = "INTERNAL"
  backend_service       = "${google_compute_region_backend_service.backend.self_link}"
  all_ports             = true
  network               = "${google_compute_network.default.name}"
  subnetwork            = "${google_compute_subnetwork.default.name}"
}

resource "google_compute_route" "route-ilb-beta" {
  provider     = "google-beta"
  name         = "route-ilb-beta%{random_suffix}"
  dest_range   = "0.0.0.0/0"
  network      = "${google_compute_network.default.name}"
  next_hop_ilb = "${google_compute_forwarding_rule.default.self_link}"
  priority     = 2000
}
`, context)
}

func testAccCheckComputeRouteDestroy(s *terraform.State) error {
	for name, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_route" {
			continue
		}
		if strings.HasPrefix(name, "data.") {
			continue
		}

		config := testAccProvider.Meta().(*Config)

		url, err := replaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/global/routes/{{name}}")
		if err != nil {
			return err
		}

		_, err = sendRequest(config, "GET", "", url, nil)
		if err == nil {
			return fmt.Errorf("ComputeRoute still exists at %s", url)
		}
	}

	return nil
}
