---
name: concurrency
description: Solve, explain, review, and implement concurrent code, especially in Go. Use when Codex needs to work on goroutines, channels, mutexes, atomics, contexts, worker pools, pipelines, bounded parallelism, cancellation, race conditions, deadlocks, goroutine leaks, or interview-style concurrency exercises.
---

# Concurrency

## Overview

Use this skill to solve or explain concurrency tasks with an emphasis on Go interview problems and practical debugging. Start by framing correctness requirements, then choose the simplest concurrency primitive that satisfies them.

## Workflow

1. Classify the task.
   - Explanation: define the model, trade-offs, and failure modes before showing code.
   - Implementation: state assumptions explicitly, then write the smallest correct concurrent design.
   - Review/debugging: prioritize races, deadlocks, leaks, cancellation gaps, ordering bugs, and unbounded goroutine creation.
2. Extract the hard requirements.
   - Preserve input order or allow completion order.
   - Fail fast on first error or collect all errors.
   - Support cancellation or timeout.
   - Limit concurrency or allow one goroutine per unit of work.
   - Return partial results or all-or-nothing.
   - Coordinate shared state or stream values between stages.
3. Choose the primitive with the lowest complexity.
   - Prefer `sync.Mutex` or `sync.RWMutex` for shared mutable state.
   - Prefer channels for ownership transfer, work queues, or pipelines.
   - Prefer `sync/atomic` for simple counters, flags, or single-word state.
   - Prefer `errgroup.Group` with `SetLimit` for peer tasks that should fail fast together.
   - Prefer a worker pool for many homogeneous jobs with bounded concurrency.
4. Make cancellation behavior explicit.
   - Thread `context.Context` through every function that can block or start work.
   - Stop accepting new work after `ctx.Done()`.
   - Do not claim a running task can stop early unless the task itself observes `ctx`.
5. Verify the result.
   - Check ordering, error semantics, close ownership, and leak risk.
   - Run `go test -race` when tests exist.
   - Add tests for cancellation, first-error behavior, worker bounds, and empty input.

## Decision Rules

- Choose `mutex` over channels when multiple goroutines need to mutate or read the same in-memory structure and no queue semantics are needed.
- Choose channels over `mutex` when the program naturally passes work or results between stages.
- Choose `errgroup` over manual `WaitGroup` management when sibling goroutines share one lifecycle and should stop on the first failure.
- Choose a worker pool over spawning one goroutine per job when job count may be large or bounded parallelism is required.
- Choose an atomic-index worker pool over a task channel when jobs already live in a slice and workers can claim indices directly.

## Go-Specific Rules

- Close a channel only from the sending side that owns it.
- Do not close a channel to signal "no more receives" when multiple senders still exist.
- Avoid unbuffered channels unless rendezvous semantics are intended.
- Avoid mixing channels and locks unless there is a clear ownership boundary.
- Avoid `time.Sleep` for synchronization in production logic or tests.
- Treat `go test -race` as mandatory for non-trivial concurrent code.

## Communication Pattern

When answering, start with the concurrency model and the invariants:
- what can run in parallel
- what state is shared
- what establishes ordering
- what stops work
- what happens on error

Then give the implementation or review findings.

## References

- Read [go-patterns.md](./references/go-patterns.md) for standard Go patterns, code templates, and primitive-selection guidance.
- Read [review-checklist.md](./references/review-checklist.md) when reviewing or debugging concurrency code.
