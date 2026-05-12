# Copilot Instructions for toolbox

## Project Overview

- This repository is a personal toolbox that integrates multiple small utilities.
- Main focus areas are social integrations for Mastodon and Bluesky.
- Keep changes practical, incremental, and easy to review.

## Coding Guidelines

- Language: Go.
- Prefer small, focused changes over broad refactors.
- Preserve existing public behavior unless explicitly asked to change it.
- Keep source-code comments in English.
- Follow existing package boundaries and naming style in this repository.

## Validation

- Use Taskfile tasks for local validation.
- Primary check command:

```sh
task test
```

- If needed, run additional checks:

```sh
task govulncheck
```

## Dependencies

- Do not add new external dependencies unless there is a clear benefit.
- If a dependency is added, explain why in the change summary.

## Documentation and Maintenance

- Update related documentation when behavior, options, or workflows change.
- Keep README and package comments consistent with the current implementation.
- Prefer explicit notes about assumptions and constraints when they are not obvious from code.

## Current Maintenance Context

- Ongoing maintenance includes workflow cleanup, potential integration with `github.com/goark/webinfo`, and documentation refresh.
- During this phase, prioritize readability and risk reduction over feature expansion.
