## Push

Pushes container images to one or multiple image registries. 

#### Synopsis

This command pushes all specified images in the push specification to a corresponding registry. You can list many images to push
and many registries.

In order to complete a push to a repository, both the login and push specifications need to be filled in. A push is run with the
authentication passed in the login spec.

The registry, image and tag are used to create a full qualified image path 

``e.g. my-image-registry.docker.com:1111/myimage/cool-image:0.0.1``

#### Command

```
ocictl push [flags] [...]
```

#### Options

```
-p, --path          Path to your spec.yaml or push.yaml. By default will look in the current working directory
-b, --builder       Choose either docker and buildah as the targetted image builder. By default the builder is docker.
-d, --debug         Turn on debug logging
```

### Example

spec.yaml
```yaml
push:
    - registry: my-image-registry.docker.com:1111
      image: myimage/cool-image
      tag: 0.0.1
```
Command
```
ocictl push --path ./spec.yaml
```

or you can use shorthand option like this: 

```
ocictl push -p ./spec.yaml
```