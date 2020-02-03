package config

import (
	"aws-go-helper/models"
	"os"
	"strconv"
)

// Instance Configuration instance
var Instance models.Configurations

// SetupEnv Setup Env
func SetupEnv(filenames ...string) {

	var cfg = &models.Configurations{
		AppName:            os.Getenv("app_name"),
		Version:            os.Getenv("version"),
		AWSS3Bucket:        os.Getenv("aws_s3_bucket"),
		AWSS3ResizedBucket: os.Getenv("aws_s3_resized_bucket"),
		AWSS3Region:        os.Getenv("aws_s3_region"),
		AWSS3ACL:           os.Getenv("aws_s3_acl"),
	}

	if value, err := strconv.ParseUint(os.Getenv("aws_max_file_size"), 10, 64); err == nil {
		if value > 0 {
			cfg.AWSMaxFileSize = uint(value)

		}
	}

	if value, err := strconv.ParseInt(os.Getenv("aws_signature_expiry_in_minutes"), 10, 64); err == nil {
		if value > 0 {
			cfg.AWSSignatureExpiryInMinutes = int(value)
		}
	}

	Instance = *cfg

}
