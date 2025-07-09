# igor

A simple command-line tool to quickly fetch and copy `.gitignore` files from the official `github/gitignore` repository into your current project directory.

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