# Folded aliases: an alias's identity IS this state machine (one alias set
# per machine), so aliases live here rather than as their own kind. Each
# entry is keyed by its name -- adding, renaming, or removing one alias
# never touches its siblings -- and routes 100% of traffic to the version
# THIS deployment published (spec CEL guarantees publish: true whenever
# aliases exist). Weighted canary routing between two specific versions is
# an imperative deployment-shift operation and is deliberately not modeled.
resource "aws_sfn_alias" "this" {
  for_each = { for a in var.spec.aliases : a.name => a }

  name        = each.value.name
  description = each.value.description != "" ? each.value.description : null

  routing_configuration {
    state_machine_version_arn = aws_sfn_state_machine.this.state_machine_version_arn
    weight                    = 100
  }
}
