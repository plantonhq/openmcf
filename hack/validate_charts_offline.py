#!/usr/bin/env python3
"""Offline InfraChart validation: render templates with default values and
prove the rendered Cloud Resource manifests against the repo's own contracts.

The authoritative chart proof is the platform's `planton chart build`, which
renders with the production Jinjava engine and validates against the schemas
the target control plane runs. That proof needs a backend whose schemas match
this repo's protos. This tool closes the gap for chart authoring inside this
repo: it validates a chart offline, against THIS checkout's protos, with no
backend at all.

Four checks per chart:

1. RENDER  - every file under templates/ renders with the chart's default
   values (plus the platform-injected `env`) and parses as YAML. Rendering
   uses Python Jinja2 with the Jinjava-compatible filters charts rely on
   (`bool`); Jinjava (JVM) remains the production renderer, so a construct
   outside the shared dialect can still diverge - keep templates to the
   documented constructs and treat `chart build` as the final word.
2. SCHEMA  - every rendered document passes `planton validate-manifest`
   (protovalidate over the spec, using the proto contracts compiled into the
   binary). Point PLANTON_BIN at a binary built from this checkout so the
   schemas match the tree you are editing.
3. REFERENCES - every `valueFrom` resolves: the referenced kind + name must
   be a document defined by the same chart. Charts must be self-contained;
   a reference to a resource the chart does not create would only fail at
   deploy time, mid-environment. Charts must use the explicit
   {kind, name, fieldPath} form (the authoring contract): a name-only
   shorthand (legal in standalone manifests, where the field's default-kind
   annotation fills the gap) is rejected here because it cannot be resolved
   or output-checked without compiling the proto annotations -- it would
   pass this gate unverified.
4. OUTPUT FIELDS - every `valueFrom.fieldPath` of the `status.outputs.<field>`
   form names a field that actually exists in the referenced kind's
   stack_outputs.proto. This catches the silent-breakage class where a chart
   references an output that was renamed or never existed.

Usage:
    hack/validate_charts_offline.py [chart-dir ...]

With no arguments, validates every chart under charts/ (directories holding a
Chart.yaml). Exits non-zero if any check fails. Environment:
    PLANTON_BIN  path to the planton CLI (default: "planton" on PATH)
    CHART_ENV    value injected as `values.env` (default: "dev")

Dependencies: python3 with jinja2 + pyyaml (see charts/Makefile, which runs
this through a self-managed virtualenv).
"""

import os
import re
import subprocess
import sys
import tempfile
from pathlib import Path

import jinja2
import yaml

REPO_ROOT = Path(__file__).resolve().parent.parent
CHARTS_ROOT = REPO_ROOT / "charts"
PLANTON_BIN = os.environ.get("PLANTON_BIN", "planton")
CHART_ENV = os.environ.get("CHART_ENV", "dev")


def jinjava_bool(value):
    """Jinjava's `bool` filter: charts use it because rendered param values
    are strings; accept native bools so defaults typed as YAML booleans in
    values.yaml behave identically."""
    if isinstance(value, bool):
        return value
    return str(value).strip().lower() in ("true", "1", "yes")


def load_default_values(chart_dir: Path) -> dict:
    """Build the render context from values.yaml defaults.

    Values render as strings (matching the platform, where params are string
    typed and Jinjava stringifies them); booleans normalize to "true"/"false"
    so `{{ ... }}` emissions parse as YAML booleans exactly as the JVM
    renderer's output does.
    """
    with open(chart_dir / "values.yaml") as f:
        doc = yaml.safe_load(f) or {}
    values = {}
    for param in doc.get("params", []):
        v = param.get("value", "")
        if isinstance(v, bool):
            v = "true" if v else "false"
        values[param["name"]] = "" if v is None else str(v)
    # The platform injects the deployment environment at render time; charts
    # use it to prefix every resource name.
    values.setdefault("env", CHART_ENV)
    return values


def render_chart(chart_dir: Path, values: dict):
    """Render every template with defaults. Returns (docs, errors) where each
    doc is (source_file, parsed_yaml_dict)."""
    env = jinja2.Environment(
        undefined=jinja2.StrictUndefined,  # an undeclared param is a defect, not an empty string
        keep_trailing_newline=True,
    )
    env.filters["bool"] = jinjava_bool

    docs, errors = [], []
    template_files = sorted((chart_dir / "templates").rglob("*.yaml"))
    if not template_files:
        errors.append((chart_dir, "templates/ contains no .yaml files"))
        return docs, errors

    for tf in template_files:
        rel = tf.relative_to(chart_dir)
        try:
            rendered = env.from_string(tf.read_text()).render(values=values)
        except jinja2.TemplateError as e:
            errors.append((rel, f"render failed: {e}"))
            continue
        try:
            for parsed in yaml.safe_load_all(rendered):
                if parsed is None:
                    continue
                if not isinstance(parsed, dict):
                    errors.append((rel, f"rendered document is not a mapping: {type(parsed).__name__}"))
                    continue
                docs.append((rel, parsed))
        except yaml.YAMLError as e:
            errors.append((rel, f"rendered output is not valid YAML: {e}"))
    return docs, errors


