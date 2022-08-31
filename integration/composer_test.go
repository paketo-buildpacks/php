package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testComposer(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack().WithVerbose()
		docker = occam.NewDocker()
	})

	context("building a PHP app that contains vendored Composer packages", func() {
		var (
			image     occam.Image
			container occam.Container

			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
			source, err = occam.Source(filepath.Join("testdata", "vendored_composer_app"))
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("creates a working OCI image", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(phpBuildpack).
				WithEnv(map[string]string{
					"BP_PHP_SERVER": "httpd",
				}).
				WithPullPolicy("never").
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			container, err = docker.Container.Run.
				WithEnv(map[string]string{"PORT": "8080"}).
				WithPublish("8080").
				WithPublishAll().
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(Serve(ContainSubstring("This is a PHP app.")).OnPort(8080))

			Expect(logs).To(ContainLines(ContainSubstring("Detected existing vendored packages, will run 'composer install' with those packages")))

			Expect(logs).To(ContainLines(ContainSubstring("CA Certificates Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("PHP Distribution Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Composer")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Composer Install")))
			Expect(logs).To(ContainLines(ContainSubstring("Apache HTTP Server Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("PHP FPM Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("PHP HTTPD Buildpack")))
			Expect(logs).To(ContainLines(ContainSubstring("PHP Start Buildpack")))
		})
	})
}
