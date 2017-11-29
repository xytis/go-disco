# go-disco
Unified service discovery library

# Quickstart

Dependency management: [Dep](https://github.com/golang/dep)

    dep ensure
    go build

Test application accepts ENV variable DISCOVERY_BACKEND in format `backend:config_url`

Examples:

    DISCOVERY_BACKEND=dns:127.0.0.1:53?suffix=service.consul ./go-disco consul
    DISCOVERY_BACKEND=consul:default ./go-disco consul