def validate_schema(rel, doc):
    """Run `planton validate-manifest` on one rendered document."""
    problems = []
    for key in ("apiVersion", "kind"):
        if key not in doc:
            problems.append((rel, f"document missing {key}"))
    name = (doc.get("metadata") or {}).get("name")
    if not name:
        problems.append((rel, "document missing metadata.name"))
    if problems:
        return problems

    with tempfile.NamedTemporaryFile("w", suffix=".yaml", delete=False) as tmp:
        yaml.safe_dump(doc, tmp, sort_keys=False)
        tmp_path = tmp.name
    try:
        result = subprocess.run(
            [PLANTON_BIN, "validate-manifest", tmp_path],
            capture_output=True, text=True,
        )
        if result.returncode != 0:
            detail = (result.stdout + result.stderr).strip()
            problems.append((rel, f"{doc.get('kind')}/{name}: validate-manifest failed:\n{detail}"))
    finally:
        os.unlink(tmp_path)
    return problems


def iter_value_from(node, path="spec"):
    """Yield every (path, valueFrom-dict) in a rendered document's spec --
    including incomplete ones, so the reference checks can reject the
    name-only shorthand instead of silently skipping it."""
    if isinstance(node, dict):
        vf = node.get("valueFrom")
        if isinstance(vf, dict):
            yield path, vf
        for k, v in node.items():
            if k != "valueFrom":
                yield from iter_value_from(v, f"{path}.{k}")
    elif isinstance(node, list):
        for i, item in enumerate(node):
            yield from iter_value_from(item, f"{path}[{i}]")


def stack_outputs_proto(api_version: str, kind: str) -> Path:
    """Map a manifest's apiVersion + kind to the kind's stack_outputs.proto.
    Convention: <provider>.planton.dev/v1 + Kind ->
    apis/dev/planton/provider/<provider>/<kind lowercased>/v1/stack_outputs.proto
    (provider directories drop hyphens: digital-ocean -> digitalocean)."""
    provider = api_version.split(".")[0].replace("-", "")
    return (REPO_ROOT / "apis" / "dev" / "planton" / "provider"
            / provider / kind.lower() / "v1" / "stack_outputs.proto")


def validate_references(docs):
    """Checks 3 + 4: intra-chart resolution and output-field existence."""
    problems = []
    defined = {(d.get("kind"), (d.get("metadata") or {}).get("name")) for _, d in docs}
    proto_cache = {}

    for rel, doc in docs:
        holder = f"{doc.get('kind')}/{(doc.get('metadata') or {}).get('name')}"
        for path, vf in iter_value_from(doc.get("spec") or {}):
            target_kind, target_name = vf.get("kind"), vf.get("name")
            field_path = vf.get("fieldPath", "")

            # Charts must spell references out (the authoring contract): a
            # name-only valueFrom relies on default-kind proto annotations
            # this validator does not compile, so it would bypass both the
            # resolution and the output-field checks below.
            if not target_kind or not field_path:
                problems.append((rel, f"{holder}: {path} uses a name-only valueFrom -- charts must "
                                      "use the explicit form {kind, name, fieldPath} so references "
                                      "are verifiable"))
                continue

            if (target_kind, target_name) not in defined:
                problems.append((rel, f"{holder}: {path} references {target_kind}/{target_name}, "
                                      "which this chart does not define"))
                continue

            m = re.match(r"^status\.outputs\.([A-Za-z0-9_]+)", field_path)
            if not m:
                problems.append((rel, f"{holder}: {path} fieldPath '{field_path}' does not follow "
                                      "the status.outputs.<field> form"))
                continue
            field = m.group(1)

            target_api = next(d.get("apiVersion") for _, d in docs
                              if d.get("kind") == target_kind
                              and (d.get("metadata") or {}).get("name") == target_name)
            proto = stack_outputs_proto(target_api, target_kind)
            if proto not in proto_cache:
                proto_cache[proto] = proto.read_text() if proto.exists() else None
            content = proto_cache[proto]
            if content is None:
                problems.append((rel, f"{holder}: cannot locate stack_outputs.proto for {target_kind} "
                                      f"(expected {proto.relative_to(REPO_ROOT)})"))
                continue
            if not re.search(rf"\b{re.escape(field)} = \d+", content):
                problems.append((rel, f"{holder}: {path} references output '{field}' which does not "
                                      f"exist in {target_kind}'s stack_outputs.proto"))
    return problems


def validate_chart(chart_dir: Path) -> list:
    values = load_default_values(chart_dir)
    docs, problems = render_chart(chart_dir, values)
    for rel, doc in docs:
        problems.extend(validate_schema(rel, doc))
    problems.extend(validate_references(docs))
    return problems


def main():
    if len(sys.argv) > 1:
        chart_dirs = [Path(a).resolve() for a in sys.argv[1:]]
    else:
        chart_dirs = sorted(p.parent for p in CHARTS_ROOT.glob("*/*/Chart.yaml"))

    failed = 0
    for chart_dir in chart_dirs:
        label = str(chart_dir.relative_to(CHARTS_ROOT)) if chart_dir.is_relative_to(CHARTS_ROOT) else str(chart_dir)
        problems = validate_chart(chart_dir)
        if problems:
            failed += 1
            print(f"✘ {label}")
            for rel, msg in problems:
                print(f"    [{rel}] {msg}")
        else:
            print(f"✔ {label}")

    print(f"\nOffline chart validation: {len(chart_dirs) - failed} passed, {failed} failed "
          f"out of {len(chart_dirs)} charts")
    sys.exit(1 if failed else 0)


if __name__ == "__main__":
    main()
