# Web-Test Availability Alert

This preset pages when an Application Insights availability (web) test fails from multiple locations at once -- the outside-in signal that real users cannot reach the site, independent of any server-side metric looking healthy.

## When to Use

- Every user-facing endpoint with an availability test (the test itself is configured on the Application Insights resource in Azure)
- The severity-0 pager: multi-location failure means the site is down, not flaky

## Key Configuration Choices

- **3 failed locations** (`failedLocationCount: 3`) -- tolerates single-location blips; the convention is the test's location count minus two
- **Severity 0** -- a confirmed outage is the highest urgency Azure models
- **Both the test and the scope reference the web test by FK** -- the scope tells Azure what to evaluate and the criteria's `webTestId` names the test; wiring both to the `AzureApplicationInsightsStandardWebTest`'s `web_test_id` output means the alert follows the test with no hand-copied ARM id. The criteria's `componentId` binds it to the Application Insights resource the same way.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-homepage-web-test` | The availability test to watch | `AzureApplicationInsightsStandardWebTest` status outputs |
| `my-web-app-insights` | The Application Insights resource | `AzureApplicationInsights` status outputs |
| `my-platform-oncall` | The action group to page | `AzureMonitorActionGroup` status outputs |

## Related Presets

- **01-static-threshold** -- Server-side metric lines
- **02-dynamic-anomaly** -- Learned-band anomaly detection
