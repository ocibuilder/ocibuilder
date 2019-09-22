## Environment Variables

The ocibuilder supports two ways of parametizing your specification file. This is either through values specified in your
spec.yaml directly or through referring to a system environment variable.

Parameters are defined in the *params* field in your spec.yaml and refer to a destination of field to replace and the value you
want the field to  replaced with.

>**NOTE**: A specific array item is referred to by index  in the dest field. For example, if you want to access the first step
element you would have ``steps.0``

*e.g.*

```yaml
params:
  # Replaces the value in location build.steps.0.tag with 0.0.3
  - dest: build.steps.0.tag
    value: 0.0.3
  # Replaces the value in location build.steps.0.metadata.name with the environment variable $BUILD_DEV
  - dest: build.steps.0.metadata.name
    valueFromEnv: BUILD_DEV
```

If you specify a valueFromEnv with a value that has not been set, a warning will be returned stating that your environment
variable is empty.