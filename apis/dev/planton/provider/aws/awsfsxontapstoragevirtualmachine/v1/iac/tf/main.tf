# ---------------------------------------------------------------------------
# AWS FSx for NetApp ONTAP Storage Virtual Machine (SVM)
# ---------------------------------------------------------------------------
# One aws_fsx_ontap_storage_virtual_machine resource carries the whole spec.
# ForceNew attributes (replace the SVM when changed): file_system_id, name,
# and root_volume_security_style. The admin password and the ENTIRE Active
# Directory block (domain join details included) update in place — the SVM
# re-joins the new domain without replacement.
# ---------------------------------------------------------------------------

resource "aws_fsx_ontap_storage_virtual_machine" "this" {
  # Parent file system (ForceNew — an SVM cannot move between file systems).
  file_system_id = var.spec.file_system_id

  # The ONTAP SVM name (ForceNew) — distinct from metadata.name, which only
  # becomes the Name tag. This name appears in DNS endpoints, junction paths,
  # and SnapMirror relationships.
  name = var.spec.name

  # Default security style for volumes created under this SVM (ForceNew).
  # Sent explicitly so the plan states the decision — the spec default (UNIX)
  # is materialized before the module runs.
  root_volume_security_style = var.spec.root_volume_security_style

  # vsadmin password for SVM-scoped ONTAP CLI access. Updatable in place;
  # omitted entirely when unset.
  svm_admin_password = local.svm_admin_password

  # Self-managed Active Directory join — the switch that turns on the SMB
  # endpoint. ONTAP SVMs support only self-managed AD (no AWS Managed
  # Microsoft AD). The whole block updates in place.
  dynamic "active_directory_configuration" {
    for_each = var.spec.active_directory_configuration != null ? [var.spec.active_directory_configuration] : []

    content {
      # NetBIOS computer-object name; AWS generates one when omitted.
      netbios_name = active_directory_configuration.value.netbios_name != "" ? active_directory_configuration.value.netbios_name : null

      self_managed_active_directory_configuration {
        domain_name = active_directory_configuration.value.domain_name
        dns_ips     = active_directory_configuration.value.dns_ips
        username    = active_directory_configuration.value.username
        password    = active_directory_configuration.value.password
        # The AD group granted share-administration rights; the spec default
        # ("Domain Admins") is materialized before the module runs.
        file_system_administrators_group = active_directory_configuration.value.file_system_administrators_group
        # Omitted → the computer object lands in the default Computers
        # container.
        organizational_unit_distinguished_name = active_directory_configuration.value.organizational_unit_distinguished_name != "" ? active_directory_configuration.value.organizational_unit_distinguished_name : null
      }
    }
  }

  tags = local.aws_tags
}
