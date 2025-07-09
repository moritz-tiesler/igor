# igor

A simple command-line tool to quickly fetch and copy `.gitignore` files from the official [**github/gitignore**](https://github.com/github/gitignore) repository into your current project directory.

##  Installation

### 1. Pre-compiled Binaries 

1.  Go to the [**Releases page**](https://github.com/moritz-tiesler/igor/releases).
2.  Download the appropriate `.zip` or `.tar.gz` file for your system.

### 2. From Source using `go install`

If you have Go (1.18 or higher recommended) installed and prefer to build from source, you can use `go install`:

```bash
go install github.com/moritz-tiesler/igor@latest
```

### Copy a `.gitignore` file

To copy a `.gitignore` file for a specific language or framework, simply provide its name as an argument:

```bash
igor Go
# Copies Go.gitignore to .gitignore in your current directory

igor Python
# Copies Python.gitignore to .gitignore

igor Node
# Copies Node.gitignore to .gitignore
```
### List all available `.gitignore` files

If you're unsure of the exact name or want to see all options, use the `--list` flag:

```bash
igor --list
```
