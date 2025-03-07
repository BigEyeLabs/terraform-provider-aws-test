package lightsail_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tflightsail "github.com/hashicorp/terraform-provider-aws/internal/service/lightsail"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// LightSail only allows 5 load balancers per account. Serializing these tests simplifies running all
// LoadBalancer tests without risk of hitting the account limit.
func TestAccLightsailLoadBalancer_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]map[string]func(t *testing.T){
		"lb": {
			"basic":             testAccLoadBalancer_basic,
			"disappears":        testAccLoadBalancer_disappears,
			"name":              testAccLoadBalancer_name,
			"health_check_path": testAccLoadBalancer_healthCheckPath,
			"tags":              testAccLoadBalancer_tags,
		},
		"lb_attachment": {
			"basic":      testAccLoadBalancerAttachment_basic,
			"disappears": testAccLoadBalancerAttachment_disappears,
		},
		"lb_certificate": {
			"basic":                     testAccLoadBalancerCertificate_basic,
			"disappears":                testAccLoadBalancerCertificate_disappears,
			"domain_validation_records": testAccLoadBalancerCertificate_domainValidationRecords,
			"subject_alternative_names": testAccLoadBalancerCertificate_subjectAlternativeNames,
		},
		"lb_certificate_attachment": {
			"basic": testAccLoadBalancerCertificateAttachment_basic,
		},
		"lb_https_redirection_policy": {
			"basic": testAccLoadBalancerHTTPSRedirectionPolicy_basic,
		},
		"lb_stickiness_policy": {
			"basic":           testAccLoadBalancerStickinessPolicy_basic,
			"cookie_duration": testAccLoadBalancerStickinessPolicy_cookieDuration,
			"enabled":         testAccLoadBalancerStickinessPolicy_enabled,
			"disappears":      testAccLoadBalancerStickinessPolicy_disappears,
		},
	}

	acctest.RunSerialTests2Levels(t, testCases, 0)
}
func testAccLoadBalancer_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_lightsail_lb.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, strings.ToLower(lightsail.ServiceID))
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, strings.ToLower(lightsail.ServiceID)),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "health_check_path", "/"),
					resource.TestCheckResourceAttr(resourceName, "instance_port", "80"),
					resource.TestCheckResourceAttrSet(resourceName, "dns_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLoadBalancer_name(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	lightsailNameWithSpaces := fmt.Sprint(rName, "string with spaces")
	lightsailNameWithStartingDigit := fmt.Sprintf("01-%s", rName)
	lightsailNameWithUnderscore := fmt.Sprintf("%s_123456", rName)
	resourceName := "aws_lightsail_lb.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, strings.ToLower(lightsail.ServiceID))
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, strings.ToLower(lightsail.ServiceID)),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config:      testAccLoadBalancerConfig_basic(lightsailNameWithSpaces),
				ExpectError: regexp.MustCompile(`must contain only alphanumeric characters, underscores, hyphens, and dots`),
			},
			{
				Config:      testAccLoadBalancerConfig_basic(lightsailNameWithStartingDigit),
				ExpectError: regexp.MustCompile(`must begin with an alphabetic character`),
			},
			{
				Config: testAccLoadBalancerConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLoadBalancerExists(ctx, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "health_check_path"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_port"),
				),
			},
			{
				Config: testAccLoadBalancerConfig_basic(lightsailNameWithUnderscore),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLoadBalancerExists(ctx, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "health_check_path"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_port"),
				),
			},
		},
	})
}

func testAccLoadBalancer_healthCheckPath(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_lightsail_lb.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, strings.ToLower(lightsail.ServiceID))
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, strings.ToLower(lightsail.ServiceID)),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_healthCheckPath(rName, "/"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "health_check_path", "/"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoadBalancerConfig_healthCheckPath(rName, "/healthcheck"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "health_check_path", "/healthcheck"),
				),
			},
		},
	})
}

func testAccLoadBalancer_tags(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_lightsail_lb.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, strings.ToLower(lightsail.ServiceID))
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, strings.ToLower(lightsail.ServiceID)),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoadBalancerConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccLoadBalancerConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckLoadBalancerExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No LightsailLoadBalancer ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).LightsailClient(ctx)

		resp, err := tflightsail.FindLoadBalancerById(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		if resp == nil {
			return fmt.Errorf("Load Balancer %q does not exist", rs.Primary.ID)
		}

		return nil
	}
}

func testAccLoadBalancer_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_lightsail_lb.test"

	testDestroy := func(*terraform.State) error {
		// reach out and DELETE the LoadBalancer
		conn := acctest.Provider.Meta().(*conns.AWSClient).LightsailClient(ctx)
		_, err := conn.DeleteLoadBalancer(ctx, &lightsail.DeleteLoadBalancerInput{
			LoadBalancerName: aws.String(rName),
		})

		if err != nil {
			return fmt.Errorf("error deleting Lightsail LoadBalancer in disappear test")
		}

		// sleep 7 seconds to give it time, so we don't have to poll
		time.Sleep(7 * time.Second)

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, strings.ToLower(lightsail.ServiceID))
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, strings.ToLower(lightsail.ServiceID)),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(ctx, resourceName),
					testDestroy,
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckLoadBalancerDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_lightsail_lb" {
				continue
			}

			conn := acctest.Provider.Meta().(*conns.AWSClient).LightsailClient(ctx)

			_, err := tflightsail.FindLoadBalancerById(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return create.Error(names.Lightsail, create.ErrActionCheckingDestroyed, tflightsail.ResLoadBalancer, rs.Primary.ID, errors.New("still exists"))
		}

		return nil
	}
}

func testAccLoadBalancerConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_lightsail_lb" "test" {
  name              = %[1]q
  health_check_path = "/"
  instance_port     = "80"
}
`, rName)
}

func testAccLoadBalancerConfig_healthCheckPath(rName string, rPath string) string {
	return fmt.Sprintf(`
resource "aws_lightsail_lb" "test" {
  name              = %[1]q
  health_check_path = %[2]q
  instance_port     = "80"
}
`, rName, rPath)
}

func testAccLoadBalancerConfig_tags1(rName string, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_lightsail_lb" "test" {
  name              = %[1]q
  health_check_path = "/"
  instance_port     = "80"
  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccLoadBalancerConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_lightsail_lb" "test" {
  name              = %[1]q
  health_check_path = "/"
  instance_port     = "80"
  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
