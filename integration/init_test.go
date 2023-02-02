package integration_test

import (
	gocontext "context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

var phpBuildpack string

func TestIntegration(t *testing.T) {
	pack := occam.NewPack()
	Expect := NewWithT(t).Expect

	format.MaxLength = 0

	output, err := exec.Command("bash", "-c", "../scripts/package.sh --version 1.2.3").CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), string(output))

	phpBuildpack, err = filepath.Abs("../build/buildpackage.cnb")
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(10 * time.Second)

	builder, err := pack.Builder.Inspect.Execute()
	Expect(err).NotTo(HaveOccurred())

	suite := spec.New("Integration", spec.Parallel(), spec.Report(report.Terminal{}))

	// This test will only run on the Bionic full stack, to test the stack upgrade scenario.
	// All other tests will run against the Bionic full stack and Jammy full stack
	if builder.BuilderName == "paketobuildpacks/builder:buildpackless-full" {
		suite("StackUpgrades", testStackUpgrades)
	}

	suite("Builtin Server", testPhpBuiltinServer)
	suite("Composer", testComposer)
	suite("HTTPD", testPhpHttpd)
	suite("Memcached Session Handler", testMemcachedSessionHandler)
	suite("Nginx", testPhpNginx)
	suite("Redis Session Handler", testRedisSessionHandler)
	suite("Reproducible Builds", testReproducibleBuilds)
	suite.Run(t)

	// Clean up memcached image
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	Expect(err).NotTo(HaveOccurred())

	_, err = dockerClient.ImageRemove(gocontext.Background(), "memcached:latest", types.ImageRemoveOptions{Force: true})
	Expect(err).NotTo(HaveOccurred())

	_, err = dockerClient.ImageRemove(gocontext.Background(), "redis:latest", types.ImageRemoveOptions{Force: true})
	Expect(err).NotTo(HaveOccurred())
}
