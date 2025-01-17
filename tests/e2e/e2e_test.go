package e2e

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAdd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Add Nums Suite")
}

// to run test, run command 'ginkgo'
var _ = Describe("Adding two nums", func() {
	Context("When  3 is added to 6", func() {
		It("Should return 9", func() {
			result := add_nums(3, 6)
			file, err := os.Open("/var/run/oadp-credentials")
			Expect(err).NotTo(HaveOccurred())
			fmt.Println("FILE CONTENTS")
			fmt.Print(file)
			expected := 9
			Expect(expected).To(Equal(result))
		})
	})
})
