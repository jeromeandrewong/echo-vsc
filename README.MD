
# Echo

A CLI tool to convert your existing **VSCode themes** to **iTerm colors**.

## Demo

![demo](demo.gif)

## Installation

```bash
go install github.com/jeromeandrewong/echo-vsc/cmd/echo-vsc@latest
```

## Usage

```bash
echo-vsc
```

## Project Structure

```txt
echo-vsc/
├── cmd/
│   └── echo-vsc/
│       └── main.go
├── internal/
│   ├── constants/
│   │   └── constants.go
│   ├── logger/
│   │   └── logger.go
│   ├── theme/
│   │   └── theme.go
│   ├── vsc/
│   │   └── vsc.go
│   └── converter/
│       └── converter.go
├── pkg/
│   └── utils/
│       └── utils.go
└── go.mod
└── go.sum
└── README.MD
```

## Worker Pool Concurrency Pattern

The worker pool pattern helps prevent system overload while maintaining high performance through parallel processing.

- Controlled concurrency with fixed number of workers
- Efficient resource usage
- Better error handling and monitoring
- Consistent performance even with many extensions
>
## How it works

1. When user starts the program, we get all of the user's VSCode themes from the `~/.vscode/extensions` folder.

2. A worker pool of 5 concurrent workers is created to process the extensions:
    - Workers wait for jobs through a job channel
    - Each worker processes extensions by reading and unmarshalling package.json files
    - Results (themes and any errors) are sent through a results channel

3. Extensions are distributed as jobs:
    - Each extension directory is sent as a job to the worker pool
    - Workers pick up jobs as they finish their previous work
    - This ensures controlled, efficient processing

4. The main goroutine:
    - Coordinates the distribution of work
    - Collects results from all workers
    - Aggregates themes into the final collection
    - Handles any errors that occur during processing

5. we display the themes to the user in a list from [bubbletea's library](https://github.com/charmbracelet/bubbletea), just for aesthetic purposes (`vsctheme_picker.go`)

6. If the user-selected theme does not have a `themeType` set, we ask the user to select one either `light` or `dark`. (`theme_type_picker.go`)

7. After getting the file path of the selected theme, we convert it to an iTerm theme using the `convertTheme` function (`converter.go`)
    - read theme file
    - removes comments and trailing commas or unmarshalling errors will occur
    - unmarshal the cleaned file
    - iterate through ANSI color mappings and retrive corresponding color from vscode theme + add fallback colors if missing
    - convert hex color to RGBA (iterm uses RGBA)
    - generate iTerm theme XML
