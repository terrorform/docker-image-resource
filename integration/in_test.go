package integration_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const inRequest = `{
	"source": {
		"repository": "library/alpine",
		"tag": "3.1"
	},
	"params": {
		"rootfs": true,
		"skip_download": false,
		"save": false
	},
	"version": {
		"digest": %q
	}
}`

const digest = `d9477888b78e8c6392e0be8b2e73f8c67e2894ff9d4b8e467d1488fcceec21c8`

var _ = FDescribe("In", func() {
	Context("when using dockerhub", func() {
		Context("when rootfs is true", func() {
			var (
				contentDir string
				cmd        *exec.Cmd
			)

			BeforeEach(func() {
				var err error
				contentDir, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())

				req := fmt.Sprintf(inRequest, digest)
				cmd = exec.Command(builtIn, contentDir)
				cmd.Stdin = bytes.NewBufferString(req)
			})

			AfterEach(func() {
				err := os.RemoveAll(contentDir)
				Expect(err).NotTo(HaveOccurred())
			})

			It("downloads the rootfs and metadata", func() {
				sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess).Should(gexec.Exit(0))

				Expect(filepath.Join(contentDir, "image")).NotTo(BeARegularFile())

				repoName, err := ioutil.ReadFile(filepath.Join(contentDir, "repository"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(repoName)).To(Equal("library/alpine"))

				tag, err := ioutil.ReadFile(filepath.Join(contentDir, "tag"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(tag)).To(Equal("3.1"))

				imageID, err := ioutil.ReadFile(filepath.Join(contentDir, "image-id"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(imageID)).To(Equal(""))

				digest, err := ioutil.ReadFile(filepath.Join(contentDir, "digest"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(digest)).To(Equal(digest))

				metadata, err := ioutil.ReadFile(filepath.Join(contentDir, "metadata.json"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(metadata)).To(Equal(`{}`))

				dockerInspect, err := ioutil.ReadFile(filepath.Join(contentDir, "metadata_inspect.json"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(dockerInspect)).To(Equal(`{}`))

				Expect(filepath.Join(contentDir, "rootfs")).To(BeARegularFile())
			})
		})
	})

	Context("when using ecr", func() {
		It("downloads the rootfs and metadata", func() {
		})
	})

	Context("when using gcs", func() {
		It("downloads the rootfs and metadata", func() {
		})
	})
})
