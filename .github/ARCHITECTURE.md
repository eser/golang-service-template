# Architecture Guide

This guide explains how this full-stack codebase is organized.

## File Structure

| Folder          | Description                 |
| --------------- | --------------------------- |
| `cmd/`          | Command-line applications   |
| `docs/`         | Documentation               |
| `ops/`          | Operations                  |
| `pkg/`          | Packages                    |
| `pkg/bliss/`    | Common go packages          |
| `pkg/proto-go/` | Go protocol buffers         |
| `pkg/proto-ts/` | TypeScript protocol buffers |
| `pkg/web/`      | Web site                    |
| `specs/`        | Specification files         |
| `tmp/`          | Temporary files             |

- All command-line applications and packages are placed in the `cmd` and `pkg` folders, respectively. They are
  language-independent.
