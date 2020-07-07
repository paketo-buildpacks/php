<?php
    require '../vendor/autoload.php';

    $log = new Monolog\Logger('my-log');
    $log->pushHandler(new Monolog\Handler\StreamHandler('php://stdout', Monolog\Logger::WARNING));
    $log->addWarning('SUCCESS');

    print "This is an HTTPD app.";
?>
