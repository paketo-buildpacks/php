package integration_test

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testPhpNginx(t *testing.T, context spec.G, it spec.S) {
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

	context("building a PHP app that uses Nginx as a web server and Composer as a package manager", func() {
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
				WithEnv(map[string]string{
					"BP_PHP_SERVER": "nginx",
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

			Eventually(container).Should(Serve(ContainSubstring("SUCCESS: date loads.")).OnPort(8080).WithEndpoint("/index.php?date"))

			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for CA Certificates")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Distribution")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Composer")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Composer Install")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Nginx Server")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP FPM")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Nginx")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Start")))
			Expect(logs).NotTo(ContainLines(ContainSubstring("Paketo Buildpack for Procfile")))
			Expect(logs).NotTo(ContainLines(ContainSubstring("Paketo Buildpack for Environment Variables")))
			Expect(logs).NotTo(ContainLines(ContainSubstring("Paketo Buildpack for Image Labels")))

			// Ensure FPM is running as well
			Eventually(func() string {
				cLogs, err := docker.Container.Logs.Execute(container.ID)
				Expect(err).NotTo(HaveOccurred())
				return cLogs.String()
			}).Should(
				And(
					ContainSubstring("NOTICE: fpm is running"),
					ContainSubstring("NOTICE: ready to handle connections"),
				),
			)
		})

		context("using optional utility buildpacks", func() {
			it.Before(func() {
				Expect(os.WriteFile(filepath.Join(source, "Procfile"), []byte("web: procmgr-binary /layers/paketo-buildpacks_php-start/php-start/procs.yml && echo hi"), 0644)).To(Succeed())
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
						"BP_PHP_SERVER":          "nginx",
						"BPE_SOME_VARIABLE":      "fish-n-chips",
						"BP_IMAGE_LABELS":        "cool-label=cool-value",
						"BP_LIVE_RELOAD_ENABLED": "true",
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

				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for CA Certificates")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Watchexec")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Distribution")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Composer")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Composer Install")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Nginx Server")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP FPM")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Nginx")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Start")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Procfile")))
				Expect(logs).To(ContainLines(ContainSubstring("web: procmgr-binary /layers/paketo-buildpacks_php-start/php-start/procs.yml && echo hi")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Environment Variables")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Image Labels")))

				Expect(image.Buildpacks[10].Key).To(Equal("paketo-buildpacks/environment-variables"))
				Expect(image.Buildpacks[10].Layers["environment-variables"].Metadata["variables"]).To(Equal(map[string]interface{}{"SOME_VARIABLE": "fish-n-chips"}))
				Expect(image.Labels["cool-label"]).To(Equal("cool-value"))
			})
		})

		context("when using CA certificates", func() {
			var (
				client *http.Client
			)

			it.Before(func() {
				var err error
				name, err = occam.RandomName()
				Expect(err).NotTo(HaveOccurred())
				source, err = occam.Source(filepath.Join("testdata", "ca_cert_apps"))
				Expect(err).NotTo(HaveOccurred())

				caCert, err := os.ReadFile(fmt.Sprintf("%s/nginx_app/certs/ca.pem", source))
				Expect(err).ToNot(HaveOccurred())

				caCertPool := x509.NewCertPool()
				caCertPool.AppendCertsFromPEM(caCert)

				cert, err := tls.LoadX509KeyPair(fmt.Sprintf("%s/nginx_app/certs/cert.pem", source), fmt.Sprintf("%s/nginx_app/certs/key.pem", source))
				Expect(err).ToNot(HaveOccurred())

				client = &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							RootCAs:      caCertPool,
							Certificates: []tls.Certificate{cert},
							MinVersion:   tls.VersionTLS12,
						},
					},
				}
			})

			it("builds a working OCI image with given CA cert added to trust store", func() {
				var err error
				var logs fmt.Stringer
				image, logs, err = pack.WithNoColor().Build.
					WithBuildpacks(phpBuildpack).
					WithEnv(map[string]string{
						"BP_PHP_SERVER":             "nginx",
						"BP_PHP_NGINX_ENABLE_HTTPS": "true",
					}).
					WithPullPolicy("never").
					Execute(name, filepath.Join(source, "nginx_app"))
				Expect(err).NotTo(HaveOccurred())

				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for CA Certificates")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Distribution")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Nginx Server")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP FPM")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Nginx")))
				Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Start")))

				container, err = docker.Container.Run.
					WithPublish("8080").
					WithEnv(map[string]string{
						"PORT":                         "8080",
						"BP_PHP_ENABLE_HTTPS_REDIRECT": "false",
						"SERVICE_BINDING_ROOT":         "/bindings",
					}).
					WithVolumes(fmt.Sprintf("%s:/bindings/ca-certificates", filepath.Join(source, "binding"))).
					Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				Eventually(func() string {
					cLogs, err := docker.Container.Logs.Execute(container.ID)
					Expect(err).NotTo(HaveOccurred())
					return cLogs.String()
				}).Should(
					ContainSubstring("Added 1 additional CA certificate(s) to system truststore"),
				)

				Eventually(container).Should(Serve(ContainSubstring("Hello world, Authenticated User!")).OnPort(8080).WithProtocol("https").WithEndpoint("/").WithClient(client))
			})
		})
	})
}
