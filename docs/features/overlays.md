## Overlays

The ocibuilder supports yaml overlays which can be applied at runtime using the `--overlay` command line flag and passing in an overlay yaml file.

The overlaying functionality of ocibuilder is built on top of [ytt](https://github.com/k14s/ytt) which makes use of annotations for matching and overlaying yaml on top of a template. 

The ocibuilder abstracts away these annotations and takes advantage of labels which can be defined in your build specification. You specify what array item you want to overlay
by adding an `overlay: <NAME>` label to your specification file. The ocibuilder will then match this label with the same label in your overlay file and apply
all the necessary annotations for you.


>**NOTE:** If you want to specify your own ytt annotations, you are able to do so by passing a [standard annotated ytt overlay file](https://get-ytt.io/#example:example-overlay-files). Annotation
reference for ytt can be found [here](https://github.com/k14s/ytt/blob/master/docs/lang-ref-ytt-overlay.md).


Example:

*spec.yaml*
```yaml
build:
  templates:
    - name: template-1
      cmd:
        - docker:
            inline:
              - ADD . /src
              - RUN cd /src && go build -o goapp
  steps:
    - metadata:
        name: go-build
        labels:
          overlay: first-step
      stages:
        - metadata:
            name: build-env
          base:
            image: golang
            platform: alpine
          template: template-1
      tag: v0.1.0
```

*overlay.yaml*
```yaml
build:
  steps:
    - metadata:
        labels:
          # required overlay label to refer to specific build step
          overlay: first-step
      # the value which we want to override in our template
      tag: v0.2.0
```

*overlay applied*
```yaml
build:
  templates:
    - name: template-1
      cmd:
        - docker:
            inline:
              - ADD . /src
              - RUN cd /src && go build -o goapp
  steps:
    - metadata:
        name: go-build
        labels:
          overlay: first-step
      stages:
        - metadata:
            name: build-env
          base:
            image: golang
            platform: alpine
          template: template-1
      # the new tag value which has been overlayed onto the spec
      tag: v0.2.0
```

