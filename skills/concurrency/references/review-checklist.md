# Concurrency Review Checklist

## Use This Reference

Read this file when reviewing or debugging concurrent Go code. Lead with correctness bugs, then note performance trade-offs.

## Findings Order

1. Data races
2. Deadlocks and blocked goroutines
3. Goroutine leaks
4. Incorrect cancellation or timeout handling
5. Ordering and visibility bugs
6. Unbounded concurrency or memory growth
7. Over-complicated design relative to requirements

## Review Questions

- What state is shared, and what synchronizes access to it?
- Who owns each channel, and who closes it?
- Can any send or receive block forever?
- What happens if one task returns an error?
- What happens if the context is canceled before or during execution?
- Can workers keep starting new jobs after failure or cancellation?
- Is output order required, and is it actually preserved?
- Does any goroutine outlive the request without intent?

## Common Bug Patterns

- Loop variable captured incorrectly by goroutines.
- Writing results into a shared slice or map without synchronization or index ownership.
- Closing a shared channel from multiple goroutines.
- Returning before background goroutines finish, leaving leaked work.
- Using buffered channels to hide coordination bugs.
- Holding a mutex while performing blocking I/O or calling user code.
- Assuming `context.CancelFunc` can stop work that never checks the context.

## What Good Looks Like

- Clear ownership of shared state and channels.
- Explicit error policy: fail fast or collect all.
- Explicit cancellation path.
- Bounded concurrency when input size can grow.
- Tests that exercise races, cancellation, and ordering.

## Review Output Style

For review responses:
- name the invariant that is broken
- point to the exact code path
- state the user-visible consequence
- suggest the simplest correction
