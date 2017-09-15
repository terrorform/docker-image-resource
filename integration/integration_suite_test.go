package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var builtIn string

func TestDockerImageResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration")
}

var _ = BeforeSuite(func() {
	var err error
	builtIn, err = gexec.Build("github.com/concourse/docker-image-resource/cmd/in")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
