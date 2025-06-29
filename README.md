# git-branch-auto-remove

A CLI tool to automatically remove local Git branches that have been merged and deleted from the remote repository.

## Features

- Automatically identifies and lists local branches that are gone from the remote.
- Supports dry-run mode (default) to preview branches to be deleted.
- Allows forced deletion of branches.
- Supports deleting merged branches.
- Configurable protected branches (e.g., `main`, `master`, `develop`).
- Colorized output for better readability.

## Installation

To install `git-branch-auto-remove`, make sure you have Go installed (Go 1.16 or higher is recommended).

```bash
go install github.com/tkr53/git-branch-auto-remove@latest
```

This command will install the executable to your `$GOPATH/bin` (or `$GOBIN`) directory. Make sure this directory is in your system's `PATH`.

## Usage

Run the tool from within your Git repository:

```bash
git-branch-auto-remove
```

By default, this command will perform a dry-run, listing the branches that *would be* removed without actually deleting them. It will also prompt for confirmation before any deletion.

### Options

- `--force` or `-f`: Force execute deletion of branches without confirmation prompt.

  ```bash
  git-branch-auto-remove --force
  ```

- `--merged` or `-D`: Delete branches that have been merged into the current branch (similar to `git branch -d` for merged branches).

  ```bash
  git-branch-auto-remove --merged
  ```

### Configuration

You can configure `git-branch-auto-remove` by creating a `.git-branch-auto-remove.yml` file in your project's root directory or in your home directory.

Example `.git-branch-auto-remove.yml`:

```yaml
protected_branches:
  - main
  - master
  - develop
  - my-special-branch
```

The `protected_branches` list specifies branches that will never be automatically deleted by the tool.

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details. (Note: LICENSE file is not yet created.)
