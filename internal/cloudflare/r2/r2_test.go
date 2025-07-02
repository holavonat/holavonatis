package r2_test

import (
	"os"
	"testing"

	"github.com/holavonat/holavonatis/internal/cloudflare/r2"
	"github.com/holavonat/holavonatis/internal/config"
	r "github.com/stretchr/testify/require"
)

func TestUploadFile(t *testing.T) {
	cfg, err := config.GetConfig()
	r.NoError(t, err)
	t.Log("Access key:", cfg.ObjectStorage.AccessKeyID)
	t.Log("Secret access key:", cfg.ObjectStorage.SecretAccessKey)
	t.Log("Bucket name:", cfg.ObjectStorage.BucketName)
	t.Log("Object path:", cfg.ObjectStorage.ObjectPath)
	t.Log("Endpoint URL:", cfg.ObjectStorage.EndpointURL)

	cloudflare, err := r2.NewClient(r2.Cloudflare{
		BucketName:        cfg.ObjectStorage.BucketName,
		ObjectPath:        cfg.ObjectStorage.ObjectPath,
		AccessKeyID:       cfg.ObjectStorage.AccessKeyID,
		SecretAccessKey:   cfg.ObjectStorage.SecretAccessKey,
		EndpointURL:       cfg.ObjectStorage.EndpointURL,
		PublicEndpointURL: cfg.ObjectStorage.PublicEndpointURL,
	})
	r.NoError(t, err)

	file, err := os.ReadFile("test.jpg")
	r.NoError(t, err)

	uploadedUrl, err := cloudflare.UploadFile("test.jog", file, "image/jpeg", "")
	r.NoError(t, err)
	t.Log("Uploaded url:", uploadedUrl)

}
