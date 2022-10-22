# Gardena Smart System Exporter

The Gardena Smart System Exporter is a golang based application to monitor the current status of your Gardena Smart
System devices, like the mower roboter or the ignition controller.
It collects the current status of your devices and exposes the information as metrics in the Prometheus format.

The project is considered work in progress.

## Authentication

The Gardena Smart System Exporter authenticates with the Gardena api via client-id and client-secret, aka
'Application Key' and 'Application Secret', see [Developer Page](https://developer.husqvarnagroup.cloud/docs/get-started).

Once created, those credentials need to be saved to the files `client-id` and `client-secret`. The exporter expects those
files under the path `/etc/secrets/gardena-smart-system-exporter`. The path can be changed by setting the
`secret-file-path` flag. 

For development, you can also store the credentials in files provided in the `/config` directory and hide them from vcs
by running `git update-index --no-assume-unchanged <file>`.
DO NOT COMMIT THE CREDENTIALS SINCE IT GIVES ACCESS TO ALL YOUR DEVISES!
