'use client';
import { useState, useMemo, useCallback, useEffect } from 'react';
import { Typography, Button, Box, Stack, TextField } from '@mui/material';
import { create } from '@bufbuild/protobuf';
import { TableComp } from '@/components/shared/table';
import { ActionMenuProps, PAGINATION_MODE } from '@/models/table';
import { Drawer } from '@/components/shared/drawer';
import { YamlEditor } from '@/components/shared/yaml-editor';
import { AlertDialog } from '@/components/shared/alert-dialog';
import { Refresh, Add } from '@mui/icons-material';
import { useCloudResourceQuery, useCloudResourceCommand } from '@/app/cloud-resources/_services';
import { ListCloudResourcesRequestSchema } from '@/gen/org/openmcf/app/cloudresource/v1/io_pb';
import { CloudResource } from '@/gen/org/openmcf/app/cloudresource/v1/api_pb';
import { PageInfoSchema } from '@/gen/org/openmcf/app/commons/page_info_pb';
import { formatTimestampToDate } from '@/lib';
import { StackUpdatesDrawer } from '@/components/shared/stack-update';

type DrawerMode = 'view' | 'edit' | 'create' | null;

export interface CloudResourcesListProps {
  /**
   * Optional title for the cloud resources section
   */
  title?: string;
  /**
   * Whether to show the kind filter/search input field
   * @default false
   */
  showKindFilter?: boolean;
  /**
   * Initial kind filter value
   */
  initialKindFilter?: string;
  /**
   * Placeholder text for the kind filter input
   * @default "Filter by kind"
   */
  kindFilterPlaceholder?: string;
  /**
   * Label for the kind filter input (optional, will show as placeholder if not provided)
   */
  kindFilterLabel?: string;
  /**
   * Minimum width for the kind filter input
   * @default 250
   */
  kindFilterMinWidth?: number;
  /**
   * Whether to show error alerts
   * @default true
   */
  showErrorAlerts?: boolean;
  /**
   * Custom container component or styling
   */
  container?: React.ComponentType<{ children: React.ReactNode }>;
  /**
   * Callback function to be called when the cloud resources list changes
   * @default true
   */
  onChange?: () => void;
}

