package cmd

import (
	"context"
	"os"

	"github.com/blend/go-sdk/logger"
	"github.com/spf13/cobra"

	"github.com/wcharczuk/blogctl/pkg/aws"
	"github.com/wcharczuk/blogctl/pkg/aws/cloudfront"
	"github.com/wcharczuk/blogctl/pkg/aws/s3"
	"github.com/wcharczuk/blogctl/pkg/engine"
)

// Deploy returns the deploy command.
func Deploy(configPath *string, log *logger.Logger) *cobra.Command {
	var bucket, region *string
	var dryRun *bool
	cmd := &cobra.Command{
		Use:     "deploy",
		Aliases: []string{"d", "deploy"},
		Short:   "Deploy the photoblog",
		Run: func(cmd *cobra.Command, args []string) {
			deployLog := log.SubContext("deploy")
			cfg, err := engine.ReadConfig(*configPath)
			if err != nil {
				deployLog.SyncFatal(err)
				os.Exit(1)
			}

			if *bucket == "" {
				*bucket = cfg.S3.Bucket
			}
			if *bucket == "" {
				deployLog.SyncFatalf("s3 bucket not set in config or in flags, cannot continue (set at `s3 > bucket` in the config or use --bucket)")
				os.Exit(1)
			}

			if *region == "" {
				*region = cfg.S3.Region
			}
			if *region == "" {
				deployLog.SyncFatalf("s3 region not set in config or in flags, cannot continue (set at `s3 > region` in the config or use --region)")
				os.Exit(1)
			}

			mgr := s3.New(&aws.Config{
				Region: *region,
			})
			mgr.DryRun = *dryRun
			mgr.Log = deployLog
			mgr.PutObjectDefaults = s3.File{
				ACL: s3.ACLPublicRead,
			}
			paths, err := mgr.SyncDirectory(context.Background(), cfg.OutputPathOrDefault(), *bucket)

			if err != nil {
				deployLog.SyncFatal(err)
			}

			if !mgr.DryRun {
				if !cfg.Cloudfront.IsZero() && len(paths) > 0 {
					deployLog.SyncInfof("cloudfront invalidating %d paths", len(paths))
					if err := cloudfront.InvalidateMany(context.Background(), mgr.Session, cfg.Cloudfront.Distribution, paths...); err != nil {
						deployLog.SyncFatal(err)
						os.Exit(1)
					}
				}
			} else {
				deployLog.SyncDebugf("dry run; would invalidate %d files", len(paths))
			}
		},
	}

	dryRun = cmd.Flags().Bool("dry-run", false, "If we should only print the plan, and not realize changes")
	bucket = cmd.Flags().String("bucket", "", "An optional specific bucket (in the form s3://...)")
	region = cmd.Flags().String("region", "", "An optional aws region")
	return cmd
}
