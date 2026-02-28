# Contributing to Engram Go SDK

Thank you for your interest in contributing to the Engram Go SDK.

## Development Setup

```bash
# Clone the repository
git clone https://github.com/Terronex-dev/engram-go.git
cd engram-go

# Download dependencies
go mod download

# Run tests
go test -v

# Run tests with coverage
go test -cover

# Format code
go fmt ./...

# Run linter
go vet ./...
```

## Code Style

- Follow standard Go formatting (`go fmt`)
- Pass `go vet` with no warnings
- Add doc comments for exported APIs
- Include tests for new functionality
- Keep functions focused and small

## Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Write tests for your changes
4. Ensure all tests pass (`go test -v`)
5. Run `go fmt` and `go vet`
6. Commit with clear messages
7. Push and open a PR

## Cross-SDK Compatibility

Changes must maintain compatibility with the TypeScript, Python, and Rust SDKs. Run the cross-SDK test suite to verify:

```bash
cd ../engram-test-suite
npm run test:ts
python3 scripts/run_python_tests.py
```

## Testing

Tests are in `engram_test.go`. Run with:

```bash
go test -v
go test -cover
go test -race  # Check for race conditions
```

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
