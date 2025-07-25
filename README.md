# Bucket Organizer

CR~~U~~D Rest APIs for storing objects per bucket using nested maps in-memory

## Overview

This repository allows storing object IDs grouped by bucket ID using a `map[string]map[string]bool` structure. It is intended for lightweight, in-memory use cases such as testing or prototyping.

## Example

```go
type InMemoryRepo struct {
    cache map[string]map[string]bool
}

func NewInMemoryRepo() *InMemoryRepo {
    return &InMemoryRepo{
        cache: make(map[string]map[string]bool),
    }
}

func (r *InMemoryRepo) InsertObject(ctx context.Context, bucketId, objectId string) error {
    if r.cache[bucketId] == nil {
        r.cache[bucketId] = make(map[string]bool)
    }
    r.cache[bucketId][objectId] = true
    return nil
}
```

## 🔍 Debugging Note

During development, a runtime panic was encountered:

```
panic: assignment to entry in nil map
```

This happened because the `cache` map itself was not initialized before being accessed:

```go
if r.cache[bucketId] == nil {
    r.cache[bucketId] = make(map[string]bool)
}
```

Although this checks for `nil` in the nested map, it assumes that `r.cache` is already initialized — which it wasn’t.

The fix involved introducing a constructor to initialize the `cache` map:

```go
func NewInMemoryRepo() *InMemoryRepo {
    return &InMemoryRepo{
        cache: make(map[string]map[string]bool),
    }
}
```

## 🤖 AI Assistance Acknowledgement

To identify and confirm the root cause of the panic, I used **ChatGPT**. The assistant helped by pointing out the distinction between an uninitialized top-level map (`nil`) and missing nested maps (`nil entry`), which led to the runtime error.

AI was used specifically for:

- Debugging support
- Validating assumptions
- Sanity-checking the logic

The final implementation and validation were done manually.
