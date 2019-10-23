# Contributing to Ocibuilder

## Filing issues

When filing an issue, make sure to answer these five questions:

1. What version of Go are you using (`go version`)?
2. What operating system and processor architecture are you using?
3. What did you do?
4. What did you expect to see?
5. What did you see instead?

## Report a Bug

Open an issue. Please include descriptions of the following:
- Observations
- Expectations
- Steps to reproduce

## Contributing code

In general, this project follows Go project conventions, please read the [Contribution Guidelines](https://golang.org/doc/contribute.html) before sending patches.

## Contribute a Bug Fix

- Report the bug first
- Create a pull request for the fix

## Suggest a New Feature

- Create a new issue to start a discussion around new topic. Label the issue as `new-feature`

## Developer guidelines

### Download and Install

- Clone the project under `$GOPATH/src/github.com/ocibuilder/ocibuilder/`
- Run `dep ensure --vendor-only` or `dep ensure -v` to install package dependencies
- Build ctl binary using `make ocictl`, this will create `ocictl` cmd under `dist/`. Similarly, `make ocibuilder` will create package level binary.

### Run tests

For the entire package, follow this command
`go test -v <path-to-your-pkg>/`

For single function, add `-run <function-name>` after above command like this
`go test -v <path-to-your-pkg>/ -run <function-name>`

### Getting all the dependencies

```
$ make dep
```

### Re-generating Codegen

```
$ make dep
```

If you're making a change to the `pkg/apis` package, please ensure you re-run the K8 code-generator scripts found in the `/hack` folder. Ensure you have the `generate-groups.sh` script at the path: `vendor/k8s.io/code-generator/`. Next run the following command:

```
$ make codegen
```

### Re-generating OpenAPI

```
$ make dep
```

```
$ make openapigen
```

### Caveats

The docker and buildah spec template testing files are located under `/testing`.