# Ruby Reorganization Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers-executing-plans to implement this plan task-by-task.

**Goal:** Move the Ruby gem implementation into `ruby/webminion/` and replace the root `README.md` with a repo-level landing page.

**Architecture:** All Ruby files are moved via `git mv` (preserves history). `LICENSE.txt` is copied to both locations. A new root README introduces the project and links to each language implementation.

**Tech Stack:** git mv, bash cp, Markdown

---

### Task 1: Move top-level Ruby files into `ruby/webminion/`

**Files:**
- Move: `Gemfile` → `ruby/webminion/Gemfile`
- Move: `Gemfile.lock` → `ruby/webminion/Gemfile.lock`
- Move: `Rakefile` → `ruby/webminion/Rakefile`
- Move: `web_minion.gemspec` → `ruby/webminion/web_minion.gemspec`
- Move: `README.md` → `ruby/webminion/README.md`

**Step 1: Create the destination directory**

```bash
mkdir -p ruby/webminion
```

**Step 2: Move files with git mv**

```bash
git -C /Users/stevendaniels/src/stevendaniels/web_minion mv Gemfile ruby/webminion/Gemfile
git -C /Users/stevendaniels/src/stevendaniels/web_minion mv Gemfile.lock ruby/webminion/Gemfile.lock
git -C /Users/stevendaniels/src/stevendaniels/web_minion mv Rakefile ruby/webminion/Rakefile
git -C /Users/stevendaniels/src/stevendaniels/web_minion mv web_minion.gemspec ruby/webminion/web_minion.gemspec
git -C /Users/stevendaniels/src/stevendaniels/web_minion mv README.md ruby/webminion/README.md
```

**Step 3: Verify staging**

```bash
git -C /Users/stevendaniels/src/stevendaniels/web_minion status
```

Expected: all five files shown as `renamed` in the index.

---

### Task 2: Move `lib/`, `bin/`, `test/`, `vendor/` into `ruby/webminion/`

**Files:**
- Move: `lib/` → `ruby/webminion/lib/`
- Move: `bin/` → `ruby/webminion/bin/`
- Move: `test/` → `ruby/webminion/test/`
- Move: `vendor/` → `ruby/webminion/vendor/`

**Step 1: Move each directory**

```bash
git -C /Users/stevendaniels/src/stevendaniels/web_minion mv lib ruby/webminion/lib
git -C /Users/stevendaniels/src/stevendaniels/web_minion mv bin ruby/webminion/bin
git -C /Users/stevendaniels/src/stevendaniels/web_minion mv test ruby/webminion/test
git -C /Users/stevendaniels/src/stevendaniels/web_minion mv vendor ruby/webminion/vendor
```

**Step 2: Verify history is preserved on a representative file**

```bash
git -C /Users/stevendaniels/src/stevendaniels/web_minion log --follow --oneline ruby/webminion/lib/web_minion.rb
```

Expected: shows at least one commit (the original creation).

**Step 3: Spot-check the destination tree**

```bash
ls ruby/webminion/
```

Expected output contains: `bin  Gemfile  Gemfile.lock  lib  Rakefile  test  vendor  web_minion.gemspec  README.md`

---

### Task 3: Copy `LICENSE.txt` into `ruby/webminion/`

LICENSE stays at the root (repo-level) and also lives inside the gem directory.

**Step 1: Copy the file**

```bash
cp /Users/stevendaniels/src/stevendaniels/web_minion/LICENSE.txt \
   /Users/stevendaniels/src/stevendaniels/web_minion/ruby/webminion/LICENSE.txt
```

**Step 2: Stage the copy**

```bash
git -C /Users/stevendaniels/src/stevendaniels/web_minion add ruby/webminion/LICENSE.txt
```

**Step 3: Verify root LICENSE.txt is still present**

```bash
ls /Users/stevendaniels/src/stevendaniels/web_minion/LICENSE.txt
```

Expected: file exists (not moved, just copied).

---

### Task 4: Commit the file moves

**Step 1: Review staged changes**

```bash
git -C /Users/stevendaniels/src/stevendaniels/web_minion status
```

Expected: all Ruby files shown as renamed/added under `ruby/webminion/`, `LICENSE.txt` unchanged at root.

**Step 2: Commit**

```bash
git -C /Users/stevendaniels/src/stevendaniels/web_minion -c commit.gpgsign=false commit -m "refactor: move Ruby gem into ruby/webminion/"
```

---

### Task 5: Create new root `README.md`

**Files:**
- Create: `README.md`

**Step 1: Write the file**

Create `/Users/stevendaniels/src/stevendaniels/web_minion/README.md` with this content:

````markdown
# WebMinion

WebMinion is a metadata-driven browser automation library. Rather than writing a custom bot in code, you describe a browser automation flow in a JSON configuration file and hand it to WebMinion to execute. Flows are composed of named actions and steps — each step targets a page element, performs a method (navigate, fill, click, validate), and passes control to the next action.

## Implementations

| Language | Description | Directory |
|---|---|---|
| Ruby | Gem supporting Mechanize and Capybara/Selenium drivers | [ruby/webminion](ruby/webminion) |
| Go | net/http + Chrome (via go-rod) | [go-webminion](go-webminion) |

## Spec Example

A minimal flow that navigates to a URL and validates the resulting page:

```json
{
  "config": {
    "driver": "mechanize"
  },
  "flow": {
    "name": "Login Check",
    "actions": [
      {
        "name": "Go to login page",
        "starting": true,
        "steps": [
          {
            "name": "Navigate",
            "target": "https://example.com/login",
            "method": "go",
            "is_validator": false
          },
          {
            "name": "Confirm login page loaded",
            "target": "//h1[text()='Sign in']",
            "method": "element_exists",
            "is_validator": true
          }
        ]
      }
    ]
  }
}
```

See [DESIGN.md](DESIGN.md) for the full schema reference.

## License

[MIT](LICENSE.txt) · [Code of Conduct](CODE_OF_CONDUCT.md)
````

**Step 2: Stage and commit**

```bash
git -C /Users/stevendaniels/src/stevendaniels/web_minion add README.md
git -C /Users/stevendaniels/src/stevendaniels/web_minion -c commit.gpgsign=false commit -m "docs: add repo-level README"
```

**Step 3: Verify root directory is clean**

```bash
ls /Users/stevendaniels/src/stevendaniels/web_minion
```

Expected root contains only: `CODE_OF_CONDUCT.md  DESIGN.md  docs  go-webminion  LICENSE.txt  README.md  ruby`

---

### Task 6: Smoke-test the Ruby gem from its new location

This confirms nothing broke in the move — no hardcoded root-relative paths in the Ruby source.

**Step 1: Check that `test_helper.rb` paths still resolve correctly**

The file uses `./test/test_json/` and `./test/test_html/` — these are relative to wherever `rake` is invoked. After the move, `rake` must be run from `ruby/webminion/`.

```bash
cd /Users/stevendaniels/src/stevendaniels/web_minion/ruby/webminion && bundle exec rake test 2>&1 | tail -20
```

Expected: test suite runs (pass or pre-existing failures — no "No such file" load errors).

**Step 2: Verify gemspec still resolves `lib/web_minion/version.rb`**

The gemspec uses `File.expand_path("../lib", __FILE__)` — `__FILE__` is the gemspec path, so this remains correct regardless of where the gemspec lives.

```bash
cd /Users/stevendaniels/src/stevendaniels/web_minion/ruby/webminion && bundle exec ruby -e "require 'web_minion'; puts WebMinion::VERSION"
```

Expected: prints the version string (e.g. `0.x.x`).
