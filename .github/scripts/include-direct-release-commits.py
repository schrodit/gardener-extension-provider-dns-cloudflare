#!/usr/bin/env python3
"""Add commits without an associated PR to release-please's GitHub notes."""

import json
import os
import re
import subprocess
import tempfile
from pathlib import Path


START = "<!-- direct-release-commits-start -->"
END = "<!-- direct-release-commits-end -->"


def run(*args: str) -> str:
    return subprocess.check_output(args, text=True).strip()


def update_notes(text: str, section: str, *, changelog: bool) -> str:
    text = re.sub(
        rf"\n*{re.escape(START)}.*?{re.escape(END)}\n*",
        "\n\n",
        text,
        flags=re.DOTALL,
    )
    if changelog:
        heading = re.search(r"^## \d+\.\d+\.\d+[^\n]*$", text, re.MULTILINE)
        if heading is None:
            raise ValueError("release-please changelog has no version heading")
        return text[: heading.end()] + "\n\n" + section + text[heading.end() :].lstrip("\n")
    return text.rstrip() + "\n\n" + section


def main() -> None:
    release_pr = json.loads(os.environ["RELEASE_PR"])
    pr_number = str(release_pr["number"])
    head_branch = release_pr["headBranchName"]
    repository = os.environ["GITHUB_REPOSITORY"]

    previous_version = json.loads(
        run("git", "show", "origin/main:.release-please-manifest.json")
    )["."]
    previous_tag = f"v{previous_version}"
    shas = run("git", "log", "--format=%H", f"{previous_tag}..origin/main").splitlines()
    direct_commits = []
    for sha in reversed(shas):
        associated_prs = run(
            "gh", "api", f"repos/{repository}/commits/{sha}/pulls", "--jq", "length"
        )
        if associated_prs == "0":
            subject = run("git", "show", "-s", "--format=%s", sha)
            direct_commits.append(
                f"* {subject} ([{sha[:7]}](https://github.com/{repository}/commit/{sha}))"
            )

    if not direct_commits:
        return

    section = (
        f"{START}\n### Direct commits\n\n"
        + "\n".join(direct_commits)
        + f"\n{END}\n\n"
    )
    run("git", "fetch", "origin", head_branch)
    run("git", "switch", "--detach", "FETCH_HEAD")
    changelog = Path("CHANGELOG.md")
    updated = update_notes(changelog.read_text(), section, changelog=True)
    changelog.write_text(updated)
    run("git", "add", "CHANGELOG.md")
    if subprocess.run(["git", "diff", "--cached", "--quiet"], check=False).returncode:
        run("git", "config", "user.name", "github-actions[bot]")
        run("git", "config", "user.email", "41898282+github-actions[bot]@users.noreply.github.com")
        run("git", "commit", "-m", "docs: include direct commits in release notes")
        run("git", "push", "origin", f"HEAD:{head_branch}")

    body = run("gh", "pr", "view", pr_number, "--json", "body", "--jq", ".body")
    updated_body = update_notes(body, section, changelog=False)
    if updated_body != body:
        with tempfile.NamedTemporaryFile(mode="w", suffix=".md") as body_file:
            body_file.write(updated_body)
            body_file.flush()
            run("gh", "pr", "edit", pr_number, "--body-file", body_file.name)


if __name__ == "__main__":
    main()
