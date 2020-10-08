# Restructuring php buildpacks

## Proposal

The RFC proposes a new set of buildpacks and a order grouping for PHP family of
buildpacks.  This is aimed at simplifying the set of buildpacks by modularizing
them according to their function. Each buildpack should do one thing, and do it
well. These buildpacks should interact with each other via meaningful `build`
and `launch` requirements. This builds a clean separation between the concerns
for building an application and those concerns for running it. This simplified
structure would include the following buildpacks in addition to the existing
Apache HTTPD Server and Nginx Server buildpacks:

* **php-dist**:
  Installs the [`php`](https://www.php.net) interpreter, making it available on the `$PATH`
  * provides: `php`
  * requires: none

* **composer**:
  Installs the [`composer`](https://getcomposer.org), a dependency manager for PHP and makes it available on the `$PATH`
  * provides: `composer`
  * requires: none

* **composer-install**:
  Resolves project dependencies, and installs them using `composer`.
  * provides: none
  * requires: `php`, `composer`

* **php-web-server**:
  Sets a launch time start command to start PHP's [built-in web
  server](https://www.php.net/manual/en/features.commandline.webserver.php).
  * provides: none
  * requires: `php`

* **php-httpd**:
  Generates a suitable `httpd.conf` to use httpd as the web server to serve PHP
  application, and sets a launch time start command to start httpd.
  * provides: none
  * requires: `php`

* **php-nginx**:
  Generates a suitable `nginx.conf` to use nginx as the web server to serve PHP
  application, and sets a launch time start command to start nginx.
  * provides: none
  * requires: `php`

* **php-memcached-session-handler**:
  Configures the given memcached service instance as a PHP session handler.
  Memcached settings are to be provided through a suitable
  [binding](https://paketo.io/docs/buildpacks/configuration/#bindings). The
  session handler is responsible for storing data from PHP sessions. By
  default, PHP uses files but they have severe scalability/performance
  limitations.

  * provides: none
  * requires: `php`

* **php-redis-session-handler**:
  Configures the given redis service instance as a PHP session handler.  Redis
  settings are to be provided through a suitable
  [binding](https://paketo.io/docs/buildpacks/configuration/#bindings).

  * provides: none
  * requires: `php`


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
    id = "paketo-buildpacks/php-web-server"
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

## Notes

Per the [Web Server Buidpack Subteam
RFC](https://github.com/paketo-buildpacks/rfcs/blob/master/accepted/0006-web-servers.md),
the Apache HTTPD Server and Nginx Server buildpacks are no more considered to
be part of the PHP family of buildpacks.


## Unresolved Questions and Bikeshedding

{{REMOVE THIS SECTION BEFORE RATIFICATION!}}
