# gitrans - Git Transform

A tool for maintaining downstream forks by applying automated patches to upstream versions, eliminating the need to directly modify upstream code.

## Overview

gitrans enables you to maintain a clean separation between upstream code and your customizations. Instead of making direct modifications that lead to merge conflicts, you define transformation scripts that automatically generate your version from any upstream commit.

**Formula:** `UpstreamCommit + PatchScript => OurCommit`

## Key Concepts

- **Upstream**: The original repository you're forking from
- **Target Branch**: The branch tracking upstream (usually `main`)
- **Working Branch**: Your customized branch (e.g., `foo`) where patches are applied
- **PatchScript**: Transformation script `git_patch.gox` written in XGo
- **Generated Commit**: Your customized version on the working branch

## Benefits

- **Conflict-Free Updates**: Easily upgrade to new upstream versions without merge conflicts
- **Reproducible Builds**: Your customizations are codified as scripts, not ad-hoc commits
- **Clear Separation**: Upstream code remains untouched; all changes are declarative
- **Version Control for Patches**: Track your customization logic separately from the generated code
- **Simple Workflow**: Works naturally with GitHub's fork workflow

## Installation

Install XGo (if not already installed):

```bash
go install github.com/goplus/gox/xgo@latest
```

No separate installation of gitrans is needed - you'll run it directly using `xgo run .`

## Recommended Workflow

### Step 1: Fork the Upstream Repository

Fork the upstream repository on GitHub (or your Git hosting platform).

```
upstream/repo (original repository)
    ↓ fork
yourusername/repo (your fork)
```

### Step 2: Create a Working Branch

```bash
# Clone your fork
git clone https://github.com/yourusername/repo.git
cd repo

# The main branch tracks upstream
git checkout main

# Create a new working branch for your customizations
git checkout -b foo
```

### Step 3: Set Up gitrans Patch Directory

Choose one of the following locations for your patch script:

**Option A: Use `.gitrans` directory**
```bash
mkdir -p .gitrans
cd .gitrans
```

**Option B: Use `.github` directory**
```bash
mkdir -p .github
cd .github
```

### Step 4: Initialize Go Module

```bash
# Inside .gitrans or .github directory
go mod init git_patch

# Install gitrans dependency
xgo get github.com/goplus/gitrans@latest
```

This will create `go.mod` and `go.sum` files with gitrans as a dependency.

### Step 5: Create Patch Script

Create `git_patch.gox` in the `.gitrans` or `.github` directory:

```go
// TODO
```

### Step 6: Run gitrans

Execute the patch script:

```bash
# Inside .gitrans or .github directory
xgo run .
```

**What happens:**
1. gitrans reads the upstream commit from the target branch (usually `main`)
2. Applies all patches defined in `git_patch.gox`
3. Modifies files in the current working branch (`foo`)
4. Reports success or any errors

### Step 7: Review and Commit Changes

```bash
# Return to repository root
cd ..

# Review the changes made by gitrans
git status
git diff

# Add all changes
git add .

# Commit with a descriptive message
git commit -m "Apply custom patches via gitrans"

# Push the working branch to remote
git push origin foo
```

## Updating from Upstream

When the upstream repository releases new updates:

### Step 1: Update Your Target Branch

```bash
# Switch to target branch
git checkout main

# Add upstream remote (first time only)
git remote add upstream https://github.com/original/repo.git

# Fetch and merge upstream changes
git fetch upstream
git merge upstream/main

# Push updated target branch to your fork
git push origin main
```

### Step 2: Reapply Patches to Working Branch

```bash
# Switch to your working branch
git checkout foo

# Run gitrans to reapply patches
cd .gitrans  # or .github
xgo run .
cd ..
```

### Step 3: Commit and Push

```bash
git add .
git commit -m "Update to upstream v2.0.0 with custom patches"
git push origin foo
```

## gitrans API Reference

TODO

## Integration with CI/CD

Create `.github/workflows/update-patches.yml`:

```yaml
name: Update Patches

on:
  schedule:
    - cron: '0 0 * * 1'  # Weekly on Monday
  workflow_dispatch:

jobs:
  update:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          ref: foo
          
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
          
      - uses: goplus/setup-xgo@v1
        with:
          go-version: '1.21'
          xgo-version: '1.6'
        
      - name: Fetch upstream updates
        run: |
          git remote add upstream https://github.com/original/repo.git
          git fetch upstream
          git checkout main
          git merge upstream/main
          git push origin main
          
      - name: Apply patches
        run: |
          git checkout foo
          cd .gitrans
          xgo run .
          
      - name: Commit and push
        run: |
          git config user.name "GitHub Actions"
          git config user.email "actions@github.com"
          git add .
          git commit -m "Auto-update from upstream" || true
          git push origin foo
```

## FAQ

**Q: Why not just use Git patches or rebase?**

A: Traditional Git patches are brittle and break when upstream code changes structure. gitrans uses semantic transformations that adapt to code changes.

**Q: Can I have multiple working branches with different patches?**

A: Yes! Create different branches (foo, bar, baz) with different `.gitrans/git_patch.gox` files.

**Q: What happens if my patch can't be applied?**

A: gitrans will report an error and stop. You'll need to update your `git_patch.gox` to handle the new upstream structure.

**Q: Can I use gitrans with private repositories?**

A: Yes, gitrans works with any Git repository, public or private.

**Q: Do I need to commit the generated code?**

A: Yes, commit both the patch script (`.gitrans/git_patch.gox`) and the generated code on your working branch.

## Support

- **Issues**: [[GitHub Issues](https://github.com/goplus/gitrans/issues)](https://github.com/goplus/gitrans/issues)

---

**gitrans** - Transform your forks, not your workflow.
