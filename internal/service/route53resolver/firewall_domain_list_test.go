package route53resolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/route53resolver"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfroute53resolver "github.com/hashicorp/terraform-provider-aws/internal/service/route53resolver"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccRoute53ResolverFirewallDomainList_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v route53resolver.FirewallDomainList
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_route53_resolver_firewall_domain_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, route53resolver.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFirewallDomainListDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallDomainListConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallDomainListExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "domains.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccRoute53ResolverFirewallDomainList_domains(t *testing.T) {
	ctx := acctest.Context(t)
	var v route53resolver.FirewallDomainList
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_route53_resolver_firewall_domain_list.test"
	domainName1 := acctest.RandomFQDomainName()
	domainName2 := acctest.RandomFQDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, route53resolver.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFirewallDomainListDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallDomainListConfig_domains(rName, domainName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallDomainListExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "domains.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "domains.*", domainName1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFirewallDomainListConfig_domains(rName, domainName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallDomainListExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "domains.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "domains.*", domainName2),
				),
			},
			{
				Config: testAccFirewallDomainListConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallDomainListExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "domains.#", "0"),
				),
			},
		},
	})
}

func TestAccRoute53ResolverFirewallDomainList_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v route53resolver.FirewallDomainList
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_route53_resolver_firewall_domain_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, route53resolver.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFirewallDomainListDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallDomainListConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallDomainListExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfroute53resolver.ResourceFirewallDomainList(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRoute53ResolverFirewallDomainList_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var v route53resolver.FirewallDomainList
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_route53_resolver_firewall_domain_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, route53resolver.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFirewallDomainListDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallDomainListConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallDomainListExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
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
				Config: testAccFirewallDomainListConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallDomainListExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccFirewallDomainListConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallDomainListExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckFirewallDomainListDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).Route53ResolverConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_route53_resolver_firewall_domain_list" {
				continue
			}

			_, err := tfroute53resolver.FindFirewallDomainListByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Route53 Resolver Firewall Domain List still exists: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckFirewallDomainListExists(ctx context.Context, n string, v *route53resolver.FirewallDomainList) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Route53 Resolver Firewall Domain List ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).Route53ResolverConn(ctx)

		output, err := tfroute53resolver.FindFirewallDomainListByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccFirewallDomainListConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_route53_resolver_firewall_domain_list" "test" {
  name = %[1]q
}
`, rName)
}

func testAccFirewallDomainListConfig_domains(rName, domain string) string {
	return fmt.Sprintf(`
resource "aws_route53_resolver_firewall_domain_list" "test" {
  name    = %[1]q
  domains = [%[2]q]
}
`, rName, domain)
}

func testAccFirewallDomainListConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_route53_resolver_firewall_domain_list" "test" {
  name = %[1]q

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccFirewallDomainListConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_route53_resolver_firewall_domain_list" "test" {
  name = %[1]q

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
