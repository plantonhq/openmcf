output "model_id" {
  description = "The foundation model identifier the agreement covers. Matches spec.model_id -- exported so charts can order model-consuming components after access is granted."
  value       = aws_bedrock_foundation_model_agreement.this.model_id
}
