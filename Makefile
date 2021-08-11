OADP_TEST_NAMESPACE ?= oadp-operator
CREDS_SECRET_REF ?= cloud-credentials
VELERO_INSTANCE_NAME ?= example-velero
CLUSTER_PROFILE ?= aws
OADP_CRED_FILE ?= /var/run/oadp-credentials/${CLUSTER_PROFILE}-credentials
OADP_BUCKET ?= /var/run/oadp-credentials/velero-bucket-name

.PHONY:ginkgo
ginkgo: # Make sure ginkgo is in $GOPATH/bin
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega/...

test-e2e:
	ginkgo tests/e2e/ -- -cloud=$(CLUSTER_PROFILE) -credentials=$(OADP_CRED_FILE) \
	-velero_bucket=$(OADP_BUCKET) -velero_namespace=$(OADP_TEST_NAMESPACE) \
	-creds_secret_ref=$(CREDS_SECRET_REF) \
	-velero_instance_name=$(VELERO_INSTANCE_NAME)
