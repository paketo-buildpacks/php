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

func testPhpBuiltinServer(t *testing.T, context spec.G, it spec.S) {
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

	context("building a PHP app that uses the built-in web server and Composer as a package manager", func() {
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
			source, err = occam.Source(filepath.Join("testdata", "simple_composer_app"))
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
				WithPullPolicy("never").
				WithEnv(map[string]string{
					"BP_PHP_WEB_DIR": "htdocs",
				}).
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			container, err = docker.Container.Run.
				WithEnv(map[string]string{"PORT": "8080"}).
				WithPublish("8080").
				WithPublishAll().
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(Serve(ContainSubstring("SUCCESS: date loads.")).OnPort(8080).WithEndpoint("/index.php?date"))

			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Distribution")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Composer")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Composer Install")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Built-in Server")))
			Expect(logs).NotTo(ContainLines(ContainSubstring("Paketo for Procfile")))
			Expect(logs).NotTo(ContainLines(ContainSubstring("Paketo for Environment Variables")))
			Expect(logs).NotTo(ContainLines(ContainSubstring("Paketo for Image Labels")))
		})

		context("using optional utility buildpacks", func() {
			it.Before(func() {
				Expect(os.WriteFile(filepath.Join(source, "Procfile"), []byte(`web: php -S 0.0.0.0:"${PORT:-80}" -t htdocs && echo hi`), 0644)).To(Succeed())
			})

			it.After(func() {
				Expect(os.Remove(filepath.Join(source, "Procfile"))).To(Succeed())
			})

			it("creates a working OCI image and uses the Procfile, Environment Variables, and Image Labels buildpacks", func() {
				var err error
				var logs fmt.Stringer
				image, logs, err = pack.WithNoColor().Build.
					WithBuildpacks(phpBuildpack).
					WithPullPolicy("never").
					WithEnv(map[string]string{
						"BP_PHP_WEB_DIR":    "htdocs",
						"BPE_SOME_VARIABLE": "stew-peas",
						"BP_IMAGE_LABELS":   "cool-label=cool-value",
					}).
					Execute(name, source)
				Expect(err).NotTo(HaveOccurred(), logs.String())

				container, err = docker.Container.Run.
					WithEnv(map[string]string{"PORT": "8080"}).
					WithPublish("8080").
					WithPublishAll().
					Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				Eventually(container).Should(Serve(ContainSubstring("SUCCESS: date loads.")).OnPort(8080).WithEndpoint("/index.php?date"))
				Expect(logs).To(ContainLines(ContainSubstring(`web: php -S 0.0.0.0:"${PORT:-80}" -t htdocs && echo hi`)))

				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Distribution")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Composer")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Composer Install")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Built-in Server")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Procfile")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Environment Variables")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Image Labels")))

				Expect(image.Buildpacks[5].Key).To(Equal("paketo-buildpacks/environment-variables"))
				Expect(image.Buildpacks[5].Layers["environment-variables"].Metadata["variables"]).To(Equal(map[string]interface{}{"SOME_VARIABLE": "stew-peas"}))
				Expect(image.Labels["cool-label"]).To(Equal("cool-value"))
			})
		})
	})
}
