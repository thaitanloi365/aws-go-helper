package models

// Configurations configurations
type Configurations struct {
	AppName                     string `mapstructure:"app_name" json:"app_name"`
	Version                     string `mapstructure:"version" json:"version"`
	AWSS3ResizedBucket          string `mapstructure:"aws_s3_resized_bucket" json:"aws_s3_resized_bucket"`
	AWSS3Bucket                 string `mapstructure:"aws_s3_bucket" json:"aws_s3_bucket"`
	AWSS3ACL                    string `mapstructure:"aws_s3_acl" json:"aws_s3_acl"`
	AWSS3Region                 string `mapstructure:"aws_s3_region" json:"aws_s3_region"`
	AWSMaxFileSize              uint   `mapstructure:"aws_max_file_size" json:"aws_max_file_size"`
	AWSSignatureExpiryInMinutes int    `mapstructure:"aws_signature_expiry_in_minutes" json:"aws_signature_expiry_in_minutes"`
}
