# Regression Test Harness — Adversarial Code Review

**Branch:** `agentic-regression`
**Commits reviewed:**
- `58aa20a` — Implement first draft of regression test harness.
- `40db327` — Tidying up the AI work; looks pretty sweet now.
**Range:** `1e99d0f..40db327`
**Files reviewed:** 28 files, 1,673 insertions
**Plan reference:** `docs/plans/2026-05-14-regression-test-harness-implementation.md`

---

## Strengths

1. **Self-contained IRCd** — Spinning up Ergo dynamically eliminates external dependencies. Tests are fully isolated and reproducible.

2. **Generic `Process` type** — The plan specified a `BotProcess` type in `bot.go`. The implementation generalizes to `Process` in `process.go`, reused for both the ircd and bot. Cleaner and more maintainable.

3. **Nick-scoped pattern matching** — The plan had standalone `Exact/Contains/Regex` functions with no nick awareness. The implementation makes them `Harness` methods that capture `h.BotNick`, preventing false positives from other users' traffic.

4. **Reverse-order cleanup** — `cleanup()` in `suite.go:176-225` tears down resources in reverse creation order (bot → IRC connection → ircd → temp dir), collecting all errors rather than short-circuiting.

5. **SHA256 verification** — `regtest.sh` downloads Ergo with checksums for both the tarball and extracted binary. This is proper supply-chain hygiene.

6. **Randomized identifiers** — `randSuffix()` with 7 cryptographically random bytes makes channel/nick collisions astronomically unlikely.

7. **Two-stage process shutdown** — `kill()` sends SIGINT first, waits 1s, then cancels the context (SIGKILL). Clean shutdown path.

8. **Comprehensive `process_test.go`** — Covers empty path, nonexistent binary, double-start, natural exit, and forced kill.

---

## Issues

### Critical

#### 1. `errors.AsType` incompatibility with Go 1.25

- **File:** `process.go:136`
- **Issue:** `errors.AsType[*exec.ExitError](waitErr)` is Go 1.26+ syntax. The `go.mod` declares `go 1.25.0`. `go vet` flags this as an error.
- **Impact:** Code will fail to compile or vet on the declared minimum Go version. Hard blocker.
- **Fix:**
  ```go
  var ee *exec.ExitError
  if errors.As(waitErr, &ee) {
      if ws, ok := ee.Sys().(syscall.WaitStatus); ok {
          sig := ws.Signal()
          if sig == syscall.SIGINT || sig == syscall.SIGKILL {
              return nil
          }
      }
  }
  ```

#### 2. `Remover` double-Remove in `ExpectEvent`

- **File:** `harness.go:64-72`
- **Issue:** When the pattern matches, `remover.Remove()` is called explicitly inside the handler, then `defer remover.Remove()` fires on function return. Calling `Remove()` twice on a `client.Remover` is undefined behavior in goirc v1.3.5 — it can panic or corrupt internal state.
- **Impact:** Tests can crash with a panic from goirc internals, or leave the connection in an inconsistent state causing subsequent expectations to silently fail.
- **Fix:**
  ```go
  found := false
  remover = h.HandleFunc(event, func(conn *client.Conn, line *client.Line) {
      if p.Match(line) && !found {
          found = true
          remover.Remove()
          select {
          case resultCh <- line:
          default:
          }
      }
  })
  defer remover.Remove()
  ```

#### 3. `cleanup()` blocks forever on missing DISCONNECTED event

- **File:** `suite.go:195`
- **Issue:** `<-disconnected` has no timeout. If goirc never fires the `DISCONNECTED` event (network error, library bug, connection already dead), `cleanup()` hangs forever.
- **Impact:** `cleanup()` is called from `Start()` on every failure path. A hang here means the test hangs indefinitely, and the ircd process and temp directory are leaked.
- **Fix:**
  ```go
  h.Quit("regtest: cleanup")
  select {
  case <-disconnected:
  case <-time.After(3 * time.Second):
      // Force close; DISCONNECTED event never fired
  }
  ```

#### 4. `running = false` set before `cleanup()` completes

- **File:** `suite.go:166-171`
- **Issue:** `Stop()` sets `running = false` before calling `h.cleanup()`. If `cleanup()` panics or hangs, `running` is already `false`, so a subsequent `Start()` succeeds while old resources (ircd process, temp dir) are still alive.
- **Impact:** Resource leak on cleanup failure.
- **Fix:**
  ```go
  func (h *Harness) Stop() error {
      mu.Lock()
      defer mu.Unlock()
      err := h.cleanup()
      if err == nil {
          running = false
      }
      return err
  }
  ```

### Important

#### 5. `selfValidate()` has a TOCTOU race condition

- **File:** `suite.go:273-286`
- **Issue:** `selfValidate()` sends a PRIVMSG to itself, then registers a `HandleFunc` for the echo. If Ergo delivers the echo before the handler is registered (non-zero probability on localhost), the message is silently dropped and `Start()` fails with a misleading "self-validate PRIVMSG" error.
- **Impact:** Intermittent test failures that are extremely hard to diagnose. The harness would report "self-validate failed" when the real problem is handler registration timing.
- **Fix:** Register the handler *before* sending the message, or use a sequence-numbered approach.

