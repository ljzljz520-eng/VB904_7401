# VideoLab

VideoLab is a pure-Go command-line project that plans and runs a deterministic
video grading workflow for seaside, city, and night footage. It extracts one
comparison frame per filter, selects a grading preset from fixed local scores,
saves the selected schemes in memory, and prints ffmpeg commands quoted for
Windows paths.

## Requirements

- Go 1.25.13
- `GOTOOLCHAIN=local`
- ffmpeg only when using the `execute` command

## Entry point

Run the deterministic planning workflow from the module root:

```sh
GOTOOLCHAIN=local go run ./cmd/videolab plan
```

Run the generated commands through a locally installed ffmpeg:

```sh
GOTOOLCHAIN=local go run ./cmd/videolab execute
```

Exercise the built-in deterministic decode-failure scenario without invoking
an external process:

```sh
GOTOOLCHAIN=local go run ./cmd/videolab simulate-error
```

The application writes a JSON page model to standard output. `execute` and
`simulate-error` return a non-zero status when command execution fails.

## Tests

Run the complete business workflow from the module root:

```sh
GOTOOLCHAIN=local go test -count=1 ./...
```

The repository uses only the Go standard library, embedded local fixtures,
in-memory storage, and synchronous execution.
