# Scout Reference: Research Methodology & Source Guide

## Search Strategy Framework

### The SCOUT Method

**S** — Scope: Define the exact technology, version, and context
**C** — Current: Always search with current year and "latest"
**O** — Official: Start with official docs, repos, and release notes
**U** — User Evidence: Find production case studies and real-world usage
**T** — Triangulate: Cross-reference at least 3 independent sources

### Search Query Templates

```
# Best Practices Discovery
"[technology] best practices [year] production"
"[technology] production experience lessons learned"
"[framework] recommended project structure [year]"

# Version & Compatibility
"[library] latest version changelog"
"[library] [version] migration guide"
"[library] breaking changes [year]"

# Performance & Benchmarks
"[technology] performance benchmark [year]"
"[library] vs [alternative] benchmark comparison"
"[framework] optimization techniques production"

# Security
"[library] CVE security advisory"
"[technology] security best practices [year]"
"[framework] authentication authorization guide"

# Architecture & Patterns
"[pattern name] implementation [language] [year]"
"[architecture] real world example production"
"[technology] clean architecture hexagonal"
```

## Source Reliability Matrix

### Tier 1: Authoritative Sources
| Source | Use For | URL Pattern |
|--------|---------|-------------|
| Official Documentation | API reference, features | `*.{io,dev,org,com}/docs` |
| GitHub Releases | Version history, changelogs | `github.com/{owner}/{repo}/releases` |
| RFCs & Specs | Language/library standards | `rfc-editor.org`, `go.dev/ref/spec` |
| Core Team Blogs | Roadmap, design decisions | Official engineering blogs |

### Tier 2: High-Quality Community
| Source | Use For | Reliability Notes |
|--------|---------|-------------------|
| Engineering blogs (Uber, Stripe, Netflix) | Production patterns | Tested at scale |
| Go Blog, Vue Blog | Language/framework guidance | Semi-official |
| Conference talks (GopherCon, VueConf) | Cutting-edge patterns | Expert presenters |
| GitHub Discussions | Community solutions | Check accepted answers |

### Tier 3: Community Wisdom (Verify)
| Source | Use For | Caveat |
|--------|---------|--------|
| Stack Overflow | Specific problems | Check answer date & votes |
| Reddit (r/golang, r/vuejs) | Community sentiment | Subjective |
| Dev.to, Medium articles | Tutorials, walk-throughs | Quality varies widely |
| Personal blogs | Niche solutions | Verify independently |

## Research Depth Levels

### Level 1: Quick Check (2-3 minutes)
- Single WebSearch query
- Verify latest version
- Check for known deprecations
- **Use when**: Agent needs quick validation during code review or debugging

### Level 2: Standard Research (5-8 minutes)
- 2-3 WebSearch queries from different angles
- Read 2-3 key articles/docs
- Cross-reference findings
- **Use when**: Planning implementation, choosing between approaches

### Level 3: Deep Research (10-15 minutes)
- 4-5 WebSearch queries covering all angles
- Read official docs + changelogs
- Find production case studies
- Check security advisories
- Build comparison matrix
- **Use when**: Major technology decision, architecture change, migration planning

### Level 4: Comprehensive Audit (15-20 minutes)
- Full Level 3 research
- Benchmark analysis
- Community health check (GitHub stars, issue response time, release frequency)
- Security vulnerability scan
- Long-term viability assessment
- Migration path documentation
- **Use when**: Supreme workflow Phase 1, major feature development, technology adoption

## Technology Radar Categories

### Adopt (Use with confidence)
- Proven in production at scale
- Active maintenance and community
- Well-documented
- No known critical issues
- Aligns with project architecture

### Trial (Worth pursuing)
- Promising but needs more production validation
- Good community traction
- May have rough edges
- Good fit for non-critical paths

### Assess (Worth exploring)
- Interesting but unproven
- Early stage but from reputable source
- May solve a pain point better than current approach
- Needs proof of concept before adoption

### Hold (Do not adopt)
- Deprecated or superseded
- Known critical issues
- Better alternatives exist
- Does not fit project architecture

## Project-Specific Context (The Sentinel)

When researching for this project, always consider:

### Stack Constraints
- **Backend**: Go (latest stable), GORM, Chi router
- **Frontend**: Nuxt 3, Vue 3, TailwindCSS, Pinia
- **Database**: PostgreSQL
- **Deployment**: Docker, Docker Compose
- **Architecture**: Hexagonal (backend) + FSD (frontend)

### Research Priorities
1. Patterns that align with Hexagonal architecture
2. Go idiomatic solutions (not Java/Python translated)
3. Vue 3 Composition API (not Options API)
4. Docker-first deployment strategies
5. PostgreSQL optimization techniques

### Red Flags to Watch For
- Recommendations using `interface{}` without justification in Go
- Patterns that break Hexagonal layering
- Frontend patterns that don't follow FSD structure
- Solutions that require non-Docker deployment
- Libraries with GPL licensing for commercial use