#### 6. Nil `Conn` check missing in `cleanup()`

- **File:** `suite.go:189`
- **Issue:** `h.Connected()` is called without first checking `h.Conn != nil`. If `cleanup()` runs before `h.Conn` is assigned (e.g., failure at `freeLocalAddr()`), calling a method on a nil embedded pointer will panic.
- **Impact:** Panic in cleanup propagates up, potentially leaving ircd process and temp dir leaked.
- **Fix:**
  ```go
  if h.Conn != nil && h.Connected() {
  ```

#### 7. `p.cmd.Process` nil check in `kill()`

- **File:** `process.go:107`
- **Issue:** The check is `if p.cmd != nil`, but `p.cmd.Process` could theoretically be nil if the OS refuses to fork. `p.cmd.Process.Signal()` would panic.
- **Impact:** Panic during process shutdown under resource exhaustion.
- **Fix:**
  ```go
  if p.cmd != nil && p.cmd.Process != nil {
      p.cmd.Process.Signal(os.Interrupt)
  }
  ```

#### 8. Fragile `time.Sleep(100ms)` in seen driver tests

- **File:** `drivers/seendriver/regtest/handler_record_test.go:19,32,40,42,51,53`
- **Issue:** Six `time.Sleep(100 * time.Millisecond)` calls wait for the bot to process and record events before querying. On a loaded CI machine, the bot may not process within 100ms. On a fast machine, this adds 600ms of unnecessary latency.
- **Impact:** Flaky CI results.
- **Fix:** Replace with explicit `h.Expect()` calls to wait for confirmation events, or use a retry loop with a reasonable timeout.

#### 9. `testRecordNick` may leave nick in-flight

- **File:** `drivers/seendriver/regtest/handler_record_test.go:48-58`
- **Issue:** After `h.Nick(old)`, the test only sleeps 100ms before asserting. If the nick change is still propagating, `h.Me().Nick` in `seenCmd` might return the wrong value.
- **Impact:** False negative on the nick-change assertion.
- **Fix:** After `h.Nick(old)`, call `h.Expect()` to wait for the server's nick change confirmation.

#### 10. Test ordering dependency in ignore/unignore tests

- **File:** `bot/regtest/command_ignore_test.go:9-12`
- **Issue:** `testUnignoreBasic` depends on `testIgnoreBasic` having succeeded. If `testIgnoreBasic` fails (e.g., bot response format changed), `testnick` is never added, and `testUnignoreBasic` fails with a confusing error.
- **Impact:** Cascade failures that obscure the root cause.
- **Fix:** Make `testUnignoreBasic` self-contained by doing its own ignore before unignoring.

#### 11. `p.exit` channel created after `p.cmd.Start()`

- **File:** `process.go:93-94`
- **Issue:** `p.exit = make(chan error, 1)` is created after `p.cmd.Start()`. If the process exits between `Start()` returning and `p.exit` being created, `go p.wait()` blocks on send until the channel exists.
- **Impact:** Extremely unlikely in practice, but incorrect ordering.
- **Fix:** Move `p.exit = make(chan error, 1)` before `p.cmd.Start()`.

#### 12. `regtest.sh` only supports Linux x86_64

- **File:** `regtest.sh:14`
- **Issue:** Hardcoded `ergo-2.18.0-linux-x86_64.tar.gz`. Developers on macOS, ARM, or other platforms cannot run the test harness.
- **Impact:** Limits who can contribute and test locally.
- **Fix:** Detect OS/arch and download the appropriate binary, or provide instructions for supplying `REGTEST_IRCD` manually.

#### 13. `regtest.sh` hardcodes `wget`

- **File:** `regtest.sh:35`
- **Issue:** Uses `wget` which is not available on all systems (macOS, some minimal Linux distros).
- **Impact:** Script fails on systems without `wget`.
- **Fix:** Check for `curl` as fallback, or use `curl` (more universally available).

### Minor

#### 14. Hardcoded oper password in IRCd config

- **File:** `ircd.go:155`
- **Issue:** Oper "bob" has a hardcoded bcrypt password. Authentication is disabled and no test uses this oper.
- **Impact:** Low risk (test-only), but poor practice.
- **Fix:** Remove the `opers` block entirely.

#### 15. Dead code comment

- **File:** `suite.go:257`
- **Issue:** `// return "#sp0rklf"` is unreachable dead code.
- **Fix:** Remove the comment.

#### 16. Misleading permission check error message

- **File:** `suite.go:240`
- **Issue:** Error message says "not regular readable executable" but the check is specifically for owner read+execute (`0500`).
- **Fix:** Change message to `"not a regular file with owner read+execute permissions"`.

#### 17. `testLogger` produces identical output for all severity levels

