---
name: prd
description: Design, write, update, and maintain PRD.md files for any project. Use proactively when the user asks for product requirements, feature scope, acceptance criteria, PRD creation, PRD updates, or when work should start by reading PRD before implementation.
---

# PRD Skill

## Purpose

Use this skill to create or maintain a `PRD.md` file that gives agents a shared product source of truth.

## Core rules

- Read the existing `PRD.md` before doing any design, planning, coding, or review work.
- If `PRD.md` exists, treat it as the primary product source unless the user says otherwise.
- If a change affects scope, behavior, acceptance criteria, timeline, dependencies, or risks, update `PRD.md`.
- Keep the PRD current as the work changes. Do not let implementation drift from the PRD.
- If the PRD is missing or incomplete, create or extend it before implementation starts.

## PRD workflow

1. Read the current `PRD.md`.
2. Identify goal, users, problem, scope, and constraints.
3. Clarify any missing decisions or make the safest explicit assumption.
4. Write or update the PRD.
5. Use the PRD as the reference for planning and implementation.
6. After implementation changes, update the PRD if behavior changed.

## PRD structure

Use this structure unless the project already has a stronger standard:

- Title
- Overview
- Problem statement
- Goals
- Non-goals
- Users / personas
- Scope
- User stories or use cases
- Functional requirements
- Non-functional requirements
- Acceptance criteria
- Dependencies
- Risks and edge cases
- Open questions
- Change log

## Writing guidance

- Be concrete and testable.
- Write requirements so another agent can implement them without guessing.
- Prefer bullets and short sections over long prose.
- State assumptions explicitly.
- Separate must-have requirements from nice-to-have ideas.
- Keep acceptance criteria measurable.

## Update guidance

Update `PRD.md` when any of these change:

- feature behavior
- scope boundaries
- acceptance criteria
- UX flow
- data model or API contract
- dependencies or constraints
- rollout or migration plan
- risks or open questions

## Good PRD habits

- Keep a short change log.
- Record unresolved questions instead of hiding them.
- Tie each requirement to a user need or system need.
- Avoid implementation details unless they are important constraints.

## Output expectations

When using this skill, provide:

- a clear PRD draft or PRD update
- a summary of what changed
- any assumptions made
- any follow-up questions if required
