 Block Path

[![Build Status](https://github.com/traefik/plugin-blockpath/workflows/Main/badge.svg?branch=master)](https://github.com/traefik/plugin-blockpath/actions)

Block Path is a middleware plugin for [Traefik](https://github.com/traefik/traefik) which sends an HTTP `403 Forbidden` 
response when the requested HTTP path matches one the configured [regular expressions](https://github.com/google/re2/wiki/Syntax).

## Configuration

## Static

```toml
[pilot]
    token="xxx"

[experimental.plugins.blockpath]
    modulename = "github.com/samson4649/plugin-blockpath"
    version = "v0.2.1"
```

## Dynamic

To configure the `Block Path` plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/). The following example creates and uses the `blockpath` middleware plugin to block all HTTP requests with a path starting with `/foo` with a `404` and `/swagger` with a `401`. 

```yaml
http:
  middlewares:
    blockpath:
      plugin:
        blockpath:
          elements:
            - regex: ^/swagger
              code: 401
            - regex: ^/foo
              code: 404
```
