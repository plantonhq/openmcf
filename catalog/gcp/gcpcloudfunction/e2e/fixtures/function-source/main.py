# Minimal Gen 2 HTTP function used by the live E2E scenarios. The test
# entrypoint zips this directory and stages it in a run-scoped GCS bucket
# before the scenario applies (Cloud Build needs real source to build).
import functions_framework


@functions_framework.http
def handler(request):
    return "ok"
