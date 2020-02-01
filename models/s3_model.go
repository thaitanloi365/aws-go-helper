package models

// Credentials Represents AWS credentials and config.
type Credentials struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
}

// PolicyOptions Represents policy options.
type PolicyOptions struct {
	ExpiryMinutes int
	MaxFileSize   int
}

// S3Config config
type S3Config struct {
	Credentials   Credentials
	PolicyOptions PolicyOptions
}

// Signature response
type Signature struct {
	URL        string `json:"url"`
	Policy     string `json:"policy"`
	Credential string `json:"x-amz-credential"`
	Algorithm  string `json:"x-amz-algorithm"`
	Signature  string `json:"x-amz-signature"`
	Date       string `json:"x-amz-date"`
	ACL        string `json:"acl"`
}
