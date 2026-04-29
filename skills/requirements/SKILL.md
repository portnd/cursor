---
name: requirements
description: Design, write, update, and maintain `requirements.md` files for any project. Use proactively when the user asks for requirements, functional scope, acceptance criteria, requirements creation, requirements updates, or when work should start by reading `requirements.md` before planning or implementation.
---

# Requirements Skill

## Purpose

Use this skill to create or maintain a `requirements.md` file that gives agents a shared requirements source of truth.

## Core rules

- Read the existing `requirements.md` before doing any design, planning, coding, or review work.
- If `requirements.md` exists, treat it as the primary requirements source unless the user says otherwise.
- If a change affects scope, behavior, acceptance criteria, timeline, dependencies, or risks, update `requirements.md`.
- Keep the requirements current as the work changes. Do not let implementation drift from the requirements.
- If the requirements are missing or incomplete, create or extend them before implementation starts.

## Requirements workflow

1. Read the current `requirements.md`.
2. Identify goal, users, problem, scope, and constraints.
3. Clarify any missing decisions or make the safest explicit assumption.
4. Write or update the requirements.
5. Use the requirements as the reference for planning and implementation.
6. After implementation changes, update the requirements if behavior changed.

## Requirements structure

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

Update `requirements.md` when any of these change:

- feature behavior
- scope boundaries
- acceptance criteria
- UX flow
- data model or API contract
- dependencies or constraints
- rollout or migration plan
- risks or open questions

## Good habits

- Keep a short change log.
- Record unresolved questions instead of hiding them.
- Tie each requirement to a user need or system need.
- Avoid implementation details unless they are important constraints.

## Output expectations

When using this skill, provide:

- a clear requirements draft or requirements update
- a summary of what changed
- any assumptions made
- any follow-up questions if required
