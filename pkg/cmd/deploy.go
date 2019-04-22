package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/blend/go-sdk/logger"
	"github.com/spf13/cobra"

	"github.com/wcharczuk/blogctl/pkg/aws"
	"github.com/wcharczuk/blogctl/pkg/aws/cloudfront"
	"github.com/wcharczuk/blogctl/pkg/aws/s3"
	"github.com/wcharczuk/blogctl/pkg/config"
	"github.com/wcharczuk/blogctl/pkg/engine"
)

// Deploy returns the deploy command.
func Deploy(flags *config.PersistentFlags) *cobra.Command {
	var bucket, region *string
	var dryRun *bool
	cmd := &cobra.Command{
		Use:     "deploy",
		Aliases: []string{"d", "deploy"},
		Short:   "Deploy the photoblog",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, cfgPath, err := engine.ReadConfig(flags)
			if err != nil {
				logger.FatalExit(err)
			}
			log := Logger(cfg, "deploy")
			fmt.Fprintf(log.Logger.Output, banner)
			if cfgPath != "" {
				log.Infof("using config path: %s", cfgPath)
			}

			if *bucket == "" {
				*bucket = cfg.S3.Bucket
			}
			if *bucket == "" {
				log.Fatalf("s3 bucket not set in config or in flags, cannot continue (set at `s3 > bucket` in the config or use --bucket)")
				os.Exit(1)
			}

			if *region == "" {
				*region = cfg.S3.Region
			}
			if *region == "" {
				log.Fatalf("s3 region not set in config or in flags, cannot continue (set at `s3 > region` in the config or use --region)")
				os.Exit(1)
			}

			mgr := s3.New(&aws.Config{
				Region: *region,
			})
			mgr.DryRun = *dryRun
			mgr.Log = log
			mgr.PutObjectDefaults = s3.File{
				ACL: s3.ACLPublicRead,
			}
			paths, err := mgr.SyncDirectory(context.Background(), cfg.OutputPathOrDefault(), *bucket)

			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			if !mgr.DryRun {
				if !cfg.Cloudfront.IsZero() && len(paths) > 0 {
					log.Infof("cloudfront invalidating %d paths", len(paths))
					if err := cloudfront.InvalidateMany(context.Background(), mgr.Session, cfg.Cloudfront.Distribution, paths...); err != nil {
						log.Fatal(err)
						os.Exit(1)
					}
				}
			} else {
				log.Debugf("dry run; would invalidate %d files", len(paths))
			}
			log.Info("complete")
		},
	}

	dryRun = cmd.Flags().Bool("dry-run", false, "If we should only print the plan, and not realize changes")
	bucket = cmd.Flags().String("bucket", "", "An optional specific bucket (in the form s3://...)")
	region = cmd.Flags().String("region", "", "An optional aws region")
	return cmd
}
