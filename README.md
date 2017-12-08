# sk8 [skate] - A quick way to deploy many microservices in a similar way

sk8 is a yaml-preprocessor for [kubectl][kctl]

sk8 uses simplified yaml-files describing a service and using templates to generate different [kubernetes][k8] resources.

## Usage
```
Usage: sk8 [option(s)] {file(s)}

Options:  -apply              Call 'kubectl' to update Kubernetes
          -kubeconfig | -kc   Show the current server according to the KUBECONFIG-file
          -verbose    | -v    Show verbose output
          -{template-tag}     Applies the template indicated by the {template-tag} key (i.e is prefixed with)
          -all                Applies all templates that are specified

Example:  sk8 simpleservice.yaml -all -apply
```

## How it works

sk8 works on folders in the current directory and merges any file you use with all the files in the '.sk8' folder.

You can also, per file, add designated 'parents' to include extra definitions into the service.

Se the [examples] for more information on how merging and specified inheritance works.

### Requirements

You must use the env-var named KUBECONFIG pointing to the kubeconfig-file you wish to use.


[k8]: https://github.com/kubernetes/kubernetes
[kctl]: https://github.com/kubernetes/kubectl