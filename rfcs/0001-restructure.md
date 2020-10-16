# Restructuring PHP Buildpacks

## Proposal

The RFC proposes a new set of smaller buildpacks that, by themselves provides
more limited functionality, but in composition provides a powerful way to build
various kinds of images related to PHP.

## Motivation

1. Reducing the responsibility of each buildpack will help to make them easier
   to understand and maintain. We will end up with more buildpacks, each with
   simpler implementations.

1. Align buildplan philosophy with our other buildpacks. We have moved away
   from the buildpack unilaterally marking its layers as `build` or `launch`.
   Buildpacks should interact with each other via meaningful `build` and
   `launch` requirements. This builds a clean separation between the concerns
   of building and running an application.

The new structure would include the following buildpacks in addition to the existing
Apache HTTPD Server and Nginx Server buildpacks[<sup>1</sup>](#note-1):

* **php-dist**:
  Installs the [`php`](https://www.php.net) distribution, making it available on the `$PATH`
  * provides: `php`
  * requires: none

  This distribution must include popular modules like memcached, mongodb,
  mysqli, openssl, phpiredis, zlib etc; and a suitable `php.ini`. Downstream
  buildpacks are expected to obtain info (like extension dir) through the
  [`php-config`](https://www.php.net/manual/en/install.pecl.php-config.php)
  executable.

* **composer**:
  Installs [`composer`](https://getcomposer.org), a dependency manager for PHP
  and makes it available on the `$PATH`
  * provides: `composer`
  * requires: none

* **composer-install**:
  Resolves project dependencies, and installs them using `composer`.
  * provides: none
  * requires: `php`, `composer` at `build`

* **php-fpm**: (??)
  Configures `php-fpm.conf` (config file in `php.ini` syntax), and sets a start
  command.
  FPM[<sup>2</sup>](#note-2).
  * provides: `php-fpm`
  * requires: `php` during launch

  Separation of php-fpm into a separate buildpack lets users run FPM in one
  container and web server in another container.
  The generated `php-fpm.conf` must accept requests on a unix socket, must
  "include" custom php-fpm configs specified by the user (e.g. in buildpack.yml
  or specific location in the app) and write the generated config file to a
  well known location (like setting an env var or writing to <workingDir>).

* **php-builtin-server**:
  Set up PHP's [built-in web
  server](https://www.php.net/manual/en/features.commandline.webserver.php) to
  serve PHP applications.
  * provides: none
  * requires: `php` at `launch`

  This buildpack generates `php.ini`, sets up env variables and a start command
  to start the built-in web server.

* **php-httpd**:
  Sets up HTTPD as the web server to serve PHP applications.
  * provides: none
  * requires: `php`, `php-fpm`, `httpd`, at `launch`

  This buildpack generates `httpd.conf` and sets up a start command to run PHP
  FPM and HTTPD Server. Apps need to declare the intention to use httpd in
  `buildpack.toml`:

* **php-nginx**:
  Sets up Nginx as the web server to serve PHP applications.
  * provides: none
  * requires: `php`, `php-fpm`, `nginx` at `launch`

  This buildpack generates `nginx.conf`, sets up a start command to run PHP FPM
  and Nginx Server. Apps need to declare the intention to use nginx in
  `buildpack.toml`:

* **php-memcached-session-handler**:
  Configures the given memcached service instance as a PHP session
  handler[<sup>2</sup>](#note-2). Memcached settings are to be provided through
  a suitable
  [binding](https://paketo.io/docs/buildpacks/configuration/#bindings).
  * provides: none
  * requires: php at `launch`

* **php-redis-session-handler**:
  Configures the given redis service instance as a PHP session
  handler[<sup>3</sup>](#note-3). Redis settings are to be provided through a
  suitable
  [binding](https://paketo.io/docs/buildpacks/configuration/#bindings).
  * provides: none
  * requires: php at `launch`


This would result in the following order groupings in the PHP language family meta-buildpack:

```toml
[[order]] # HTTPD web server

  [[order.group]]
    id = "paketo-buildpacks/php-dist"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/composer"
    version = ""
    optional = true

  [[order.group]]
    id = "paketo-buildpacks/composer-install"
    version = ""
    optional = true

  [[order.group]]
    id = "paketo-buildpacks/httpd"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/php-fpm"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/php-httpd"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/php-memcached-session-handler"
    version = ""
    optional = true

  [[order.group]]
    id = "paketo-buildpacks/php-redis-session-handler"
    version = ""
    optional = true


[[order]] # Nginx web server

  [[order.group]]
    id = "paketo-buildpacks/php-dist"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/composer"
    version = ""
    optional = true

  [[order.group]]
    id = "paketo-buildpacks/composer-install"
    version = ""
    optional = true

  [[order.group]]
    id = "paketo-buildpacks/nginx"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/php-fpm"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/php-nginx"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/php-memcached-session-handler"
    version = ""
    optional = true

  [[order.group]]
    id = "paketo-buildpacks/php-redis-session-handler"
    version = ""
    optional = true

[[order]] # Built-in web server

  [[order.group]]
    id = "paketo-buildpacks/php-dist"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/composer"
    version = ""
    optional = true

  [[order.group]]
    id = "paketo-buildpacks/composer-install"
    version = ""
    optional = true

  [[order.group]]
    id = "paketo-buildpacks/php-builtin-server"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/php-memcached-session-handler"
    version = ""
    optional = true

  [[order.group]]
    id = "paketo-buildpacks/php-redis-session-handler"
    version = ""
    optional = true
```

## Unresolved Questions and Bikeshedding

* php-httpd/php-nginx overrides the start cmd set by php-fpm. Does the spec allow for a way for
start command set by 2 buildpacks to run together? (Different process-types?)

* How do you have a group such that it produces an image who will simply run php-fpm? There's no
way to differentiate that group from php-built-in group.

* Is php-fpm worth the separation? (If not, we can document it as a possible alternate solution)

* What's the best way to detect if the app should be run on builtin/httpd/nginx?

{{REMOVE THIS SECTION BEFORE RATIFICATION!}}

## Notes

<a name="note-1">1</a>. Per the [Web Server Buidpack Subteam
RFC](https://github.com/paketo-buildpacks/rfcs/blob/master/accepted/0006-web-servers.md),
the Apache HTTPD Server and Nginx Server buildpacks are no more considered to
be part of the PHP family of buildpacks.

<a name="note-2">2</a>. There are a few ways for adding support for PHP to a
web server – as a native web server module, using CGI, using FastCGI. PHP-FPM
(FastCGI Process Manager) is a FastCGI implementation for PHP, bundled with the
official PHP distribution since version 5.3.3

<a name="note-3">3</a>. The session handler is responsible for storing data
from PHP sessions. By default, PHP uses files but they have severe scalability
limitations. With external session handlers, multiple application nodes can
connect to a central data store.
