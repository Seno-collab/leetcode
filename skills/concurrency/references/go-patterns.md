# Go Concurrency Patterns

## Use This Reference

Read this file when the task requires choosing or implementing a concurrency pattern in Go. Prefer the simplest pattern that preserves the required invariants.

## Primitive Selection

- Shared map, cache, or accumulator: start with `sync.Mutex`.
- Mostly reads, rare writes: consider `sync.RWMutex`, but only if contention is real.
- Simple counter or flag: use `sync/atomic`.
- Parallel peer tasks with shared cancellation: use `errgroup.Group`.
- Bounded processing of many jobs: use a worker pool.
- Multi-stage processing: use a pipeline with clearly owned channels.

## Pattern: Worker Pool

Use when many homogeneous jobs must run with bounded concurrency.

Questions to settle first:
- Must results preserve input order?
- Should the first error stop remaining work?
- Can running work be canceled, or only future work?

Two common variants:
- Channel queue: best when producers and consumers are decoupled.
- Atomic index over a slice: best when all jobs already exist in memory.

Guidelines:
- Cap worker count at `min(workers, len(jobs))`.
- Store results by index when output order matters.
- Cancel on the first error if the contract is fail-fast.
- Do not claim full cancellation if `fn` lacks `context.Context`.

## Pattern: errgroup for Peer Tasks

Use when several tasks start together and any failure should stop the group.

Typical fit:
- fetch N independent resources
- run validation checks in parallel
- fan out to several independent computations

Guidelines:
- Derive the group from the request context.
- Use `SetLimit` to bound concurrency when task count is large.
- Return the first meaningful error with enough context to identify the failing task.

## Pattern: Mutex-Protected Shared State

Use when goroutines collaborate on one shared structure and queue semantics add needless complexity.

Good fits:
- concurrent cache
- metrics aggregation
- dedup tables
- in-memory session state

Guidelines:
- Keep critical sections short.
- Hold the lock only while touching protected state.
- Avoid calling slow or blocking functions while holding the lock.
- Document which fields the mutex protects.

## Pattern: Pipeline

Use when data naturally moves through stages such as parse -> validate -> transform -> write.

Guidelines:
- Give each stage one clear responsibility.
- Decide who owns closing each outbound channel.
- Check `ctx.Done()` in every stage that can block.
- Drain or stop upstream cleanly on failure to avoid leaks.

## Pattern: Fan-Out / Fan-In

Use when one input stream is processed by multiple workers and merged later.

Guidelines:
- Make ownership of the merge channel explicit.
- Close the output only after all senders are done.
- Beware of slow consumers causing blocked senders.

## Testing Guidance

Write tests for:
- empty input
- one item
- worker limit greater than job count
- deterministic ordering when required
- first error vs collect-all behavior
- canceled context before start
- canceled context during processing

Run:

```bash
go test -race ./...
```

## Red Flags

- goroutine spawned in a loop with no bound and no join path
- send on a channel that may never be received
- receive from a channel that nobody closes or writes
- multiple goroutines closing the same channel
- using `time.Sleep` to "wait" for correctness
- shared variable mutated without synchronization
