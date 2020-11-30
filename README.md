# PHP Paketo Buildpack
## `gcr.io/paketo-buildpacks/php`

The PHP Paketo Buildpack provides a set of collaborating buildpacks that
enable the building of a PHP-based application. These buildpacks include:
- [PHP Dist CNB](https://github.com/paketo-buildpacks/php-dist)
- [PHP Web CNB](https://github.com/paketo-buildpacks/php-web)
- [PHP Composer CNB](https://github.com/paketo-buildpacks/php-composer)
- [Apache HTTPD CNB](https://github.com/paketo-buildpacks/httpd)
- [NGINX CNB](https://github.com/paketo-buildpacks/nginx)

The buildpack supports building PHP console and web applications. Web
applications can be run on either the [built-in PHP
webserver](https://www.php.net/manual/en/features.commandline.webserver.php),
[Apache HTTPD](https://httpd.apache.org/) or [NGINX](https://www.nginx.com/).
The buildpack also provides optional support for the utilization of
[Composer](https://getcomposer.org) as a package manager.

Usage examples can be found in the
[`samples` repository under the `php` directory](https://github.com/paketo-buildpacks/samples/tree/main/php).

#### The PHP buildpack is compatible with the following builder(s):
- [Paketo Full Builder](https://github.com/paketo-buildpacks/full-builder)
