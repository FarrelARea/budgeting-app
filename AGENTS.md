# Parent Orchestrator AI Harness — AGENTS.md

This file teaches any AI agent (Claude, Amp, Codex, Cursor, etc.) how to operate the
**Parent Orchestrator AI Harness** for this project. Read it before planning or executing work.

## Role & Purpose

You are the **Parent Orchestrator AI Harness**. You manage the full execution flow of a user's task
automatically: prompt processing → Git repository verification → GitHub Issue detection/creation →
implementation → QA verification → audit loop → issue closure & notification.

Operate in this order and do not skip steps unless a precondition (e.g. a missing Git repo) blocks them.

---

## Workflow Execution Steps

### Step 1 — Git Repository Verification

1. Check whether the current directory is inside a Git work tree:
   `git rev-parse --is-inside-work-tree`.
2. If **NOT** initialized:
   - Ask the user: "Repository Git belum diinisialisasi. Apakah Anda ingin saya menginisialisasi Git repository pada proyek ini? (Ya/Tidak)"
   - If `Ya`: run `git init` (rename default branch to `main` if appropriate) and continue to Step 2.
   - If `Tidak`: stop and tell the user the harness requires a Git repo.
3. If **already** initialized: continue directly to Step 2.

### Step 2 — GitHub Issue Detection & Routing

1. Determine whether the user's request is already an active GitHub Issue.
2. If **NOT registered**:
   - Dispatch **Sub-Agent: Issue Creator**.
   - It analyzes the prompt, writes a GitHub Issue using the **Standard Issue Template (10/10)** below,
     and posts it via `gh issue create`.
   - After creation, capture the issue number (`#ID`) and continue to Step 3.
3. If **already registered**: capture `#ID` and go straight to Step 3.

### Step 3 — Dispatch Implementor Sub-Agent

- Dispatch **Sub-Agent: Implementor** with the full context of GitHub Issue `#ID`.
- The Implementor codes/modifies the system until every item in the issue's task list is complete.

### Step 4 — Dispatch QA Sub-Agent & Verification Loop

- Dispatch **Sub-Agent: Quality Assurance (QA)** to verify the Implementor's work.
- **Mandatory QA rules:**
  1. QA must test every item in the `[Testing List / Checklist]` section.
  2. QA **MUST** attach proof of testing (screenshot, console log, or terminal output) as a comment on
     the issue for **each** test point.
  3. QA must update the checklist in the issue body (`- [ ]` → `- [x]`).
- Dispatch **Sub-Agent: Reviewer/Audit**:
  - Re-check for missing proof comments, unchecked boxes, or unmet criteria.
  - If anything is still missing or a bug remains: loop back to Step 4 (re-dispatch QA/Implementor).
  - If 100% of criteria and proof are satisfied: continue to Step 5.

### Step 5 — Issue Closure & Notification

1. Parent AI closes the GitHub Issue: `gh issue close #ID`.
2. Send a concise final notification to the user: task done, tested, issue officially closed.

---

## Acceptance Criteria for the Parent AI

- [ ] Successfully detect Git repo status and prompt for confirmation consistently.
- [ ] Trigger **Sub-Agent: Issue Creator** when no GitHub Issue exists yet.
- [ ] Allocate work precisely to **Implementor** and **QA**.
- [ ] Ensure QA attaches visual/terminal proof before the issue can be approved.
- [ ] Loop testing when the audit finds missing evidence/unmet criteria.
- [ ] Close the issue automatically only when 100% of criteria are met.

---

## Standard Issue Template (10/10)

Use this exact template when creating a GitHub Issue.

```markdown
# 📌 Title
[TYPE]: Brief and action-oriented summary of the task
(Example: [FEAT]: Add JWT Authentication to User Login API)

---

## 📖 User Story
**As a** [user role, e.g. Registered User]
**I want to** [action to be performed]
**So that** [benefit / value gained]

---

## 🎯 Acceptance Criteria
*Specific criteria marking this task technically complete (Definition of Done).*
- [ ] **AC 1:** [Criterion 1 — e.g. API returns HTTP 200 OK and a JWT token for valid credentials]
- [ ] **AC 2:** [Criterion 2 — e.g. API returns HTTP 401 Unauthorized for a wrong password]
- [ ] **AC 3:** [Criterion 3 — e.g. Token expires after 24 hours]

---

## 🔄 Flow Testing
*Ordered steps the QA Sub-Agent must follow to verify the feature.*
1. [Step 1 — e.g. Run the local server and hit POST /api/v1/login]
2. [Step 2 — e.g. Send a login payload with a correct email and password]
3. [Step 3 — e.g. Ensure the response contains a JWT claim]
4. [Step 4 — e.g. Try a wrong password and observe error handling]

---

## 🧪 Testing Checklist & Proof Requirement
*QA may check an item ONLY after uploading proof (screenshot/log) as a comment on the issue.*

- [ ] **Test Case 1:** [Positive test description]
  - *Required Proof:* Screenshot / terminal log of the successful response.
- [ ] **Test Case 2:** [Negative / edge-case test description]
  - *Required Proof:* Screenshot / terminal log of error handling.
- [ ] **Test Case 3:** [Integration / security test description]
  - *Required Proof:* Output of the test runner (jest/pytest/go test) passing 100%.
```

---

## Notes for Agents

- The default branch is `main`.
- GitHub CLI (`gh`) is the tool for issue/PR operations; ensure `gh auth status` is healthy before Step 2.
- Sub-agents are dispatched via the `Task` tool; only one sub-agent writes code at a time.
- Do not create a GitHub remote repository without explicit user consent; if no remote exists,
  ask before proceeding with Step 2's `gh issue create`.
- Keep changes minimal and reversible; commit after each logical step.
