package e2e

import (
	"flag"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Common vars obtained from flags passed in ginkgo.
var credentials, namespace, bucket, bucketFilePath, credSecretRef, instanceName, region, provider string

func init() {
	flag.StringVar(&credentials, "credentials", "", "Cloud Credentials file path location")
	flag.StringVar(&bucketFilePath, "velero_bucket", "myBucket", "AWS S3 data file path location")
	flag.StringVar(&namespace, "velero_namespace", "oadp-operator", "Velero Namespace")
	flag.StringVar(&region, "region", "us-east-1", "BSL region")
	flag.StringVar(&provider, "provider", "aws", "BSL provider")
	flag.StringVar(&credSecretRef, "creds_secret_ref", "cloud-credentials", "Credential secret ref for backup storage location")
	flag.StringVar(&instanceName, "velero_instance_name", "example-velero", "Velero Instance Name")
}

func TestOADPE2E(t *testing.T) {
	flag.Parse()
	s3Buffer, err := getJsonData(bucketFilePath)
	if err != nil {
		t.Fatalf("Error getting bucket json file: %v", err)
	}
	s3Data, err := decodeJson(s3Buffer) // Might need to change this later on to create s3 for each tests
	if err != nil {
		t.Fatalf("Error decoding json file: %v", err)
	}
	bucket = s3Data["velero-bucket-name"].(string)

	RegisterFailHandler(Fail)
	RunSpecs(t, "OADP E2E Suite")
}
