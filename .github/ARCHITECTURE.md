# Architecture Guide

This guide explains how this full-stack codebase is organized.

## File Structure

| Folder          | Description                 |
| --------------- | --------------------------- |
| `cmd/`          | Command-line applications   |
| `ops/`          | Operations                  |
| `pkg/`          | Packages                    |
| `pkg/bliss/`    | Common go packages          |
| `tmp/`          | Temporary files             |

- All command-line applications and packages are placed in the `cmd` and `pkg` folders, respectively. They are
  language-independent.
