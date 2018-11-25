package cmd

import (
	"context"

	"github.com/blend/go-sdk/logger"
	"github.com/spf13/cobra"
	"github.com/wcharczuk/photoblog/pkg/aws"
	"github.com/wcharczuk/photoblog/pkg/aws/s3"
)

// Deploy returns the deploy command.
func Deploy(configPath *string, log *logger.Logger) *cobra.Command {
	var bucket, region *string
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy the photoblog",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := ReadConfig(*configPath)
			if err != nil {
				log.SyncFatalExit(err)
			}

			if *bucket == "" {
				*bucket = cfg.S3.Bucket
			}
			if *bucket == "" {
				log.SyncFatalf("s3 bucket not set in config or in flags, cannot continue (set at `s3 > bucket` in the config or use --bucket)")
			}

			if *region == "" {
				*region = cfg.S3.Region
			}
			if *region == "" {
				log.SyncFatalf("s3 region not set in config or in flags, cannot continue (set at `s3 > region` in the config or use --region)")
			}

			mgr := s3.New(&aws.Config{
				Region: *region,
			})
			mgr.Log = log
			mgr.PutObjectDefaults = s3.File{
				ACL: s3.ACLPublicRead,
			}
			if err := mgr.SyncDirectory(context.Background(), cfg.OutputPathOrDefault(), *bucket); err != nil {
				log.SyncFatal(err)
			}
		},
	}

	bucket = cmd.Flags().String("bucket", "", "An optional specific bucket (in the form s3://...)")
	region = cmd.Flags().String("region", "", "An optional aws region")
	return cmd
}