export function CloudResourcesList({
  title = 'Cloud Resources',
  showKindFilter = false,
  initialKindFilter = '',
  kindFilterPlaceholder = 'Filter by kind',
  kindFilterLabel,
  kindFilterMinWidth = 250,
  showErrorAlerts = true,
  container: Container,
  onChange,
}: CloudResourcesListProps) {
  const { query } = useCloudResourceQuery();
  const { command } = useCloudResourceCommand();
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [cloudResources, setCloudResources] = useState<CloudResource[]>([]);
  const [totalPages, setTotalPages] = useState(0);
  const [apiLoading, setApiLoading] = useState(true);
  const [kindFilter, setKindFilter] = useState<string>(initialKindFilter);

  // Drawer state
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerMode, setDrawerMode] = useState<DrawerMode>(null);
  const [selectedResource, setSelectedResource] = useState<CloudResource | null>(null);
  const [formData, setFormData] = useState<{ manifest: string }>({ manifest: '' });

  // Delete confirmation dialog state
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [resourceToDelete, setResourceToDelete] = useState<CloudResource | null>(null);

  // Stack updates drawer state
  const [stackUpdatesDrawerOpen, setStackUpdatesDrawerOpen] = useState(false);
  const [selectedResourceForStackUpdates, setSelectedResourceForStackUpdates] =
    useState<CloudResource | null>(null);

  // Function to call the API
  const handleLoadCloudResources = useCallback(() => {
    setApiLoading(true);
    if (query) {
      query
        .listCloudResources(
          create(ListCloudResourcesRequestSchema, {
            kind: kindFilter.trim() || undefined,
            pageInfo: create(PageInfoSchema, {
              num: page,
              size: rowsPerPage,
            }),
          })
        )
        .then((result) => {
          setCloudResources(result.resources);
          setTotalPages(result.totalPages || 0);
          setApiLoading(false);
        })
        .finally(() => {
          setApiLoading(false);
        });
    }
  }, [query, kindFilter, page, rowsPerPage]);

  useEffect(() => {
    if (onChange) {
      onChange();
    }
  }, [cloudResources]);

  // Auto-load on mount and when dependencies change
  useEffect(() => {
    if (query) {
      handleLoadCloudResources();
    }
  }, [query, handleLoadCloudResources]);

  // Handle kind filter change
  const handleKindFilterChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    const newKind = event.target.value;
    setKindFilter(newKind);
    setPage(0); // Reset to first page when filter changes
  }, []);

  // Convert to TableComp format
  const tableHeaders = useMemo(() => ['Name', 'Kind', 'Created At', 'Updated At', 'Actions'], []);
  const tableDataPath = useMemo<Array<'name' | 'kind' | 'createdAt' | 'updatedAt'>>(
    () => ['name', 'kind', 'createdAt', 'updatedAt'],
    []
  );

  const tableDataTemplate = useMemo(
    () => ({
      createdAt: (val: string, row: CloudResource) => {
        return row.createdAt ? formatTimestampToDate(row.createdAt, 'DD/MM/YYYY') : '-';
      },
      updatedAt: (val: string, row: CloudResource) => {
        return row.updatedAt ? formatTimestampToDate(row.updatedAt, 'DD/MM/YYYY') : '-';
      },
    }),
    []
  );

  // Drawer handlers
  const handleOpenDrawer = useCallback((mode: DrawerMode, resource?: CloudResource) => {
    setDrawerMode(mode);
    if (resource) {
      setSelectedResource(resource);
      setFormData({
        manifest: resource.manifest || '',
      });
    } else {
      setSelectedResource(null);
      setFormData({ manifest: '' });
    }
    setDrawerOpen(true);
  }, []);

  const handleCloseDrawer = useCallback(() => {
    setDrawerOpen(false);
    setDrawerMode(null);
    setSelectedResource(null);
    setFormData({ manifest: '' });
  }, []);

  const handleSave = useCallback(() => {
    if (command) {
      if (drawerMode === 'create') {
        command.create(formData.manifest).then(() => {
          handleLoadCloudResources();
          handleCloseDrawer();
        });
      } else if (drawerMode === 'edit' && selectedResource) {
        command.update(selectedResource.id, formData.manifest).then(() => {
          handleLoadCloudResources();
          handleCloseDrawer();
        });
      }
    }
  }, [
    command,
    drawerMode,
    formData.manifest,
    selectedResource,
    handleLoadCloudResources,
    handleCloseDrawer,
  ]);

  const handleConfirmDelete = useCallback((row: CloudResource) => {
    setResourceToDelete(row);
    setDeleteDialogOpen(true);
  }, []);

  const handleDelete = useCallback(() => {
    if (command) {
      command
        .delete(resourceToDelete.id)
        .then(() => {
          handleLoadCloudResources();
          setDeleteDialogOpen(false);
          setResourceToDelete(null);
        })
        .catch((err: any) => {
          console.error('Delete error:', err);
          setDeleteDialogOpen(false);
          setResourceToDelete(null);
        });
    }
  }, [command, resourceToDelete, handleLoadCloudResources, showErrorAlerts]);

  const handleCancelDelete = useCallback(() => {
    setDeleteDialogOpen(false);
    setResourceToDelete(null);
  }, []);

  const handleOpenStackUpdates = useCallback((row: CloudResource) => {
    setSelectedResourceForStackUpdates(row);
    setStackUpdatesDrawerOpen(true);
  }, []);

  const handleCloseStackUpdates = useCallback(() => {
    setStackUpdatesDrawerOpen(false);
    setSelectedResourceForStackUpdates(null);
  }, []);

  const tableActions: ActionMenuProps<CloudResource>[] = useMemo(
    () => [
      {
        text: 'View',
        handler: (row: CloudResource) => {
          handleOpenDrawer('view', row);
        },
        isMenuAction: true,
      },
      {
        text: 'Edit',
        handler: (row: CloudResource) => {
          handleOpenDrawer('edit', row);
        },
        isMenuAction: true,
      },
      {
        text: 'Stack Updates',
        handler: (row: CloudResource) => {
          handleOpenStackUpdates(row);
        },
        isMenuAction: true,
      },
      {
        text: 'Delete',
        handler: (row: CloudResource) => {
          handleConfirmDelete(row);
        },
        isMenuAction: true,
      },
    ],
    [handleOpenDrawer, handleConfirmDelete, handleOpenStackUpdates]
  );

  const handlePageChange = useCallback((newPage: number, newRowsPerPage: number) => {
    setPage(newPage);
    setRowsPerPage(newRowsPerPage);
  }, []);

  // Function to reload data from API
  const handleRefresh = useCallback(() => {
    handleLoadCloudResources();
  }, [handleLoadCloudResources]);

  // Handle YAML editor change
  const handleYamlChange = useCallback(
    (value: string) => {
      if (drawerMode !== 'view') {
        setFormData({ manifest: value });
      }
    },
    [drawerMode]
  );

  // Handle create button click
  const handleCreateClick = useCallback(() => {
    handleOpenDrawer('create');
  }, [handleOpenDrawer]);

  const content = (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
        <Typography variant="h5">{title}</Typography>
        <Stack flexDirection="row" gap={1}>
          {showKindFilter && (
            <TextField
              size="small"
              value={kindFilter}
              onChange={handleKindFilterChange}
              placeholder={kindFilterPlaceholder}
              label={kindFilterLabel}
              sx={{ minWidth: kindFilterMinWidth }}
            />
          )}
          <Button variant="contained" startIcon={<Add />} onClick={handleCreateClick}>
            Create
          </Button>
          <Button
            variant="outlined"
            color="secondary"
            startIcon={<Refresh />}
            onClick={handleRefresh}
            disabled={apiLoading || !query}
          >
            {apiLoading ? 'Loading...' : 'Refresh'}
          </Button>
        </Stack>
      </Box>

      <TableComp
        data={cloudResources}
        loading={apiLoading}
        options={{
          headers: tableHeaders,
          dataPath: tableDataPath,
          dataTemplate: tableDataTemplate,
          actions: tableActions,
          showPagination: true,
          paginationMode: PAGINATION_MODE.SERVER,
          currentPage: page,
          rowsPerPage: rowsPerPage,
          totalPages: totalPages,
          onPageChange: handlePageChange,
          emptyMessage: 'No cloud resources found',
          border: true,
        }}
      />
    </Box>
  );

  return (
    <>
      {Container ? <Container>{content}</Container> : content}

      {/* Drawer for View/Edit/Create */}
      <Drawer
        open={drawerOpen}
        onClose={handleCloseDrawer}
        title={
          drawerMode === 'view'
            ? 'View Cloud Resource'
            : drawerMode === 'edit'
              ? 'Edit Cloud Resource'
              : 'Create Cloud Resource'
        }
        width={800}
      >
        <Stack spacing={3}>
          <Box>
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
              <Typography variant="subtitle2">Manifest (YAML)</Typography>
              {drawerMode !== 'view' && (
                <Box display="flex" gap={1}>
                  <Button variant="contained" color="secondary" onClick={handleCloseDrawer}>
                    Cancel
                  </Button>
                  <Button variant="contained" color="primary" onClick={handleSave}>
                    {drawerMode === 'edit' ? 'Update' : 'Create'}
                  </Button>
                </Box>
              )}
            </Box>
            <YamlEditor
              value={drawerMode === 'view' ? selectedResource?.manifest || '' : formData.manifest}
              onChange={handleYamlChange}
              readOnly={drawerMode === 'view'}
              height="500px"
            />
          </Box>
        </Stack>
      </Drawer>

      {/* Delete Confirmation Dialog */}
      <AlertDialog
        open={deleteDialogOpen}
        onClose={handleCancelDelete}
        onSubmit={handleDelete}
        title="Delete Cloud Resource"
        subTitle={`Are you sure you want to delete "${resourceToDelete?.name}"? This action cannot be undone.`}
        submitLabel="Delete"
        submitBtnColor="error"
        cancelLabel="Cancel"
      />

      {/* Stack Updates Drawer */}
      {selectedResourceForStackUpdates && (
        <StackUpdatesDrawer
          open={stackUpdatesDrawerOpen}
          cloudResourceId={selectedResourceForStackUpdates.id}
          onClose={handleCloseStackUpdates}
        />
      )}
    </>
  );
}
