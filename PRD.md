PRD: chaos — DAG-Based Multi-Directory Command Runner

1. Overview

chaos is a CLI tool for executing arbitrary commands across directories in a Git repository, ordered by declared dependencies.

It treats directories as components, builds a dependency graph (DAG) from config files, and executes a user-provided command in each component’s working directory.

Core philosophy:
	•	Minimal constraints
	•	Explicit dependencies
	•	Execution everywhere
	•	Failures do not block unrelated work
	•	Users are responsible for correctness

2. Goals
	•	Execute a command across many directories in dependency order
	•	Support both sequential and parallel execution
	•	Require minimal configuration (zero-config works)
	•	Be fast, predictable, and simple
	•	Integrate cleanly with CI and devcontainers

3. Non-Goals
	•	No task definition language (not a build system)
	•	No caching
	•	No incremental builds
	•	No sandboxing beyond allow/block lists
	•	No cross-repo orchestration

4. Core Concepts

4.1 Component

A component is:
	•	Any directory with a chaos.yaml config file
	•	Plus any directory referenced as a dependency

Components:
	•	Are identified by repo-root-relative paths
	•	Always exist within the Git repository root

Examples:

/frontend
/backend/api
/libs/utils
/


4.2 Config File (chaos.yaml)

Optional file located in a component directory.

If absent:
	•	dependencies = none
	•	allowlist = *
	•	blocklist = none

Schema

dependencies:
  - /path/to/dependency
allow:
  - "*"
block:
  - "rm -rf *"

Rules
	•	All fields optional
	•	Missing fields use defaults
	•	Paths must be absolute (/...)
	•	Paths must stay within repo root
	•	Blocklist overrides allowlist

4.3 Dependency Graph (DAG)
	•	Built from all discovered chaos.yaml files
	•	Includes:
	•	all config-defined components
	•	all referenced dependencies (even without configs)

Validation:
	•	Cycles → error
	•	Missing paths → error
	•	Non-directory paths → error
	•	Paths outside repo → error

5. CLI

5.1 Commands

chaos exec <cmd> [args...]
chaos list
chaos dag
chaos dry-run exec <cmd> [args...]

5.2 Command Semantics
	•	Everything after exec is passed as raw argv
	•	No shell interpretation
	•	Equivalent to executing directly in each directory

Example:

chaos exec npm test

Runs:

npm test

in each component directory

6. Execution Model

6.1 Modes

Mode	Behavior
Sequential (default)	One component at a time
Parallel (--parallel)	All ready components run concurrently


6.2 Ordering
	•	Dependencies always execute before dependents
	•	Independent components have no ordering guarantees
	•	Topological sorting is used

6.3 Failure Behavior
	•	A component failure does not stop execution
	•	All components are attempted
	•	Final result:
	•	success if all succeed
	•	failure if any fail

6.4 Dependency Failure Handling
	•	Dependents still execute even if dependencies fail
	•	DAG enforces ordering only, not correctness gating

6.5 Ctrl+C Behavior

Mode	Behavior
Sequential	Stop current process and exit
Parallel	Terminate all running processes and exit

Exit code on interrupt: 130

7. Output

7.1 Sequential
	•	Direct streaming of stdout/stderr
	•	No modification

7.2 Parallel
	•	Each line prefixed with component path:

/frontend | running tests...
/backend  | building...

	•	Output streamed in real-time

7.3 Partial Lines
	•	Buffered until newline or process exit

8. Command Filtering

8.1 Allow / Block Lists
	•	Match against full command (argv string)
	•	* matches all
	•	Blocklist overrides allowlist

8.2 Behavior
	•	If blocked:
	•	command is skipped
	•	optionally logged in verbose mode

8.3 Skip Types
	•	Policy skip (blocked command)
	•	Dependency failure skip (not enforced, but tracked)
	•	Interrupt skip

9. Discovery

9.1 Process
	1.	Locate Git repo root using Git
	2.	Find all chaos.yaml files
	3.	Build component set:
	•	config directories
	•	dependency targets

9.2 Path Rules
	•	Normalized repo-root-relative paths
	•	Format: /a/b/c
	•	No ..
	•	No absolute filesystem paths

9.3 Symlinks
	•	Normalized by default
	•	Optional flag to disable normalization

10. Commands

10.1 list

Outputs all components:
	•	path
	•	whether config exists

10.2 dag

Outputs dependency graph:
	•	edges or adjacency list
	•	optionally grouped by levels

10.3 dry-run

Shows:
	•	execution order
	•	skipped components
	•	allow/block decisions
	•	no commands executed

11. Environment
	•	Commands inherit parent environment unchanged
	•	Working directory = component directory

12. Exit Codes

Condition	Code
Success	0
Any failure	1
Interrupted	130
Config / DAG error	non-zero (e.g. 2)


13. Distribution

13.1 Build
	•	Language: Rust
	•	Single static binary (where possible)

13.2 Releases
	•	Managed via Cargo Dist
	•	Published to GitHub Releases

Artifacts:

chaos-<version>-<target>.tar.gz
checksums.txt

Example targets:
	•	x86_64-apple-darwin
	•	aarch64-apple-darwin
	•	x86_64-unknown-linux-gnu
	•	aarch64-unknown-linux-gnu

13.3 Install Script

Usage:

curl -sSL ... | sh

Behavior:
	•	Detect OS/arch
	•	Map to Cargo Dist target triple
	•	Download matching release artifact
	•	Verify checksum (if available)
	•	Install to:
	•	$HOME/.local/bin or
	•	/usr/local/bin if writable

13.4 Devcontainer Feature

Purpose:
	•	Install chaos inside other devcontainers

Behavior:
	•	Accept version parameter
	•	Download matching GitHub Release binary
	•	Install to /usr/local/bin/chaos
	•	No source builds

14. Versioning
	•	Semantic versioning
	•	Git tags: vX.Y.Z
	•	Version embedded at build time

Command:

chaos version


15. Future Considerations (Not in v1)
	•	Component filtering / targeting
	•	Reverse dependency execution
	•	Change detection via git diff
	•	Tags / labels
	•	Concurrency limits
	•	Structured output (JSON)

16. Summary

chaos is:
	•	a DAG-aware command runner
	•	for executing arbitrary commands across directories
	•	with minimal structure and maximum flexibility

It prioritizes:
	•	simplicity over safety
	•	execution over orchestration
	•	developer control over guardrails
