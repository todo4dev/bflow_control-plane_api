package storage

import "src/core/validator"

type StorageConfig struct {
	Endpoint        string
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
}

var _ validator.IValidable = (*StorageConfig)(nil)

func (config *StorageConfig) Validate() error {
	return validator.Object(config,
		validator.String(&config.Endpoint).Trim().Required().URI().Default("http://localhost:9000"),
		validator.String(&config.Region).Trim().Required().Default("us-east-1"),
		validator.String(&config.Bucket).Trim().Required().Default("storage"),
		validator.String(&config.AccessKeyID).Trim().Required().Default("accessKeyId"),
		validator.String(&config.SecretAccessKey).Trim().Required().Default("secretAccessKey"),
	).Validate()
}
