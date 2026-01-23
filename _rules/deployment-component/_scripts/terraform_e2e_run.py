#!/usr/bin/env python3
"""
Deterministic tool: Run Terraform (tofu) E2E for a provider/kind using local CLI and hack manifest.

Steps:
  1) make local
  2) project-planton tofu init --manifest <manifest>
  3) project-planton tofu plan --manifest <manifest>
  4) project-planton tofu apply --manifest <manifest> --auto-approve
  5) project-planton tofu destroy --manifest <manifest> --auto-approve

Usage:
  python3 _rules/deployment-component/_scripts/terraform_e2e_run.py --provider aws --kindfolder awscloudfront --manifest <path>

Outputs JSON with each step's exit_code/stdout/stderr and success booleans.
"""

import argparse
import json
import os
import subprocess
import sys


def find_repo_root(start_dir: str) -> str:
    current = os.path.abspath(start_dir)
    while True:
        if os.path.isdir(os.path.join(current, ".git")) or os.path.isfile(os.path.join(current, "go.mod")):
            return current
        parent = os.path.dirname(current)
        if parent == current:
            return start_dir
        current = parent


def run(cmd, cwd=None, env=None):
    try:
        p = subprocess.run(cmd, cwd=cwd, env=env, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=False)
        return {"exit_code": p.returncode, "stdout": p.stdout, "stderr": p.stderr}
    except Exception as exc:
        return {"exit_code": 127, "stdout": "", "stderr": str(exc)}


def main() -> int:
    parser = argparse.ArgumentParser(description="Run Terraform E2E via ProjectPlanton CLI")
    parser.add_argument("--provider", required=True)
    parser.add_argument("--kindfolder", required=True)
    parser.add_argument("--manifest", required=True)
    args = parser.parse_args()

    provider = args.provider.strip().lower().replace("_", "")
    kind = args.kindfolder.strip().lower().replace("_", "")
    if any(x in provider for x in ["..", "/", "~"]) or any(x in kind for x in ["..", "/", "~"]):
        print(json.dumps({"error": "invalid provider/kindfolder"}))
        return 2

    repo_root = os.environ.get("REPO_ROOT", find_repo_root(os.getcwd()))
    env = os.environ.copy()

    result = {
        "repo_root": repo_root,
        "manifest": os.path.abspath(args.manifest),
        "make_local": {},
        "tofu_init": {},
        "tofu_plan": {},
        "tofu_apply": {},
        "tofu_destroy": {},
    }

    # 1) make local
    result["make_local"] = run(["make", "-C", repo_root, "local"], cwd=repo_root, env=env)

    # 2) init
    result["tofu_init"] = run(["project-planton", "tofu", "init", "--manifest", result["manifest"]], cwd=repo_root, env=env)
    # 3) plan
    result["tofu_plan"] = run(["project-planton", "tofu", "plan", "--manifest", result["manifest"]], cwd=repo_root, env=env)
    # 4) apply
    result["tofu_apply"] = run(["project-planton", "tofu", "apply", "--manifest", result["manifest"], "--auto-approve"], cwd=repo_root, env=env)
    # 5) destroy
    result["tofu_destroy"] = run(["project-planton", "tofu", "destroy", "--manifest", result["manifest"], "--auto-approve"], cwd=repo_root, env=env)

    print(json.dumps(result))
    return 0 if all(step.get("exit_code", 1) == 0 for step in [result["tofu_init"], result["tofu_plan"], result["tofu_apply"], result["tofu_destroy"]]) else 4


if __name__ == "__main__":
    sys.exit(main())


