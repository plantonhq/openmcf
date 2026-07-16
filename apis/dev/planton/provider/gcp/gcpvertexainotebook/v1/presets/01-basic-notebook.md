# Basic Notebook

The minimal managed JupyterLab environment: a CPU-only Workbench instance
in the ambient project on the service's default Workbench image.

## What this preset creates

A Workbench instance named `data-exploration` in `us-central1-a` on an
`e2-standard-4` VM (4 vCPU, 16 GB RAM) with a 200 GB SSD boot disk and
the service's latest Workbench image (JupyterLab with the common Python
data stack — pandas, scikit-learn — pre-installed).

## When to use

- Data exploration and visualization
- Feature engineering and preprocessing
- Light scikit-learn model training
- Prototyping before scaling to GPU

## Remix ideas

- Add `acceleratorConfig` and a GPU-capable machine type for deep
  learning training (see the gpu-ml-notebook preset).
- Add `desiredState: STOPPED` to pre-provision a notebook that starts
  billing only when someone flips it to ACTIVE.
- Add `dataDisk` to separate datasets and checkpoints from the OS disk.
- Add `instanceOwners` with one email to lock access to a single user.
