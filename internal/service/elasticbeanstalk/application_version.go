package elasticbeanstalk

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_elastic_beanstalk_application_version", name="Application Version")
// @Tags(identifierAttribute="arn")
func ResourceApplicationVersion() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceApplicationVersionCreate,
		ReadWithoutTimeout:   resourceApplicationVersionRead,
		UpdateWithoutTimeout: resourceApplicationVersionUpdate,
		DeleteWithoutTimeout: resourceApplicationVersionDelete,

		CustomizeDiff: verify.SetTagsDiff,

		Schema: map[string]*schema.Schema{
			"application": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},
	}
}

func resourceApplicationVersionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ElasticBeanstalkConn(ctx)

	application := d.Get("application").(string)
	description := d.Get("description").(string)
	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	name := d.Get("name").(string)

	s3Location := elasticbeanstalk.S3Location{
		S3Bucket: aws.String(bucket),
		S3Key:    aws.String(key),
	}

	createOpts := elasticbeanstalk.CreateApplicationVersionInput{
		ApplicationName: aws.String(application),
		Description:     aws.String(description),
		SourceBundle:    &s3Location,
		Tags:            getTagsIn(ctx),
		VersionLabel:    aws.String(name),
	}

	_, err := conn.CreateApplicationVersionWithContext(ctx, &createOpts)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Elastic Beanstalk Application Version (%s): %s", name, err)
	}

	d.SetId(name)

	return append(diags, resourceApplicationVersionRead(ctx, d, meta)...)
}

func resourceApplicationVersionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ElasticBeanstalkConn(ctx)

	resp, err := conn.DescribeApplicationVersionsWithContext(ctx, &elasticbeanstalk.DescribeApplicationVersionsInput{
		ApplicationName: aws.String(d.Get("application").(string)),
		VersionLabels:   []*string{aws.String(d.Id())},
	})
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Elastic Beanstalk Application Version (%s): %s", d.Id(), err)
	}

	if len(resp.ApplicationVersions) == 0 {
		log.Printf("[DEBUG] Elastic Beanstalk application version read: application version not found")

		d.SetId("")

		return diags
	} else if len(resp.ApplicationVersions) != 1 {
		return sdkdiag.AppendErrorf(diags, "reading application version properties: found %d versions of label %q, expected 1",
			len(resp.ApplicationVersions), d.Id())
	}

	arn := aws.StringValue(resp.ApplicationVersions[0].ApplicationVersionArn)
	d.Set("arn", arn)
	d.Set("description", resp.ApplicationVersions[0].Description)

	return diags
}

func resourceApplicationVersionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ElasticBeanstalkConn(ctx)

	if d.HasChange("description") {
		if err := resourceApplicationVersionDescriptionUpdate(ctx, conn, d); err != nil {
			return sdkdiag.AppendErrorf(diags, "updating Elastic Beanstalk Application Version (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceApplicationVersionRead(ctx, d, meta)...)
}

func resourceApplicationVersionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ElasticBeanstalkConn(ctx)

	application := d.Get("application").(string)
	name := d.Id()

	if !d.Get("force_delete").(bool) {
		environments, err := versionUsedBy(ctx, application, name, conn)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "deleting Elastic Beanstalk Application Version (%s): %s", d.Id(), err)
		}

		if len(environments) > 1 {
			return sdkdiag.AppendErrorf(diags, "Unable to delete Application Version, it is currently in use by the following environments: %s.", environments)
		}
	}
	_, err := conn.DeleteApplicationVersionWithContext(ctx, &elasticbeanstalk.DeleteApplicationVersionInput{
		ApplicationName:    aws.String(application),
		VersionLabel:       aws.String(name),
		DeleteSourceBundle: aws.Bool(false),
	})

	// application version is pending delete, or no longer exists.
	if tfawserr.ErrCodeEquals(err, "InvalidParameterValue") {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Elastic Beanstalk Application version (%s): %s", d.Id(), err)
	}

	return diags
}

func resourceApplicationVersionDescriptionUpdate(ctx context.Context, conn *elasticbeanstalk.ElasticBeanstalk, d *schema.ResourceData) error {
	application := d.Get("application").(string)
	description := d.Get("description").(string)
	name := d.Get("name").(string)

	_, err := conn.UpdateApplicationVersionWithContext(ctx, &elasticbeanstalk.UpdateApplicationVersionInput{
		ApplicationName: aws.String(application),
		Description:     aws.String(description),
		VersionLabel:    aws.String(name),
	})

	return err
}

func versionUsedBy(ctx context.Context, applicationName, versionLabel string, conn *elasticbeanstalk.ElasticBeanstalk) ([]string, error) {
	now := time.Now()
	resp, err := conn.DescribeEnvironmentsWithContext(ctx, &elasticbeanstalk.DescribeEnvironmentsInput{
		ApplicationName:       aws.String(applicationName),
		VersionLabel:          aws.String(versionLabel),
		IncludeDeleted:        aws.Bool(true),
		IncludedDeletedBackTo: aws.Time(now.Add(-1 * time.Minute)),
	})

	if err != nil {
		return nil, err
	}

	var environmentIDs []string
	for _, environment := range resp.Environments {
		environmentIDs = append(environmentIDs, *environment.EnvironmentId)
	}

	return environmentIDs, nil
}
