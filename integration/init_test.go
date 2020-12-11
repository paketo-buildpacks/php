package integration_test

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/paketo-buildpacks/occam"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var phpBuildpack string

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect

	bash := pexec.NewExecutable("bash")
	buffer := bytes.NewBuffer(nil)
	err := bash.Execute(pexec.Execution{
		Args:   []string{"-c", "../scripts/package.sh --version 1.2.3"},
		Stdout: buffer,
		Stderr: buffer,
	})
	Expect(err).NotTo(HaveOccurred(), buffer.String)

	phpBuildpack, err = filepath.Abs("../build/buildpackage.cnb")
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(10 * time.Second)

	suite := spec.New("Integration", spec.Parallel(), spec.Report(report.Terminal{}))
	suite("composer with nginx", testPhpNginx)
	suite("composer with httpd", testPhpHttpd)
	suite.Run(t)
}

func ContainerLogs(id string) func() string {
	docker := occam.NewDocker()

	return func() string {
		logs, _ := docker.Container.Logs.Execute(id)
		return logs.String()
	}
}

func BeAvailableAndReady() types.GomegaMatcher {
	return &BeAvailableAndReadyMatcher{
		Docker: occam.NewDocker(),
	}
}

type BeAvailableAndReadyMatcher struct {
	Docker occam.Docker
}

func (*BeAvailableAndReadyMatcher) Match(actual interface{}) (bool, error) {
	container, ok := actual.(occam.Container)
	if !ok {
		return false, fmt.Errorf("BeAvailableMatcher expects an occam.Container, received %T", actual)
	}

	response, err := http.Get(fmt.Sprintf("http://localhost:%s", container.HostPort("8080")))
	if err != nil {
		return false, nil
	}

	if response.StatusCode != http.StatusOK {
		return false, nil
	}

	defer response.Body.Close()

	return true, nil
}

func (m *BeAvailableAndReadyMatcher) FailureMessage(actual interface{}) string {
	container := actual.(occam.Container)
	message := fmt.Sprintf("Expected\n\tdocker container id: %s\nto be available.", container.ID)

	if logs, _ := m.Docker.Container.Logs.Execute(container.ID); logs != nil {
		message = fmt.Sprintf("%s\n\nContainer logs:\n\n%s", message, logs)
	}

	return message
}

func (m *BeAvailableAndReadyMatcher) NegatedFailureMessage(actual interface{}) string {
	container := actual.(occam.Container)
	message := fmt.Sprintf("Expected\n\tdocker container id: %s\nnot to be available.", container.ID)

	if logs, _ := m.Docker.Container.Logs.Execute(container.ID); logs != nil {
		message = fmt.Sprintf("%s\n\nContainer logs:\n\n%s", message, logs)
	}

	return message
}
