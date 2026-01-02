# Go Backend Implementation Plan

## Goal Description
Professionalize the `goBackend` project by adding critical missing components such as testing, environment configuration, API documentation, and linting. The current codebase lacks these standard engineering practices, which are essential for maintainability and reliability.

## Missing Components Analysis
The following components are currently missing or need improvement:

1.  **Testing Strategy**:
    *   **Current State**: `Makefile` has a `test` target, but there are no `*_test.go` files in the codebase.
    *   **Missing**: Unit tests for services (starting with `auth-service`), Integration tests for API endpoints.
2.  **Configuration Management**:
    *   **Current State**: Database credentials and ports are hardcoded in `docker-compose.yml` and `Makefile`.
    *   **Missing**: `.env` file support, `.env.example` template, and robust config loading in Go code (e.g., using `viper` or `godotenv`).
3.  **API Documentation**:
    *   **Current State**: Only `README.md` provides a text-based summary of endpoints.
    *   **Missing**: Standard Swagger/OpenAPI specification (v3) for the `bff-gateway`.
4.  **Code Quality & Linting**:
    *   **Current State**: No linting configuration.
    *   **Missing**: `golangci-lint` integration in `Makefile` and CI pipeline.
5.  **Observability**:
    *   **Current State**: Basic logging.
    *   **Missing**: Structured logging (e.g., `zap`) and distributed tracing (e.g., OpenTelemetry).

## Proposed Changes

### 1. Configuration Management
#### [NEW] [.env.example](file:///C:/Users/user/Documents/antiGoogle/goBackend/.env.example)
*   Create a template for environment variables.
#### [MODIFY] [docker-compose.yml](file:///C:/Users/user/Documents/antiGoogle/goBackend/docker-compose.yml)
*   Update to use `${VAR}` syntax for secrets and ports.
#### [MODIFY] [Makefile](file:///C:/Users/user/Documents/antiGoogle/goBackend/Makefile)
*   Update to load `.env` variables if present.

### 2. Testing Framework
#### [NEW] [services/auth-service/internal/usecase/auth_usecase_test.go](file:///C:/Users/user/Documents/antiGoogle/goBackend/services/auth-service/internal/usecase/auth_usecase_test.go)
*   Implement unit tests for `Login` and `Register` use cases using `testify`.
#### [NEW] [test/integration/auth_test.go](file:///C:/Users/user/Documents/antiGoogle/goBackend/test/integration/auth_test.go)
*   Add integration tests hitting the `bff-gateway` (running in Docker).

### 3. API Documentation
#### [NEW] [bff-gateway/api/openapi.yaml](file:///C:/Users/user/Documents/antiGoogle/goBackend/bff-gateway/api/openapi.yaml)
*   Create OpenAPI v3 specification describing all BFF endpoints.
#### [MODIFY] [README.md](file:///C:/Users/user/Documents/antiGoogle/goBackend/README.md)
*   Link to the OpenAPI spec and explain how to view it (e.g., using Swagger UI).

### 4. Tooling
#### [MODIFY] [Makefile](file:///C:/Users/user/Documents/antiGoogle/goBackend/Makefile)
*   Add `lint` target using `golangci-lint`.
*   Add `swagger` target to serve or generate docs.

## Verification Plan

### Automated Tests
*   Run `make test` to execute the newly added unit tests.
    *   Expected output: `PASS` for all packages.
*   Run `make lint` to verify code quality.
    *   Expected output: No linting errors.

### Manual Verification
1.  **Config**: Create `.env` from `.env.example`, run `docker-compose up`, and verify services start correctly.
2.  **Docs**: Open the Swagger UI (if added) or raw YAML to verify API definitions match the `README.md`.
