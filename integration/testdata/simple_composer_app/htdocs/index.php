<?php
    require '../vendor/autoload.php';

    $log = new Monolog\Logger('my-log');
    $log->pushHandler(new Monolog\Handler\StreamHandler('php://stdout', Monolog\Logger::WARNING));
    $log->addWarning('SUCCESS');

    $names = $_SERVER['QUERY_STRING'];
    foreach (explode(",", $names) as $name) {
      if (extension_loaded($name)) {
        echo 'SUCCESS: ' . $name . ' loads.';
      }
      else {
        echo 'ERROR: ' . $name . ' failed to load.';
      }
    }
?>
