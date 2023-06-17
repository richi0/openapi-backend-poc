# OpenAPI Backend POC

## Goals

1. get an OpenAPI description document
2. use a code generator to create Go codebase on the description document
3. get input validation based on the description document
4. write the remaining code to create a mock backend service
5. generate the api documentation based on the description document

## Results

### Goal 1

The Petstore description document provided as an example by [OpenAPI](https://github.com/OAI/OpenAPI-Specification/blob/main/examples/v3.0/petstore-expanded.yaml) was used for this POC

### Goal 2

As a code generate [oapi-codegen](https://github.com/deepmap/oapi-codegen) was used.

The following command generated the code in the `generated` folder.

```shell
oapi-codegen -generate "server,types,spec"  -package generated petstore-expanded.yaml > generated/openapi.gen.go
```

### Goal 3

For validation, the middleware package from the [oapi-codegen project](https://github.com/deepmap/oapi-codegen/blob/master/pkg/middleware/oapi_validate.go) was used

```go
oapimiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
```

This middleware validates all input (path, query parameters, body, headers) based on the description document.

### Goal 4

Code in `main.go` was written to have a running server

### Goal 5

Api documentation was created by running the following command and is served protected by basic auth on path `/documentation`

```shell
npx @redocly/cli build-docs petstore-expanded.yaml
```
