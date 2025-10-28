# Testing Guide for mcp-nomad

This document describes the testing infrastructure and how to run tests for the mcp-nomad project.

## Test Structure

The test suite is organized into several categories:

```
test/
├── unit/                    # Unit tests
│   ├── nomad_client_test.go
│   ├── tools_test.go
│   └── benchmark_test.go
├── integration/             # Integration tests
│   ├── nomad_client_integration_test.go
│   └── real_nomad_test.go
├── mocks/                   # Mock implementations
│   └── nomad_client_mock.go
├── testdata/                # Test data and fixtures
│   └── sample_data.go
├── config.go                # Test configuration
└── utils.go                 # Test utilities
```

## Test Categories

### Unit Tests
- **Location**: `test/unit/`
- **Purpose**: Test individual components in isolation
- **Mocking**: Uses mock implementations to isolate units under test
- **Speed**: Fast execution (< 1 second)
- **Dependencies**: No external dependencies

### Integration Tests
- **Location**: `test/integration/`
- **Purpose**: Test component interactions and API contracts
- **Mocking**: Uses mock HTTP servers to simulate Nomad API
- **Speed**: Medium execution (1-5 seconds)
- **Dependencies**: Minimal external dependencies

### Real Integration Tests
- **Location**: `test/integration/real_nomad_test.go`
- **Purpose**: Test against actual Nomad server
- **Mocking**: None - uses real Nomad API
- **Speed**: Slow execution (5-30 seconds)
- **Dependencies**: Requires running Nomad server

## Running Tests

### Using Makefile (Recommended)

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race

# Run benchmark tests
make test-benchmark

# Run all test types
make test-all

# Clean test artifacts
make clean-test
```

### Using Go Commands

```bash
# Run unit tests
go test -v ./test/unit/...

# Run integration tests
go test -v ./test/integration/...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./test/...
go tool cover -html=coverage.out -o coverage.html

# Run tests with race detection
go test -v -race ./test/...

# Run benchmark tests
go test -v -bench=. ./test/...
```

## Test Configuration

### Environment Variables

- `NOMAD_ADDR`: Nomad server address (default: http://localhost:4646)
- `NOMAD_TOKEN`: Nomad ACL token (optional)
- `SKIP_INTEGRATION`: Skip integration tests (default: false)

### Test Data

Test data is centralized in `test/testdata/sample_data.go` and includes:
- Sample job specifications (HCL and JSON)
- Sample jobs, nodes, namespaces, allocations
- Sample variables, ACL tokens, logs
- Sample cluster data

## Mock Implementation

The mock implementation (`test/mocks/nomad_client_mock.go`) provides:
- Complete interface implementation
- Configurable behavior via function pointers
- Support for all NomadClient methods
- Easy setup for different test scenarios

### Example Usage

```go
mockClient := &mocks.MockNomadClient{}
mockClient.ListJobsFunc = func(namespace, status string) ([]types.JobSummary, error) {
    return testdata.SampleJobs, nil
}

jobs, err := mockClient.ListJobs("default", "")
// Test assertions...
```

## Test Utilities

The `test/utils.go` file provides:
- Assertion helpers for common types
- Test data creation functions
- JSON comparison utilities
- Log format validation

### Example Usage

```go
// Assert job equality
test.AssertJobEqual(t, expectedJob, actualJob)

// Create test data
job := test.CreateTestJob("test-id", "test-name", "service")

// Assert JSON equality
test.AssertJSONEqual(t, expectedJSON, actualJSON)
```

## Continuous Integration

Tests are automatically run on:
- Push to main/develop branches
- Pull requests
- Multiple Go versions (1.21, 1.22, 1.23)

### GitHub Actions Workflow

The `.github/workflows/test.yml` file defines:
- Unit and integration tests
- Race detection
- Coverage reporting
- Benchmark tests
- Linting
- Security scanning

## Coverage Goals

- **Unit Tests**: > 80% coverage
- **Integration Tests**: > 60% coverage
- **Overall**: > 70% coverage

## Best Practices

### Writing Tests

1. **Use descriptive test names** that explain the scenario
2. **Follow the Arrange-Act-Assert pattern**
3. **Test both success and failure cases**
4. **Use table-driven tests for multiple scenarios**
5. **Mock external dependencies**
6. **Keep tests focused and atomic**

### Example Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name           string
        input          InputType
        expectedOutput OutputType
        expectedError  string
    }{
        {
            name:           "successful case",
            input:          validInput,
            expectedOutput: expectedResult,
            expectedError:  "",
        },
        {
            name:           "error case",
            input:          invalidInput,
            expectedOutput: nil,
            expectedError:  "expected error message",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockClient := setupMock()
            
            // Act
            result, err := functionUnderTest(tt.input)
            
            // Assert
            if tt.expectedError != "" {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
            } else {
                require.NoError(t, err)
                assert.Equal(t, tt.expectedOutput, result)
            }
        })
    }
}
```

## Debugging Tests

### Verbose Output

```bash
go test -v ./test/unit/...
```

### Run Specific Test

```bash
go test -v -run TestSpecificFunction ./test/unit/...
```

### Debug Mode

```bash
go test -v -race -count=1 ./test/unit/...
```

## Performance Testing

### Benchmark Tests

```bash
# Run all benchmarks
make test-benchmark

# Run specific benchmark
go test -v -bench=BenchmarkListJobs ./test/unit/...

# Run benchmark with memory profiling
go test -v -bench=BenchmarkListJobs -benchmem ./test/unit/...
```

### Performance Goals

- **API Calls**: < 100ms average
- **Memory Usage**: < 10MB per operation
- **Concurrent Requests**: Support 100+ concurrent requests

## Troubleshooting

### Common Issues

1. **Test fails with "connection refused"**
   - Check if Nomad server is running
   - Verify NOMAD_ADDR environment variable

2. **Test fails with "permission denied"**
   - Check NOMAD_TOKEN environment variable
   - Verify ACL permissions

3. **Race condition detected**
   - Review concurrent access patterns
   - Use proper synchronization

4. **Test timeout**
   - Check network connectivity
   - Increase timeout values if needed

### Getting Help

- Check the test logs for detailed error messages
- Review the GitHub Actions workflow logs
- Consult the Nomad API documentation
- Check the project's issue tracker
