## Pull

Pulls an image passed with the name flag.

#### Synopsis

This command pulls an image that you have passed in by name. The name should include the path to the image but not
the image registry itself.

``e.g. myimage/cool-image:0.0.1``

The pull command looks to pull from any registries that have been specified in the login specification. Once the image has
been found in any of the specified registries, a pull is executed.

#### Command

```
ocictl pull [flags] [...]
```

#### Options

```
-i, --name          Specify the name of the image you want to pull
-p, --path          Path to your spec.yaml. By default will look in the current working directory
-b, --builder       Choose either docker and buildah as the targetted image puller. By default the builder is docker.
-d, --debug         Turn on debug logging
```

### Example

Command
```
ocictl pull --path ./spec.yaml --name myimage/cool-image:0.0.1
```