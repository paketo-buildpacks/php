# Restructuring PHP Buildpacks

## Proposal

The RFC proposes a new set of buildpacks and a order grouping for PHP family of
buildpacks.  This is aimed at simplifying the set of buildpacks by modularizing
them according to their function. Each buildpack should do one thing, and do it
well. These buildpacks should interact with each other via meaningful `build`
and `launch` requirements. This builds a clean separation between the concerns
for building an application and those concerns for running it. This simplified
structure would include the following buildpacks in addition to the existing
Apache HTTPD Server and Nginx Server buildpacks[<sup>1</sup>](#note-1):

* **php-dist**:
  Installs the [`php`](https://www.php.net) distribution, making it available on the `$PATH`
  * provides: `php`
  * requires: none

  This distribution must include popular modules like ldap, memcached, mongodb,
  mysqli, openssl, phpiredis, zlib etc.

* **composer**:
  Installs [`composer`](https://getcomposer.org), a dependency manager for PHP and makes it available on the `$PATH`
  * provides: `composer`
  * requires: none

* **composer-install**:
  Resolves project dependencies, and installs them using `composer`.
  * provides: none
  * requires: `php`, `composer` at `build`

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
  * requires: `php`, `httpd`, at `launch`

  This buildpack generates `php.ini`, `php-fpm.conf` and `httpd.conf`; sets up
  env variables and a start command to run PHP FPM[<sup>2</sup>](#note-2) and
  HTTPD Server.

* **php-nginx**:
  Sets up Nginx as the web server to serve PHP applications.
  * provides: none
  * requires: `php`, `nginx` at `launch`

  This buildpack generates `php.ini`, `php-fpm.conf` and `nginx.conf`; sets up
  env variables and a start command to run PHP FPM[<sup>2</sup>](#note-2) and
  Nginx Server.

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

{{REMOVE THIS SECTION BEFORE RATIFICATION!}}

## Notes

<a name="note-1">1</a>. Per the [Web Server Buidpack Subteam
RFC](https://github.com/paketo-buildpacks/rfcs/blob/master/accepted/0006-web-servers.md),
the Apache HTTPD Server and Nginx Server buildpacks are no more considered to
be part of the PHP family of buildpacks.

<a name="note-2">2</a>. There are two primary ways for adding support for PHP
to a web server – as a native web server module, or as a CGI executable.
PHP-FPM (FastCGI Process Manager) is a FastCGI implementation for PHP, bundled
with the official PHP distribution since version 5.3.3

<a name="note-3">3</a>. The session handler is responsible for storing data
from PHP sessions. By default, PHP uses files but they have severe
scalability/performance limitations.
