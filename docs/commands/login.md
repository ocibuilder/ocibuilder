## Login

Logs into all registries defined in the specifcation.

>**NOTE:** Login functionality with the docker flavour of the builder does not currently maintain logged in sessions. This
will be resolved in a future release.


#### Synopsis

This command logs into all registries that have been defined in the specification. You can login with a number of different credentials.
These can be plain, taken from environment variables or kubernetes secrets.


#### Command

```
ocictl login [flags] [...]
```

#### Options

```
-p, --path          Path to your spec.yaml or login.yaml. By default will look in the current working directory
-b, --builder       Choose either docker and buildah as the targetted image puller. By default the builder is docker.
-d, --debug         Turn on debug logging
```

### Example

Command
```
ocictl login --path ./spec.yaml
```