# One EBS volume - a fresh volume XOR a copy of an existing one -
# with attachments managed in-line.
#
# Lifecycle facts the render below depends on:
#   - the two arms are two provider resources (aws_ebs_volume /
#     aws_ebs_volume_copy); exactly one exists per the spec's union
#     CEL, and both expose the same downstream surface (id/arn);
#   - a copy inherits the source's availability zone, encryption
#     posture, and snapshot lineage - the provider offers no override,
#     which is why the spec forbids those fields on the copy arm;
#   - size/type/iops/throughput update in place on both arms;
#     everything else replaces;
#   - attachments are ForceNew per (volume, instance, device) - a
#     changed device name detaches and re-attaches;
#   - final_snapshot is config-only at AWS (never read back), so
#     imports do not round-trip it.

locals {
  is_copy = var.spec.copy_from != null

  # Optional numeric dials: 0 means "let AWS default it" and is never
  # sent.
  size       = var.spec.size_gb > 0 ? var.spec.size_gb : null
  iops       = var.spec.iops > 0 ? var.spec.iops : null
  throughput = var.spec.throughput_mibps > 0 ? var.spec.throughput_mibps : null
  type       = var.spec.type != "" ? var.spec.type : null
}

# The CREATE arm: a fresh volume in a chosen availability zone.
resource "aws_ebs_volume" "this" {
  count = local.is_copy ? 0 : 1

  availability_zone = var.spec.availability_zone
  size              = local.size
  type              = local.type
  iops              = local.iops
  throughput        = local.throughput

  snapshot_id = var.spec.snapshot_id != "" ? var.spec.snapshot_id : null
  encrypted   = var.spec.encrypted ? true : null
  kms_key_id  = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null

  multi_attach_enabled = var.spec.multi_attach_enabled ? true : null
  final_snapshot       = var.spec.final_snapshot ? true : null

  volume_initialization_rate = var.spec.volume_initialization_rate > 0 ? var.spec.volume_initialization_rate : null

  tags = local.aws_tags
}

# The COPY arm: clone an existing volume (same zone as the source).
resource "aws_ebs_volume_copy" "this" {
  count = local.is_copy ? 1 : 0

  source_volume_id = var.spec.copy_from.source_volume_id
  size             = local.size
  volume_type      = local.type
  iops             = local.iops
  throughput       = local.throughput

  tags = local.aws_tags
}

locals {
  # The one volume either arm produced - the downstream surface is
  # identical.
  volume_id                = local.is_copy ? aws_ebs_volume_copy.this[0].id : aws_ebs_volume.this[0].id
  volume_arn               = local.is_copy ? aws_ebs_volume_copy.this[0].arn : aws_ebs_volume.this[0].arn
  volume_availability_zone = local.is_copy ? aws_ebs_volume_copy.this[0].availability_zone : aws_ebs_volume.this[0].availability_zone
  volume_size              = local.is_copy ? aws_ebs_volume_copy.this[0].size : aws_ebs_volume.this[0].size
  volume_create_time       = local.is_copy ? "" : aws_ebs_volume.this[0].create_time
}

# In-line attachments, keyed "device//instance" (stable across list
# reorders, and the "//" separator is the import bridge's
# address-key-segment convention - multi-attach conventionally reuses
# the SAME device name per instance, so the instance id
# disambiguates).
resource "aws_volume_attachment" "this" {
  for_each = {
    for attachment in var.spec.attachments :
    "${attachment.device_name}//${attachment.instance_id}" => attachment
  }

  device_name = each.value.device_name
  volume_id   = local.volume_id
  instance_id = each.value.instance_id

  force_detach                   = each.value.force_detach ? true : null
  skip_destroy                   = each.value.skip_destroy ? true : null
  stop_instance_before_detaching = each.value.stop_instance_before_detaching ? true : null
}
