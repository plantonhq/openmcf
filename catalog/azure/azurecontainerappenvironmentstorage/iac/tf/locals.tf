locals {
  # Access-mode wire values. The storage registration carries no tags
  # (ARM does not support them on managedEnvironments/storages), so no
  # tag locals exist here.
  access_mode_map = {
    "READ_ONLY"  = "ReadOnly"
    "READ_WRITE" = "ReadWrite"
  }
}
