# docs/changes

Per-release changelog. One file per provider release. The release notes are assembled from these files at release time.

## How it works

- **`unreleased.md`** — the file currently being edited during the release-in-progress cycle. PRs that change user-visible behavior append an entry here.
- **`vX.Y.Z.md`** — frozen snapshot of `unreleased.md` at the moment the release was cut. Renamed from `unreleased.md`; a new empty `unreleased.md` is created for the next cycle.

## When to add an entry

Mandatory for any change that's user-visible and not obvious from the auto-generated schema markdown:

- Resource behavior changes (filtering, drift fixes, new Read fields)
- Schema additions, removals, or type changes
- Backend contract changes the provider relies on
- `ConflictsWith` / `Required` / `Optional` / `Default` / `Computed` transitions
- State migrations (`SchemaVersion` bumps)

Skip for: pure refactor with no behavior change, test-only changes, doc-only changes, or dependency bumps without functional impact.

## What to write

Short. Release-note tight. Cover only what shipped — not future plans.

- One H2 heading per change, format: `## <Resource or area> — <one-line summary> (<JIRA-ID>)`
- What changed
- User impact (drift, breaking changes, behavior shifts)
- Migration steps if any (HCL changes, state import notes)

Don't narrate the why or the history. Don't preview future tickets — they get their own entry when they ship.