- **File:** `suite.go:19-41`
- **Issue:** `Debug`, `Info`, `Warn`, `Error` all call `t.Logf` identically. No way to distinguish severity in test output.
- **Fix:** Prefix with level: `tl.t.Logf("[DEBUG] "+fmt, args...)`, etc.

#### 18. `cfg.Flood = true` is undocumented

- **File:** `suite.go:115`
- **Issue:** Disabling flood protection is necessary but lacks explanation.
- **Fix:** Add comment: `"disable flood protection so rapid command-and-expect sequences don't trigger rate limiting"`.

#### 19. `TestFullLifecycle` couples to bot's help command

- **File:** `regtest/integration_test.go:28`
- **Issue:** Tests `CommandAndExpect("help", ...)` which depends on the bot's help command implementation. If the help response changes, the integration test breaks.
- **Fix:** Use a more generic assertion or a different self-test command.

#### 20. `testRecordKick` declared but skipped

- **File:** `drivers/seendriver/regtest/suite_test.go:36`
- **Issue:** `t.Run("kick", testRecordKick)` is declared but `testRecordKick` immediately calls `t.Skip()`. Appears as a test that ran but was skipped.
- **Fix:** Remove the declaration until a real implementation exists.

#### 21. `freeLocalAddr` has a TOCTOU race on port binding

- **File:** `ircd.go:18-26`
- **Issue:** Opens a listener on `localhost:0`, extracts the port, closes it, then the ircd binds to it. Between close and bind, another process could grab the port.
- **Impact:** Extremely unlikely on localhost, but technically a race.
- **Fix:** Accept as-is (extremely low risk), or pass the listener directly to the ircd if the config supported it.

#### 22. `Process.wait()` double-read from closed channel

- **File:** `process.go:99-105`
- **Issue:** `kill()` reads from `p.exit` twice in worst case. The second read from a closed channel returns `nil`, which could mask an actual error.
- **Impact:** Confusing code, but safe in practice.
- **Fix:** Restructure `kill()` to avoid the double-read pattern.

#### 23. `set -x` in `regtest.sh` is noisy

- **File:** `regtest.sh:3`
- **Issue:** `set -x` prints every command executed, which is very verbose and makes it harder to spot actual errors in CI logs.
- **Fix:** Remove `set -x` and add explicit `echo` statements for key steps, or use `set -x` only for the download/extract section.

#### 24. Global mutex prevents parallel test suites

- **File:** `suite.go:47-50`
- **Issue:** Package-level `mu` and `running` flag means only one harness can be active at a time. `go test ./bot/regtest/ ./drivers/seendriver/regtest/` would fail the second suite.
- **Impact:** Test suites cannot be parallelized.
- **Fix:** Document the limitation, or refactor to per-Harness state.

#### 25. `regtest.sh` comment typo

- **File:** `regtest.sh:50`
- **Issue:** `# disable parallel testing because it has a sad?` — unclear comment.
- **Fix:** Replace with `# disable parallel testing because the global harness mutex only allows one harness at a time`.

---

## Plan vs. Implementation Differences

| Plan | Implementation | Assessment |
|------|---------------|------------|
| `BotProcess` type in `bot.go` | Generic `Process` type in `process.go` | **Improvement** — reusable for both ircd and bot |
| `Start()` uses `REGTEST_SERVER` env var | `Start()` uses `REGTEST_BOT` + `REGTEST_IRCD` | **Improvement** — more explicit, better design |
| `findBotBinary()` falls back to `go build` | `getBinaryPath()` requires env vars, no fallback | **Deviation** — stricter, faster, less convenient |
| `BotNick()` / `Channel()` methods | `BotNick` / `Channel` exported fields | **Deviation** — simpler, less encapsulation |
| `SendAndExpect(msg, pattern, timeout)` | Hardcoded 2s timeout | **Regression** — no configurable timeout |
| `Start()` / `Stop()` package functions | `Start()` package func, `Stop()` harness method | **Improvement** — Stop() on specific instance |
| Standalone `Exact/Contains/Regex` | Harness methods with nick-scoping | **Improvement** — prevents false positives |

---

## Assessment

**Ready to merge? No — requires fixes.**

**Reasoning:** Four critical issues must be addressed before merging:

1. **`errors.AsType`** (`process.go:136`) — compile/vet failure on Go 1.25. One-line fix.
2. **`Remover` double-Remove** (`harness.go:64-72`) — can panic or corrupt goirc state. Tests cannot be trusted until fixed.
3. **`cleanup()` DISCONNECTED hang** (`suite.go:195`) — no timeout means cleanup can block forever, leaking ircd processes and temp directories.
4. **`running = false` before cleanup** (`suite.go:166-171`) — panicked cleanup leaves system in inconsistent state.

The important issues (nil `Conn` in cleanup, `p.cmd.Process` nil check, `selfValidate()` race, fragile `time.Sleep` calls) should also be addressed before merging to prevent intermittent CI failures.

The core architecture is sound: the self-contained IRCd, generic process management, nick-scoped patterns, and SHA256-verified downloads are well-designed. The implementation in several ways improves on the plan (generic Process, nick-scoping, explicit env vars).
