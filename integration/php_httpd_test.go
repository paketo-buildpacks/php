package integration_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testPhpHttpd(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	context("when building a php app using php-web, php-composer, and httpd", func() {
		var (
			image     occam.Image
			container occam.Container

			name string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
		})

		it("creates a working OCI image", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(phpBuildpack).
				WithPullPolicy("never").
				Execute(name, filepath.Join("testdata", "offline_composer_httpd"))
			Expect(err).NotTo(HaveOccurred(), logs.String())

			container, err = docker.Container.Run.Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(BeAvailableAndReady(), ContainerLogs(container.ID))

			response, err := http.Get(fmt.Sprintf("http://localhost:%s", container.HostPort()))
			Expect(err).NotTo(HaveOccurred())
			defer response.Body.Close()

			content, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(MatchRegexp("This is an HTTPD app."))

			Expect(logs).To(ContainLines(ContainSubstring("PHP Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("Apache HTTP Server Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("PHP Web Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("PHP Composer Buildpack")))
		})
	})
}
