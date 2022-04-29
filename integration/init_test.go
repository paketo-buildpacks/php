package integration_test

import (
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

var phpBuildpack string

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect

	format.MaxLength = 0

	output, err := exec.Command("bash", "-c", "../scripts/package.sh --version 1.2.3").CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), string(output))

	phpBuildpack, err = filepath.Abs("../build/buildpackage.cnb")
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(10 * time.Second)

	suite := spec.New("Integration", spec.Parallel(), spec.Report(report.Terminal{}))
	suite("Composer", testComposer)
	suite("HTTPD", testPhpHttpd)
	suite("Nginx", testPhpNginx)
	suite("Builtin Server", testPhpBuiltinServer)
	suite("Redis Session Handler", testRedisSessionHandler)
	suite("Memcached Session Handler", testMemcachedSessionHandler)
	suite.Run(t)
}
