package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/blend/go-sdk/ansi/slant"
	"github.com/spf13/cobra"

	"github.com/wcharczuk/blogctl/pkg/aws"
	"github.com/wcharczuk/blogctl/pkg/aws/cloudfront"
	"github.com/wcharczuk/blogctl/pkg/aws/s3"
	"github.com/wcharczuk/blogctl/pkg/config"
)

// Deploy returns the deploy command.
func Deploy(flags config.Flags) *cobra.Command {
	var bucket, region *string
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy the photoblog",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, cfgPaths, err := config.ReadConfig(flags)
			Fatal(err)

			log := Logger(flags, "deploy")
			slant.Print(log.Output, "BLOGCTL")

			if len(cfgPaths) > 0 {
				log.Infof("using config path(s): %s", strings.Join(cfgPaths, ", "))
			}
			log.Infof("using parallelism: %d", *flags.Parallelism)

			if *bucket == "" {
				*bucket = cfg.S3.Bucket
			}
			if *bucket == "" {
				Fatal(fmt.Errorf("s3 bucket not set in config or in flags, cannot continue (set at `s3 > bucket` in the config or use --bucket)"))
			}

			if *region == "" {
				*region = cfg.S3.Region
			}
			if *region == "" {
				Fatal(fmt.Errorf("s3 region not set in config or in flags, cannot continue (set at `s3 > region` in the config or use --region)"))
			}

			mgr := s3.New(aws.Config{
				Region: *region,
			})
			mgr.Log = log
			mgr.Parallelism = *flags.Parallelism
			mgr.DryRun = *flags.DryRun
			mgr.PutObjectDefaults = s3.File{
				ACL: s3.ACLPublicRead,
			}
			paths, err := mgr.SyncDirectory(context.Background(), cfg.OutputPathOrDefault(), *bucket)
			Fatal(err)

			if !cfg.Cloudfront.IsZero() && len(paths) > 0 {
				if !mgr.DryRun {
					log.Infof("cloudfront invalidating %d paths", len(paths))
					if err := cloudfront.InvalidateMany(context.Background(), mgr.Session, cfg.Cloudfront.Distribution, paths...); err != nil {
						Fatal(err)
					}
				} else {
					log.Debugf("(dry run) cloudfront invalidating %d files", len(paths))
				}
			}
		},
	}

	bucket = cmd.Flags().String("bucket", "", "An optional specific bucket (in the form s3://...)")
	region = cmd.Flags().String("region", "", "An optional aws region")
	return cmd
}
