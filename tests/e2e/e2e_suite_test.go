package e2e

import (
	"flag"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Common vars obtained from flags passed in ginkgo.
var cloud, credentials, namespace, s3Bucket, BucketFilePath, credSecretRef, instanceName string

func init() {
	flag.StringVar(&cloud, "cloud", "", "Cloud Env.")
	flag.StringVar(&credentials, "credentials", "", "Cloud Credentials file path location")
	flag.StringVar(&BucketFilePath, "velero_bucket", "myBucket", "AWS S3 data file path location")
	flag.StringVar(&namespace, "velero_namespace", "oadp-operator", "Velero Namespace")
	flag.StringVar(&credSecretRef, "creds_secret_ref", "cloud-credentials", "OpenShift Credential secret ref for backup storage location")
	flag.StringVar(&instanceName, "velero_instance_name", "example-velero", "Velero Instance Name")
}

func TestOADPE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OADP E2E Suite")
}
