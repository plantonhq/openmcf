#!/usr/bin/env python3
"""
Deterministic tool: Run Pulumi E2E for a provider/kind using the local CLI.

Steps performed:
  1) make local (installs CLI against current repo APIs)
  2) pulumi login (defaults to local filesystem backend)
  3) pulumi stack select --create <stack>
  4) openmcf pulumi update --manifest <file> --stack <stack> --module-dir <dir>

Usage examples:
  python3 _rules/deployment-component/_scripts/pulumi_e2e_run.py \
    --provider aws --kindfolder awscloudfront \
    --manifest ./apis/project/planton/provider/aws/awscloudfront/v1/iac/hack/manifest.yaml \
    --stack org/project/stack \
    --pulumi-login file://${HOME}/.pulumi

Outputs JSON:
  - repo_root
  - module_dir_abs, module_dir_rel
  - make_local: {exit_code, stdout, stderr}
  - pulumi_login: {exit_code, stdout, stderr}
  - stack_select: {exit_code, stdout, stderr, stack}
  - cli_update: {exit_code, stdout, stderr}
  - update_succeeded: bool
  - error: optional string
"""

import argparse
import json
import os
import subprocess
import sys
from typing import Tuple


def find_repo_root(start_dir: str) -> str:
    current = os.path.abspath(start_dir)
    while True:
        if os.path.isdir(os.path.join(current, ".git")) or os.path.isfile(os.path.join(current, "go.mod")):
            return current
        parent = os.path.dirname(current)
        if parent == current:
            return start_dir
        current = parent


def module_dir(repo_root: str, provider: str, kind_folder: str) -> Tuple[str, str]:
    rel = os.path.join(
        "apis",
        "project",
        "planton",
        "provider",
        provider,
        kind_folder,
        "v1",
        "iac",
        "pulumi",
        "module",
    )
    return os.path.join(repo_root, rel), rel


def run(cmd, cwd=None, env=None):
    try:
        p = subprocess.run(cmd, cwd=cwd, env=env, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=False)
        return {"exit_code": p.returncode, "stdout": p.stdout, "stderr": p.stderr}
    except Exception as exc:
        return {"exit_code": 127, "stdout": "", "stderr": str(exc)}


def main() -> int:
    parser = argparse.ArgumentParser(description="Run Pulumi E2E using local CLI")
    parser.add_argument("--provider", required=True)
    parser.add_argument("--kindfolder", required=True)
    parser.add_argument("--manifest", required=True, help="Path to manifest.yaml for this resource")
    parser.add_argument("--stack", required=True, help="Stack in the form org/project/stack")
    parser.add_argument("--pulumi-login", default=None, help="Pulumi backend URL; defaults to local file backend")
    args = parser.parse_args()

    provider = args.provider.strip().lower().replace("_", "")
    kind = args.kindfolder.strip().lower().replace("_", "")
    if any(x in provider for x in ["..", "/", "~"]) or any(x in kind for x in ["..", "/", "~"]):
        print(json.dumps({"error": "invalid provider/kindfolder"}))
        return 2

    repo_root = os.environ.get("REPO_ROOT", find_repo_root(os.getcwd()))
    module_abs, module_rel = module_dir(repo_root, provider, kind)

    result = {
        "repo_root": repo_root,
        "module_dir_abs": module_abs,
        "module_dir_rel": module_rel,
        "make_local": {},
        "pulumi_login": {},
        "stack_select": {"stack": args.stack},
        "cli_update": {},
        "update_succeeded": False,
    }

    # Prepare environment: ensure non-interactive pulumi by setting a passphrase for local backend if needed
    env = os.environ.copy()
    env.setdefault("PULUMI_CONFIG_PASSPHRASE", "notsecret")

    # 1) make local
    result["make_local"] = run(["make", "-C", repo_root, "local"], cwd=repo_root, env=env)

    # 2) pulumi login (default to local file backend if not provided)
    login_url = args.pulumi_login
    if not login_url:
        home = os.path.expanduser("~")
        login_url = f"file://{os.path.join(home, '.pulumi')}"
    result["pulumi_login"] = run(["pulumi", "login", login_url], cwd=repo_root, env=env)

    # 3) pulumi stack select --create
    result["stack_select"] = {
        **result["stack_select"],
        **run(["pulumi", "stack", "select", args.stack, "--create"], cwd=repo_root, env=env),
    }

    # 4) Run openmcf CLI pulumi update
    cli_cmd = [
        "openmcf",
        "pulumi",
        "update",
        "--manifest",
        os.path.abspath(args.manifest),
        "--stack",
        args.stack,
        "--module-dir",
        module_abs,
    ]
    result["cli_update"] = run(cli_cmd, cwd=repo_root, env=env)
    result["update_succeeded"] = result["cli_update"].get("exit_code", 1) == 0

    print(json.dumps(result))
    return 0 if result["update_succeeded"] else 4


if __name__ == "__main__":
    sys.exit(main())


