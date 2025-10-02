# Dockerfile in plugins/
# Optimized for granular Docker layer caching in a Go workspace monorepo.
# This Dockerfile ensures that changes in one plugin do not invalidate
# the build cache for other, unrelated plugins.

# --- Base Stage ---
# Defines the workspace directory and the PLUGIN_NAME build argument.
# Using a base stage keeps subsequent stages cleaner.
FROM golang:1.24-alpine AS base
WORKDIR /workspace
ARG PLUGIN_NAME

# --- Deps Stage ---
# This stage's sole purpose is to download Go module dependencies.
# It's structured to maximize layer caching.
FROM base AS deps

# 1. Copy the workspace definition. A change here invalidates all builds.
# The wildcard ensures this works even if go.work.sum doesn't exist yet.
COPY go.work* ./

# 2. Copy all source and module files for the entire workspace.
# This is necessary because 'go mod download' in a workspace context
# needs to resolve all modules.
COPY pkg/ ./pkg/
COPY cmd/ ./cmd/

# 3. Download all dependencies for the entire workspace at once.
RUN go mod download

# --- Builder Stage ---
# This stage compiles the specific plugin using the pre-downloaded dependencies.
FROM base AS builder

# 1. Copy the downloaded Go modules from the 'deps' stage.
COPY --from=deps /go/pkg/mod /go/pkg/mod

# 2. Copy the workspace definition again.
# The wildcard ensures this works even if go.work.sum doesn't exist yet.
COPY go.work* ./

# 3. Copy source code, being specific to minimize cache invalidation.
# Copy the shared 'pkg' source first, as it's a dependency for all plugins.
COPY pkg ./pkg
# Then, copy the source code for all plugins.
COPY cmd ./cmd

# 4. Build the target plugin.
# The command is run from the workspace root to ensure the go.work file is found and used.
RUN go build -o /app/server ./cmd/${PLUGIN_NAME}

# --- Final Stage ---
# Creates a minimal, non-root final image for the service.
FROM gcr.io/distroless/static-debian11

COPY --from=builder /app/server /server
USER nonroot:nonroot

ENTRYPOINT ["/server"]