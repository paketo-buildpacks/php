package integration_test

import (
	gocontext "context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testMemcachedSessionHandler(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		pack   occam.Pack
		docker occam.Docker
		source string
		name   string
	)

	it.Before(func() {
		pack = occam.NewPack().WithVerbose()
		docker = occam.NewDocker()
	})

	context("building a PHP app that uses a memcached session handler", func() {
		var (
			image              occam.Image
			container          occam.Container
			memcachedContainer occam.Container
			binding            string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())

			source, err = occam.Source(filepath.Join("testdata", "session_handler_apps"))
			Expect(err).NotTo(HaveOccurred())
			binding = filepath.Join(source, "memcached_binding")

			memcachedContainer, err = docker.Container.Run.
				WithPublish("11211").
				Execute("memcached")
			Expect(err).NotTo(HaveOccurred())

			ipAddress, err := memcachedContainer.IPAddressForNetwork("bridge")
			Expect(err).NotTo(HaveOccurred())

			Expect(os.WriteFile(filepath.Join(source, "memcached_binding", "host"), []byte(ipAddress), os.ModePerm)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(binding, "servers"), []byte(ipAddress), os.ModePerm)).To(Succeed())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(memcachedContainer.ID)).To(Succeed())
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
			// Clean up memcached image
			dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			Expect(err).NotTo(HaveOccurred())

			_, err = dockerClient.ImageRemove(gocontext.Background(), "memcached:latest", types.ImageRemoveOptions{Force: true})
			Expect(err).NotTo(HaveOccurred())
		})

		it("creates a working OCI image", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(phpBuildpack).
				WithEnv(map[string]string{
					"BP_PHP_SERVER":        "nginx",
					"BP_LOG_LEVEL":         "DEBUG",
					"SERVICE_BINDING_ROOT": "/bindings",
				}).
				WithPullPolicy("never").
				WithVolumes(fmt.Sprintf("%s:/bindings/php-memcached-session", binding)).
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			container, err = docker.Container.Run.
				WithEnv(map[string]string{"PORT": "8080"}).
				WithPublish("8080").
				WithPublishAll().
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			jar, err := cookiejar.New(nil)
			Expect(err).NotTo(HaveOccurred())

			client := &http.Client{
				Jar: jar,
			}

			Eventually(container).Should(Serve(ContainSubstring("1")).WithClient(client).OnPort(8080).WithEndpoint("/index.php"))
			Eventually(container).Should(Serve(ContainSubstring("2")).WithClient(client).OnPort(8080).WithEndpoint("/index.php"))

			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Distribution")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for Nginx Server")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP FPM")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Nginx")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Start")))
			Expect(logs).To(ContainLines(ContainSubstring("Paketo Buildpack for PHP Memcached Session Handler")))
		})
	})
}
