package models

// Configurations configurations
type Configurations struct {
	AppName                     string `mapstructure:"app_name" json:"app_name"`
	Version                     string `mapstructure:"version" json:"version"`
	AWSS3OriginFolder           string `mapstructure:"aws_s3_origin_folder" json:"aws_s3_origin_folder"`
	AWSS3ResizedFolder          string `mapstructure:"aws_s3_resized_folder" json:"aws_s3_resized_folder"`
	AWSS3Bucket                 string `mapstructure:"aws_s3_bucket" json:"aws_s3_bucket"`
	AWSS3ACL                    string `mapstructure:"aws_s3_acl" json:"aws_s3_acl"`
	AWSS3Region                 string `mapstructure:"aws_s3_region" json:"aws_s3_region"`
	AWSAccessKey                string `mapstructure:"aws_access_key" json:"aws_access_key"`
	AWSSecretKey                string `mapstructure:"aws_secret_key" json:"aws_secret_key"`
	AWSMaxFileSize              uint   `mapstructure:"aws_max_file_size" json:"aws_max_file_size"`
	AWSSignatureExpiryInMinutes int    `mapstructure:"aws_signature_expiry_in_minutes" json:"aws_signature_expiry_in_minutes"`
}
