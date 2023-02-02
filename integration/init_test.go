package integration_test

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

var (
	phpBuildpack   string
	builder        occam.Builder
	memcachedImage string
	redisImage     string
)

func TestIntegration(t *testing.T) {
	pack := occam.NewPack()
	docker := occam.NewDocker()
	Expect := NewWithT(t).Expect

	format.MaxLength = 0

	output, err := exec.Command("bash", "-c", "../scripts/package.sh --version 1.2.3").CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), string(output))

	phpBuildpack, err = filepath.Abs("../build/buildpackage.cnb")
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(10 * time.Second)

	builder, err = pack.Builder.Inspect.Execute()
	Expect(err).NotTo(HaveOccurred())

	// pull and re-tag memcached image with builder-specific naming
	// this will prevent flakes in which we try to reference/remove the same image in parallel bionic/jammy builder tests
	memcachedImage = fmt.Sprintf("memcached-%s:latest", builder.LocalInfo.Stack.ID)
	Expect(docker.Pull.Execute("memcached:latest")).To(Succeed())
	Expect(docker.Image.Tag.Execute("memcached:latest", memcachedImage)).To(Succeed())
	Expect(docker.Image.Remove.WithForce().Execute("memcached:latest")).To(Succeed())

	// pull and re-tag redis image with builder-specific naming
	// this will prevent flakes in which we try to reference/remove the same image in parallel bionic/jammy builder tests
	redisImage = fmt.Sprintf("redis-%s:latest", builder.LocalInfo.Stack.ID)
	Expect(docker.Pull.Execute("redis:latest")).To(Succeed())
	Expect(docker.Image.Tag.Execute("redis:latest", redisImage)).To(Succeed())
	Expect(docker.Image.Remove.WithForce().Execute("redis:latest")).To(Succeed())

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
	Expect(docker.Image.Remove.WithForce().Execute(memcachedImage)).To(Succeed())
	// Clean up redis image
	Expect(docker.Image.Remove.WithForce().Execute(redisImage)).To(Succeed())
}
