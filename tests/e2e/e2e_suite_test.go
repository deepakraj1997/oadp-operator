package e2e

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Common vars obtained from flags passed in ginkgo.
var credFile, namespace, credSecretRef, instanceName, provider, azure_resource_file, openshift_ci, ci_cred_file, settings, bsl_profile string
var timeoutMultiplier time.Duration

func init() {
	flag.StringVar(&credFile, "credentials", "", "Cloud Credentials file path location")
	flag.StringVar(&namespace, "velero_namespace", "velero", "Velero Namespace")
	flag.StringVar(&settings, "settings", "./templates/default_settings.json", "Settings of the velero instance")
	flag.StringVar(&instanceName, "velero_instance_name", "example-velero", "Velero Instance Name")
	flag.StringVar(&bsl_profile, "cluster_profile", "aws", "Cluster profile")
	flag.StringVar(&credSecretRef, "creds_secret_ref", "cloud-credentials", "Credential secret ref for backup storage location")
	flag.StringVar(&provider, "provider", "aws", "BSL provider")
	flag.StringVar(&azure_resource_file, "azure_resource_file", "azure resource file", "Resource Group Dir for azure")
	flag.StringVar(&ci_cred_file, "ci_cred_file", credFile, "CI Cloud Cred File")
	flag.StringVar(&openshift_ci, "openshift_ci", "false", "ENV for tests")

	timeoutMultiplierInput := flag.Int64("timeout_multiplier", 1, "Customize timeout multiplier from default (1)")
	timeoutMultiplier = 1
	if timeoutMultiplierInput != nil && *timeoutMultiplierInput >= 1 {
		timeoutMultiplier = time.Duration(*timeoutMultiplierInput)
	}
}

func TestOADPE2E(t *testing.T) {
	flag.Parse()
	errString := loadDpaSettingsFromJson(settings)
	if errString != "" {
		t.Fatalf(errString)
	}

	log.Println("Using velero prefix: " + veleroPrefix)
	RegisterFailHandler(Fail)
	RunSpecs(t, "OADP E2E Suite")
}

var dpaCR *dpaCustomResource

var _ = BeforeSuite(func() {
	flag.Parse()
	errString := loadDpaSettingsFromJson(settings)
	if errString != "" {
		Expect(errors.New(errString)).NotTo(HaveOccurred())
	}

	dpaCR = &dpaCustomResource{
		Namespace:     namespace,
		Credentials:   credFile,
		CredSecretRef: credSecretRef,
	}
	dpaCR.CustomResource = dpa
	testSuiteInstanceName := "ts-" + instanceName
	dpaCR.Name = testSuiteInstanceName
	openshift_ci_bool, _ := strconv.ParseBool(openshift_ci)
	dpaCR.OpenshiftCi = openshift_ci_bool

	if openshift_ci_bool == true {
		switch dpaCR.Provider {
		case "aws":
			cloudCredData, err := readFile(dpaCR.Credentials)
			Expect(err).NotTo(HaveOccurred())
			ciCredData, err := readFile(ci_cred_file)
			Expect(err).NotTo(HaveOccurred())
			cloudCredData = append(cloudCredData, []byte("\n")...)
			credData := append(cloudCredData, ciCredData...)
			dpaCR.Credentials = "/tmp/aws-credentials"
			err = writeFile(dpaCR.Credentials, credData)
			Expect(err).NotTo(HaveOccurred())
		case "gcp":
			cloudCredData, err := readFile(dpaCR.Credentials)
			Expect(err).NotTo(HaveOccurred())
			err = createCredentialsSecret(cloudCredData, namespace, "bsl-cloud-credentials-gcp")
			Expect(err).NotTo(HaveOccurred())
			dpaCR.Credentials = ci_cred_file
		case "azure":
			cloudCredData, err := getJsonData(dpaCR.Credentials) // azure credentials need to be in json - can be changed

			Expect(err).NotTo(HaveOccurred())
			dpaCR.DpaAzureConfig = dpaAzureConfig{
				BslSubscriptionId:          fmt.Sprintf("%v", cloudCredData["subscriptionId"]),
				BslResourceGroup:           fmt.Sprintf("%v", cloudCredData["resourceGroup"]),
				BslStorageAccount:          fmt.Sprintf("%v", cloudCredData["storageAccount"]),
				BslStorageAccountKeyEnvVar: "AZURE_STORAGE_ACCOUNT_ACCESS_KEY",
				VslSubscriptionId:          fmt.Sprintf("%v", cloudCredData["subscriptionId"]),
				VslResourceGroup:           fmt.Sprintf("%v", cloudCredData["resourceGroup"]),
			}

			// bsl cloud
			cloudCreds := getAzureCreds(cloudCredData)
			err = createCredentialsSecret(cloudCreds, namespace, "bsl-cloud-credentials-azure")
			Expect(err).NotTo(HaveOccurred())
			// ci cloud
			ciJsonData, err := getJsonData(ci_cred_file)
			Expect(err).NotTo(HaveOccurred())
			if _, ok := ciJsonData["resourceGroup"]; !ok {
				resourceGroup, err := getAzureResource(azure_resource_file)
				Expect(err).NotTo(HaveOccurred())
				ciJsonData["resourceGroup"] = resourceGroup
			}
			dpaCR.DpaAzureConfig.VslSubscriptionId = fmt.Sprintf("%v", ciJsonData["subscriptionId"])
			dpaCR.DpaAzureConfig.VslResourceGroup = fmt.Sprintf("%v", ciJsonData["resourceGroup"])
			ciCreds := getAzureCreds(ciJsonData)
			dpaCR.Credentials = "/tmp/azure-credentials"
			err = writeFile(dpaCR.Credentials, ciCreds)
			Expect(err).NotTo(HaveOccurred())
		}
	}
	credData, err := readFile(dpaCR.Credentials)
	Expect(err).NotTo(HaveOccurred())
	err = createCredentialsSecret(credData, namespace, credSecretRef)
	Expect(err).NotTo(HaveOccurred())
	dpaCR.SetClient()
	Expect(doesNamespaceExist(namespace)).Should(BeTrue())
})

var _ = AfterSuite(func() {
	log.Printf("Deleting Velero CR")
	errs := deleteSecret(namespace, getSecretRef(credSecretRef))
	Expect(errs).ToNot(HaveOccurred())
	err := dpaCR.Delete()
	Expect(err).ToNot(HaveOccurred())
	Eventually(dpaCR.IsDeleted(), timeoutMultiplier*time.Minute*2, time.Second*5).Should(BeTrue())
})
