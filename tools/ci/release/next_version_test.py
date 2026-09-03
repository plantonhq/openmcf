#!/usr/bin/env python3
"""Tests for next_version.py: run with `python3 -m unittest tools/ci/release/next_version_test.py`."""
import unittest

from next_version import bump_version, latest_version


class LatestVersionTest(unittest.TestCase):
    def test_bare_tags_ignore_prerelease_and_metadata(self):
        tags = ["v0.5.26-terraform.kubernetesx.20250101.0", "v0.5.26", "v0.5.25"]
        self.assertEqual(latest_version(tags, ""), "v0.5.26")

    def test_prefix_reads_only_its_namespace(self):
        tags = ["operator/v0.8.1", "operator/v0.8.0", "v0.5.26"]
        self.assertEqual(latest_version(tags, "operator/"), "v0.8.1")
        self.assertEqual(latest_version(tags, "helm/planton/"), "v0.0.0")

    def test_bare_reads_only_bare_tags(self):
        tags = ["operator/v9.9.9", "helm/planton/v9.9.9", "v0.5.26"]
        self.assertEqual(latest_version(tags, ""), "v0.5.26")

    def test_empty_namespace_is_v0(self):
        self.assertEqual(latest_version([], "helm/planton-runner/"), "v0.0.0")


class BumpVersionTest(unittest.TestCase):
    def test_bumps(self):
        self.assertEqual(bump_version("v0.8.1", "patch"), "v0.8.2")
        self.assertEqual(bump_version("v0.8.1", "minor"), "v0.9.0")
        self.assertEqual(bump_version("v0.8.1", "major"), "v1.0.0")

    def test_first_tag_of_a_namespace_starts_at_patch_one(self):
        self.assertEqual(bump_version("v0.0.0", "patch"), "v0.0.1")

    def test_rejects_non_semver(self):
        with self.assertRaises(ValueError):
            bump_version("0.8.1", "patch")


if __name__ == "__main__":
    unittest.main()
