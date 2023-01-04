package integration_test

import (
	gocontext "context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testReproducibleBuilds(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	context("building a PHP app that uses the built-in web server and Composer as a package manager", func() {
		var (
			image occam.Image

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
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("creates a two identical images from the same input", func() {
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

			firstID := image.ID

			// Delete the first image
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())

			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(phpBuildpack).
				WithPullPolicy("never").
				WithEnv(map[string]string{
					"BP_PHP_WEB_DIR": "htdocs",
				}).
				WithClearCache().
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			Expect(firstID).To(Equal(image.ID))
		})
	})

	context("building a PHP app that contains vendored Composer packages", func() {
		var (
			image occam.Image

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
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("creates a two identical images from the same input", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(phpBuildpack).
				WithPullPolicy("never").
				WithEnv(map[string]string{
					"BP_PHP_SERVER": "httpd",
				}).
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			firstID := image.ID

			// Delete the first image
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())

			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(phpBuildpack).
				WithPullPolicy("never").
				WithEnv(map[string]string{
					"BP_PHP_SERVER": "httpd",
				}).
				WithClearCache().
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			Expect(firstID).To(Equal(image.ID))
		})
	})

	context("building a PHP app that uses HTTPD as a web server and Composer as a package manager", func() {
		var (
			image occam.Image

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
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("creates a two identical images from the same input", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(phpBuildpack).
				WithPullPolicy("never").
				WithEnv(map[string]string{
					"BP_PHP_SERVER": "httpd",
				}).
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			firstID := image.ID

			// Delete the first image
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())

			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(phpBuildpack).
				WithPullPolicy("never").
				WithEnv(map[string]string{
					"BP_PHP_SERVER": "httpd",
				}).
				WithClearCache().
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			Expect(firstID).To(Equal(image.ID))
		})
	})

	context("building a PHP app that uses Nginx as a web server and Composer as a package manager", func() {
		var (
			image occam.Image

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
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("creates a two identical images from the same input", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(phpBuildpack).
				WithPullPolicy("never").
				WithEnv(map[string]string{
					"BP_PHP_SERVER": "nginx",
				}).
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			firstID := image.ID

			// Delete the first image
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())

			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(phpBuildpack).
				WithPullPolicy("never").
				WithEnv(map[string]string{
					"BP_PHP_SERVER": "nginx",
				}).
				WithClearCache().
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			Expect(firstID).To(Equal(image.ID))
		})
	})

	context("building a PHP app that uses a memcached session handler", func() {
		var (
			image occam.Image

			memcachedContainer occam.Container

			name    string
			source  string
			binding string
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
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())

			// Clean up memcached image
			Expect(docker.Container.Remove.Execute(memcachedContainer.ID)).To(Succeed())
			dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			Expect(err).NotTo(HaveOccurred())

			_, err = dockerClient.ImageRemove(gocontext.Background(), "memcached:latest", types.ImageRemoveOptions{Force: true})
			Expect(err).NotTo(HaveOccurred())
		})

		it("creates a two identical images from the same input", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(phpBuildpack).
				WithPullPolicy("never").
				WithEnv(map[string]string{
					"BP_PHP_SERVER":        "nginx",
					"BP_LOG_LEVEL":         "DEBUG",
					"SERVICE_BINDING_ROOT": "/bindings",
				}).
				WithVolumes(fmt.Sprintf("%s:/bindings/php-memcached-session", binding)).
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			firstID := image.ID

			// Delete the first image
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())

			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(phpBuildpack).
				WithPullPolicy("never").
				WithEnv(map[string]string{
					"BP_PHP_SERVER":        "nginx",
					"BP_LOG_LEVEL":         "DEBUG",
					"SERVICE_BINDING_ROOT": "/bindings",
				}).
				WithVolumes(fmt.Sprintf("%s:/bindings/php-memcached-session", binding)).
				WithClearCache().
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			Expect(firstID).To(Equal(image.ID))
		})
	})

	context("building a PHP app that uses a redis session handler", func() {
		var (
			image occam.Image

			redisContainer occam.Container

			name    string
			source  string
			binding string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())

			source, err = occam.Source(filepath.Join("testdata", "session_handler_apps"))
			Expect(err).NotTo(HaveOccurred())

			binding = filepath.Join(source, "redis_binding")

			redisContainer, err = docker.Container.Run.
				WithPublish("6379").
				Execute("redis:latest")
			Expect(err).NotTo(HaveOccurred())

			ipAddress, err := redisContainer.IPAddressForNetwork("bridge")
			Expect(err).NotTo(HaveOccurred())

			Expect(os.WriteFile(filepath.Join(source, "redis_binding", "host"), []byte(ipAddress), os.ModePerm)).To(Succeed())
		})

		it.After(func() {
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())

			// Clean up redis image
			Expect(docker.Container.Remove.Execute(redisContainer.ID)).To(Succeed())
			dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			Expect(err).NotTo(HaveOccurred())

			_, err = dockerClient.ImageRemove(gocontext.Background(), "redis:latest", types.ImageRemoveOptions{Force: true})
			Expect(err).NotTo(HaveOccurred())
		})

		it("creates a two identical images from the same input", func() {
			var err error
			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(phpBuildpack).
				WithPullPolicy("never").
				WithEnv(map[string]string{
					"BP_PHP_SERVER":        "nginx",
					"BP_LOG_LEVEL":         "DEBUG",
					"SERVICE_BINDING_ROOT": "/bindings",
				}).
				WithVolumes(fmt.Sprintf("%s:/bindings/php-redis-session", binding)).
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			firstID := image.ID

			// Delete the first image
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())

			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(phpBuildpack).
				WithPullPolicy("never").
				WithEnv(map[string]string{
					"BP_PHP_SERVER":        "nginx",
					"BP_LOG_LEVEL":         "DEBUG",
					"SERVICE_BINDING_ROOT": "/bindings",
				}).
				WithVolumes(fmt.Sprintf("%s:/bindings/php-redis-session", binding)).
				WithClearCache().
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			Expect(firstID).To(Equal(image.ID))
		})
	})
}
