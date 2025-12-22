package storage

import (
	adapter "src/application/adapter/storage"
	"src/core/di"
	"src/core/env"
	impl "src/infrastructure/storage/s3compat"
)

func init() {
	di.RegisterAs[adapter.IStorageAdapter](func() adapter.IStorageAdapter {
		config := &adapter.StorageConfig{
			Endpoint:        env.Get("STORAGE_ENDPOINT", "http://localhost:9000"),
			Region:          env.Get("STORAGE_REGION", "us-east-1"),
			Bucket:          env.Get("STORAGE_BUCKET", "storage"),
			AccessKeyID:     env.Get("STORAGE_ACCESS_KEY_ID", "accessKeyId"),
			SecretAccessKey: env.Get("STORAGE_SECRET_ACCESS_KEY", "secretAccessKey"),
		}
		if err := config.Validate(); err != nil {
			panic(err)
		}
		return impl.NewS3CompatStorageAdapter(config)
	})
}
