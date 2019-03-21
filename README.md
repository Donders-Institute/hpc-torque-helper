# hpc-torque-helper

This package implements [gPRC](https://grpc.io/) based server and client library for retrieving Torque/Moab job information that requires elevated privileges such as root or the Torque system admin.

## Build RPMs

The server components can be build into RPMs for the CentOS 7.x.  Use the following command to build RPMs:

```bash
$ make release
```

## Make release on GitHub

Making a release on GitHub can be done with the following command:

```bash
$ VERSION={RELEASE_NUMBER} make github_release
```

where the `{RELEASE_NUMBER}` is the new release number to be created on the repository's release page.  It cannot be an existing release number.
