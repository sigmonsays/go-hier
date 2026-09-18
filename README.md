a short experiement to build a hier like tool for deployment

the idea is to feed this tool configuration about the box which can be used to
lookup node specific data.

we'll have config like this
- inputs: can be auto discovered or configured
- hierarchy: lookup order

# limitations

- No support for merging lists yet

# example config

  inputs:
    fqdn: [ fqdn from grain ]
    domain: [ domain from grain ]
    env: [ env from grain ]
    pop: [ pop from grain ]

  hierarchy:
    - name: "by fqdn"
      path: "data/fqdn/{{ .inputs.fqdn }}.yaml"
    - name: "by domain"
      path: "data/domain/{{ .inputs.domain }}.yaml"
    - name: "by pop"
      path: "data/pop/{{ .inputs.pop }}.yaml"
    - name: "by env"
      path: "data/env/{{ .inputs.env }}.yaml"
    - name: "default"
      path: "data/default.yaml"

# data files

    data/fqdn/host1.site1.example.net.yaml
    ---
    exampled:
       pkg_version: 1.2
       feature1: true

    data/domain/site1.example.net.yaml
    ---
    exampled:
       pkg_version: 1.0

    data/env/production.yaml
    ---
    exampled:
       pkg_version: 1.0


# lookup function to use in salt

For example say we're on host1.site1.example.net, to get the package version we'd just do this

  hier lookup exampled.pkg_version

To show what variables are present on the current host we could do this

  hier show

To show what variables are present for arbitrary data

  hier show --inputs inputs.yaml

  inputs.yaml
  ---
  date: 2024-11-30
  inputs:
    fqdn: host2.site1.example.net
    domain: site1.example.net
    env: production
    pop: site1
