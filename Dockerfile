# syntax=docker/dockerfile:1

# (A) Prepare base & build the Go binary
# Use the builder platform to cross-compile to the target platform.
FROM --platform=${BUILDPLATFORM} golang:1.22-bookworm AS builder
WORKDIR /transiter

RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to "/usr/bin"

FROM builder AS codegen

# (2) Next, we perform an optional step in which we re-generate all of the sqlc and
# proto code and validate that it matches what's in source control.
# The idea is to make sure the repo is internally consistent and that the Docker image
# we're building actually reflects what's in the .proto and .sql files.

# (2.1) Install all the code generation tools.
COPY justfile .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    just install-tools

# (2.2) Generate the gRPC and DB files and then move them because changes to the bind mount are
# not persisted beyond the RUN / to the build context.
COPY buf.gen.yaml .
COPY buf.lock .
COPY buf.yaml .
COPY api api
COPY sqlc.yaml .
COPY db db
COPY docs/src/api/api_docs_gen.go docs/src/api/api_docs_gen.go
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=.,target=/transiter,rw=true \
    just generate && \
    mkdir -p /out/internal /out/docs/src && \
    mv internal/gen /out/internal/gen && \
    mv docs/src/api /out/docs/src/api && \
    rm /out/docs/src/api/api_docs_gen_input.json

# (2.3) Diff the newly generated files with the ones in source control.
# If there are differences, this will fail
FROM codegen AS verify-codegen
RUN --mount=type=bind,source=./internal/gen,target=/in \
    diff --recursive /in /out/internal/gen

RUN --mount=type=bind,source=./docs/src/api,target=/in \
    diff --recursive /in /out/docs/src/api

FROM builder AS build

# (3) Build the binary.
# As soon as TARGETOS/TARGETARCH are defined, the build will differ
# for each target platform during cross-compilation.
ARG TRANSITER_VERSION
ARG TARGETOS
ARG TARGETARCH
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=.,target=/transiter \
    EXTRA_LDFLAGS='-w' EXTRA_GOFLAGS='-trimpath -o /out/transiter' \
    just build ${TRANSITER_VERSION}

# (B) Build the documentation
# This is not platform-dependendent, so can be done once on the native
# build platform and shared by all target platforms.
FROM --platform=${BUILDPLATFORM} python:3.9 AS docs-builder
WORKDIR /transiter
RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to "/usr/bin"
COPY justfile ./
COPY docs/requirements.txt docs/
RUN pip install -r docs/requirements.txt
COPY docs/mkdocs.yml docs/
COPY docs/src docs/src
RUN just docs


# (C) Pull in the Caddy binary as a dependency (for the target platform).
FROM caddy:2 AS caddy


# (D) Put it all together.
# We use this buildpack image because it already has SSL certificates installed
# No emulation is required because there are no RUN statements.
FROM buildpack-deps:bookworm-curl
COPY --link --from=caddy /usr/bin/caddy /usr/bin/
COPY --link --from=docs-builder /transiter/docs/gen /usr/share/doc/transiter
COPY --link --from=build /out/transiter /usr/bin/
ENTRYPOINT ["transiter"]
