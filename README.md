worker-ws
=========

Service receive a json payload describing infrastructure and app configuration and materialize those configurations.
terraform and ansible files are generated/updated, versionated on gitlab and if requested, call terraform and ansible to provision/configure/keep everything.
