"""
Preflight checks for the image processor.

Verifies that all required tools are installed and R2 credentials are configured.
"""

import shutil
import subprocess
from dataclasses import dataclass
from typing import Optional

from rich.console import Console
from rich.panel import Panel

from .config import R2_BUCKET, R2_ENDPOINT, R2_PROFILE

console = Console()


@dataclass
class CheckResult:
    """Result of a preflight check."""

    name: str
    passed: bool
    message: str
    fix_instruction: Optional[str] = None


def check_command_exists(command: str) -> tuple[bool, str]:
    """Check if a command exists and return its version."""
    path = shutil.which(command)
    if not path:
        return False, ""

    try:
        result = subprocess.run(
            [command, "--version"],
            capture_output=True,
            text=True,
            timeout=5,
        )
        # Extract version from output
        version_output = result.stdout.strip() or result.stderr.strip()
        # Get first line and clean it up
        version_line = version_output.split("\n")[0] if version_output else ""
        return True, version_line
    except (subprocess.TimeoutExpired, subprocess.SubprocessError):
        return True, "version unknown"


def check_jpegoptim() -> CheckResult:
    """Check if jpegoptim is installed."""
    exists, version = check_command_exists("jpegoptim")
    if exists:
        return CheckResult(
            name="jpegoptim",
            passed=True,
            message=f"jpegoptim found ({version})" if version else "jpegoptim found",
        )
    return CheckResult(
        name="jpegoptim",
        passed=False,
        message="jpegoptim not found",
        fix_instruction="brew install jpegoptim",
    )


def check_pngquant() -> CheckResult:
    """Check if pngquant is installed."""
    exists, version = check_command_exists("pngquant")
    if exists:
        return CheckResult(
            name="pngquant",
            passed=True,
            message=f"pngquant found ({version})" if version else "pngquant found",
        )
    return CheckResult(
        name="pngquant",
        passed=False,
        message="pngquant not found",
        fix_instruction="brew install pngquant",
    )


def check_aws_cli() -> CheckResult:
    """Check if AWS CLI is installed and configured with the r2 profile."""
    exists, version = check_command_exists("aws")
    if not exists:
        return CheckResult(
            name="AWS CLI",
            passed=False,
            message="AWS CLI not found",
            fix_instruction="brew install awscli",
        )

    # Check if r2 profile is configured
    try:
        result = subprocess.run(
            ["aws", "configure", "list", "--profile", R2_PROFILE],
            capture_output=True,
            text=True,
            timeout=5,
        )
        if result.returncode != 0 or "could not be found" in result.stderr.lower():
            return CheckResult(
                name="AWS CLI",
                passed=False,
                message=f"AWS CLI found but profile '{R2_PROFILE}' not configured",
                fix_instruction=f"""\
1. Get R2 credentials from Cloudflare dashboard
2. Run: aws configure --profile {R2_PROFILE}
3. Set access key, secret key
4. Region: auto
5. Output format: json""",
            )
        return CheckResult(
            name="AWS CLI",
            passed=True,
            message=f"AWS CLI configured (profile: {R2_PROFILE})",
        )
    except (subprocess.TimeoutExpired, subprocess.SubprocessError) as e:
        return CheckResult(
            name="AWS CLI",
            passed=False,
            message=f"Error checking AWS CLI: {e}",
        )


def check_r2_access() -> CheckResult:
    """Check if R2 bucket is accessible."""
    try:
        result = subprocess.run(
            [
                "aws",
                "s3",
                "ls",
                f"s3://{R2_BUCKET}",
                "--profile",
                R2_PROFILE,
                "--endpoint-url",
                R2_ENDPOINT,
            ],
            capture_output=True,
            text=True,
            timeout=10,
        )
        if result.returncode == 0:
            return CheckResult(
                name="R2 Bucket",
                passed=True,
                message=f"R2 bucket accessible ({R2_BUCKET})",
            )
        else:
            error_msg = result.stderr.strip() or "Access denied or bucket not found"
            return CheckResult(
                name="R2 Bucket",
                passed=False,
                message=f"Cannot access R2 bucket: {error_msg}",
                fix_instruction=f"""\
1. Verify bucket name: {R2_BUCKET}
2. Check R2 API token has correct permissions
3. Verify endpoint URL: {R2_ENDPOINT}""",
            )
    except subprocess.TimeoutExpired:
        return CheckResult(
            name="R2 Bucket",
            passed=False,
            message="Timeout connecting to R2",
            fix_instruction="Check network connection and R2 endpoint URL",
        )
    except subprocess.SubprocessError as e:
        return CheckResult(
            name="R2 Bucket",
            passed=False,
            message=f"Error checking R2: {e}",
        )


def run_all_checks() -> list[CheckResult]:
    """Run all preflight checks and return results."""
    return [
        check_jpegoptim(),
        check_pngquant(),
        check_aws_cli(),
        check_r2_access(),
    ]


def display_check_results(results: list[CheckResult]) -> bool:
    """Display check results with Rich formatting. Returns True if all passed."""
    console.print()
    console.print("[bold]Preflight Checks[/bold]")
    console.print()

    all_passed = True
    for result in results:
        if result.passed:
            console.print(f"  [green]✓[/green] {result.message}")
        else:
            all_passed = False
            console.print(f"  [red]✗[/red] {result.message}")
            if result.fix_instruction:
                console.print()
                console.print(
                    Panel(
                        f"[dim]{result.fix_instruction}[/dim]",
                        title=f"[yellow]To fix {result.name}[/yellow]",
                        border_style="yellow",
                        padding=(0, 2),
                    )
                )
                console.print()

    console.print()
    return all_passed


def run_preflight() -> bool:
    """Run preflight checks and display results. Returns True if all passed."""
    results = run_all_checks()
    return display_check_results(results)

