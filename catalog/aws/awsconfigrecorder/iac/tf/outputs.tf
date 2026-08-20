output "recorder_name" {
  description = "The recorder's name (AWS's regional singleton convention; also the provider's import ID)"
  value       = aws_config_configuration_recorder.this.name
}

output "delivery_channel_name" {
  description = "The delivery channel's name (set only when spec.delivery_channel is configured)"
  value       = length(aws_config_delivery_channel.this) > 0 ? aws_config_delivery_channel.this[0].name : ""
}

output "recording_enabled" {
  description = "Whether the recorder is running after apply (the folded recorder-status truth)"
  value       = aws_config_configuration_recorder_status.this.is_enabled
}
