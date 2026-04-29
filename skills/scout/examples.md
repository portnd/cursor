# Scout Examples: Research Output Templates

## Example 1: Quick Check — Library Validation

**Scenario**: Code agent wants to use `pgx` instead of `lib/pq` for PostgreSQL driver

```
## Scout Intel: PostgreSQL Driver for Go

**Recommendation**: Use pgx v5 as the primary PostgreSQL driver. lib/pq is in maintenance mode.
**Confidence**: 🟢 High
**Sources**:
  - github.com/lib/pq README: "Unless you have a specific reason to use this package, we recommend using pgx"
  - github.com/jackc/pgx — Active releases, v5 stable
**Key Insight**: lib/pq is officially recommending pgx for new projects. pgx v5 offers native PostgreSQL protocol, better performance, and supports all modern PostgreSQL features.
**Caveat**: pgx has a different API surface. If using GORM, ensure GORM's pgx driver is used (`gorm.io/driver/postgres` uses pgx under the hood since v1.25+).
```

---

## Example 2: Standard Research — Pattern Comparison

**Scenario**: Supreme workflow needs to evaluate WebSocket approach for real-time dashboard

```
## Scout Research Report: Real-time Dashboard Updates (Go + Nuxt 3)

### Executive Summary
WebSocket via gorilla/websocket remains the most battle-tested approach for Go, but the project is archived. The community is moving to `nhooyr.io/websocket` (now `github.com/coder/websocket`) or `gobwas/ws` for better performance. For Nuxt 3, the native `useWebSocket` composable from VueUse is the cleanest integration.

### Current Best Practices (as of 2026)
1. **Use `github.com/coder/websocket`** for Go server — official successor to gorilla/websocket, maintained by Coder
2. **Use `@vueuse/core` `useWebSocket`** for Nuxt 3 client — reactive, auto-reconnect, composable pattern
3. **Implement heartbeat/ping-pong** — prevent connection drops behind proxies
4. **Use Redis Pub/Sub** for multi-instance scaling — allows multiple Go instances to broadcast
5. **Binary protocol (MessagePack)** over JSON when payload is large — 30-50% bandwidth savings

### Technology Radar
| Technology | Status | Trend | Recommendation |
|-----------|--------|-------|----------------|
| coder/websocket | Adopt | ↑ Rising | Primary WebSocket library for Go |
| gorilla/websocket | Hold | ↓ Archived | No longer maintained |
| gobwas/ws | Trial | → Stable | High-performance alternative, lower-level API |
| Socket.IO (Go) | Hold | ↓ Declining | Over-engineered for Go backends |
| VueUse useWebSocket | Adopt | ↑ Rising | Best Nuxt 3 integration |

### Breaking Changes & Deprecations
- gorilla/websocket: Project archived Aug 2023, no further updates
- nhooyr.io/websocket: Migrated to github.com/coder/websocket, old import path deprecated

### Production Case Studies
1. **Coder**: Built their entire remote development platform on coder/websocket — handles thousands of concurrent terminal connections
2. **Liveblocks**: Uses WebSocket + Redis Pub/Sub for real-time collaboration at scale

### Recommended Approach for This Project
1. Backend: Use `github.com/coder/websocket` for WebSocket handler
2. Frontend: Use `@vueuse/core` `useWebSocket` composable
3. Scaling: Add Redis Pub/Sub when moving to multi-instance
4. Protocol: Start with JSON, migrate to MessagePack if bandwidth becomes an issue

### Confidence Level
🟢 High — Well-documented migration path, active community consensus
```

---

## Example 3: Deep Research — Technology Adoption Decision

**Scenario**: Evaluating whether to adopt a new ORM pattern or stick with current GORM setup

```
## Scout Research Report: Go ORM Landscape — GORM vs Alternatives (2026)

### Executive Summary
GORM remains the dominant Go ORM with the largest ecosystem, but Ent (by Facebook/Meta) has gained significant traction for complex schemas. sqlc is the preferred choice for SQL-first teams. For The Sentinel's current GORM usage, the recommendation is to stay with GORM but adopt the latest performance patterns rather than migrate.

### Current Best Practices (as of 2026)
1. **GORM v2 Performance**: Use `Session(&gorm.Session{PrepareStmt: true})` for prepared statement caching — up to 3x query speedup
2. **Batch Operations**: Use `CreateInBatches()` instead of bulk `Create()` — prevents memory spikes
3. **Context Propagation**: Always pass `ctx` to GORM operations for timeout/cancellation
4. **Soft Delete Awareness**: Be explicit about `Unscoped()` when soft deletes are configured
5. **Hook Optimization**: Minimize logic in GORM hooks — they run synchronously and block the query

### Technology Radar
| Technology | Status | Trend | Recommendation |
|-----------|--------|-------|----------------|
| GORM v2 | Adopt | → Stable | Continue using, optimize patterns |
| Ent (Meta) | Trial | ↑ Rising | Consider for new microservices with complex schemas |
| sqlc | Assess | ↑ Rising | Great for SQL-heavy, performance-critical queries |
| sqlx | Adopt | → Stable | Good for simple queries alongside GORM |
| sqlc + GORM hybrid | Trial | ↑ Rising | Use sqlc for read-heavy, GORM for write-heavy |

### Migration Risk Assessment
| Risk | GORM → Ent | GORM → sqlc | Stay + Optimize |
|------|-----------|-------------|-----------------|
| Effort | 🔴 High (full rewrite) | 🟡 Medium (query-by-query) | 🟢 Low (pattern changes) |
| ROI | 🟡 Medium (complex schemas) | 🟢 High (read performance) | 🟢 High (minimal effort) |
| Disruption | 🔴 High | 🟡 Medium | 🟢 Low |
| Learning curve | 🟡 Medium | 🟢 Low | 🟢 None |

### Recommended Approach for This Project
**Stay with GORM + optimize patterns.** The codebase already uses GORM extensively. Migration cost outweighs benefits at current scale. Focus on:
1. Enable prepared statement caching
2. Use `CreateInBatches` for bulk inserts
3. Add proper context propagation
4. Profile slow queries and add sqlc only for hot paths if needed

### Confidence Level
🟢 High — This is a mature space with clear community consensus
```

---

## Example 4: Compact Agent Feed Format

**Scenario**: Quick recommendation for another agent during code review

```
## Scout Intel: Go Error Wrapping Pattern

**Recommendation**: Use `fmt.Errorf("context: %w", err)` for wrapping. Use `errors.Is()` and `errors.As()` for checking. Avoid `pkg/errors` — it's archived.
**Confidence**: 🟢 High
**Sources**:
  - go.dev/blog/go1.13-errors (official)
  - github.com/pkg/errors README: archived notice
**Key Insight**: Since Go 1.13, the standard library has native error wrapping. pkg/errors is no longer needed.
**Caveat**: If the project already uses pkg/errors extensively, migrate incrementally — don't rewrite all at once.
```
