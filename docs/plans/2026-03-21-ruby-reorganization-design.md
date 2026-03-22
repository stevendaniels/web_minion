# Ruby Reorganization Design

**Date:** 2026-03-21

## Overview

Reorganize the repository so the Ruby gem implementation lives under `ruby/webminion/`, mirroring the existing `go-webminion/` subfolder. The root becomes a language-neutral landing page with a new README that introduces WebMinion, links to each implementation, and shows a brief spec example.

## File Destinations

| Source | Destination | Notes |
|---|---|---|
| `README.md` | `ruby/webminion/README.md` | Move via `git mv` |
| `Gemfile` | `ruby/webminion/Gemfile` | Move via `git mv` |
| `Gemfile.lock` | `ruby/webminion/Gemfile.lock` | Move via `git mv` |
| `Rakefile` | `ruby/webminion/Rakefile` | Move via `git mv` |
| `web_minion.gemspec` | `ruby/webminion/web_minion.gemspec` | Move via `git mv` |
| `lib/` | `ruby/webminion/lib/` | Move via `git mv` |
| `bin/` | `ruby/webminion/bin/` | Move via `git mv` |
| `test/` | `ruby/webminion/test/` | Move via `git mv` |
| `vendor/` | `ruby/webminion/vendor/` | Move via `git mv` |
| `LICENSE.txt` | stays at root + copy to `ruby/webminion/LICENSE.txt` | |
| `CODE_OF_CONDUCT.md` | stays at root | |
| `DESIGN.md` | stays at root | |
| `go-webminion/` | unchanged | |
| `docs/` | unchanged | |

A new `README.md` is created at the root (see below).

## New Root README Structure

1. **Header + one-paragraph description** — what WebMinion is: a metadata-driven browser automation library that lets you define browser bots via JSON configuration instead of writing custom code.

2. **Implementations table** — two rows:
   - Ruby — gem, supports Mechanize and Capybara/Selenium drivers → `[ruby/webminion](ruby/webminion)`
   - Go — net/http + Chrome via go-rod → `[go-webminion](go-webminion)`

3. **Spec example** — short annotated JSON snippet showing a minimal flow with one action and one step, followed by a "See [DESIGN.md](DESIGN.md) for the full schema" link.

4. **Footer** — links to LICENSE and CODE_OF_CONDUCT.

## Implementation Approach

Use `git mv` for all moves to preserve file history. LICENSE.txt is copied (not moved) so both the root and `ruby/webminion/` have it.

## Success Criteria

- `git log --follow ruby/webminion/lib/web_minion.rb` shows full history
- Root contains only: `README.md`, `DESIGN.md`, `CODE_OF_CONDUCT.md`, `LICENSE.txt`, `docs/`, `go-webminion/`, `ruby/`
- `ruby/webminion/` is a self-contained Ruby gem (Gemfile, gemspec, lib, bin, test, vendor)
- Root README renders correctly with working relative links
