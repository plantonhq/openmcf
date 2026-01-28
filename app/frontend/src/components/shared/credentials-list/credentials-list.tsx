'use client';
import { useState, useMemo, useCallback, useEffect } from 'react';
import { Typography, Button, Box, Stack } from '@mui/material';
import { TableComp } from '@/components/shared/table';
import { ActionMenuProps } from '@/models/table';
import { AlertDialog } from '@/components/shared/alert-dialog';
import { SimpleSelect } from '@/components/shared/simple-select';
import { CredentialDrawer, DrawerMode } from '@/app/credentials/_components/forms';
import { Refresh, Add } from '@mui/icons-material';
import { useCredentialQuery, useCredentialCommand } from '@/app/credentials/_services';
import {
  ListCredentialsRequestSchema,
  CredentialSummary,
} from '@/gen/org/openmcf/app/credential/v1/io_pb';
import { Credential, Credential_CredentialProvider } from '@/gen/org/openmcf/app/credential/v1/api_pb';
import { formatTimestampToDate } from '@/lib';
import { create } from '@bufbuild/protobuf';
import { providerConfig } from '@/app/credentials/_components/utils';

export function CredentialsList() {
  const { query } = useCredentialQuery();
  const { command } = useCredentialCommand();
  const [credentials, setCredentials] = useState<CredentialSummary[]>([]);
  const [apiLoading, setApiLoading] = useState(true);
  const [providerFilter, setProviderFilter] = useState<Credential_CredentialProvider | undefined>();

  // Drawer state
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [drawerMode, setDrawerMode] = useState<DrawerMode>(null);
  const [selectedCredential, setSelectedCredential] = useState<Credential | null>(null);

  // Delete confirmation dialog state
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [credentialToDelete, setCredentialToDelete] = useState<CredentialSummary | null>(null);

  const providerFilterOptions = useMemo(() => {
    const allOption = {
      label: 'All',
      value: Credential_CredentialProvider.CREDENTIAL_PROVIDER_UNSPECIFIED.toString(),
    };
    const providerOptions = (Object.keys(providerConfig) as unknown as Array<Credential_CredentialProvider>)
      .filter((provider) => {
        // Filter out UNSPECIFIED (value 0) by comparing numeric enum values
        return Number(provider) !== Credential_CredentialProvider.CREDENTIAL_PROVIDER_UNSPECIFIED;
      })
      .map((provider) => ({
        label: providerConfig[provider].label,
        value: provider.toString(),
      }));
    return [allOption, ...providerOptions];
  }, []);

  const handleLoadCredentials = useCallback(() => {
    setApiLoading(true);
    if (query) {
      query
        .listCredentials(
          create(ListCredentialsRequestSchema, {
            provider: providerFilter ?? Credential_CredentialProvider.CREDENTIAL_PROVIDER_UNSPECIFIED,
          })
        )
        .then((result) => {
          setCredentials(result.credentials);
          setApiLoading(false);
        })
        .catch(() => {
          setApiLoading(false);
        });
    }
  }, [query, providerFilter]);

  useEffect(() => {
    if (query) {
      handleLoadCredentials();
    }
  }, [query, handleLoadCredentials]);

  const handleRefresh = useCallback(() => {
    handleLoadCredentials();
  }, [handleLoadCredentials]);

  const handleProviderFilterChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    const value = event.target.value;
    if (value === '') {
      setProviderFilter(undefined);
    } else {
      const provider = parseInt(value, 10) as Credential_CredentialProvider;
      setProviderFilter(provider);
    }
  }, []);

  // Drawer handlers
  const handleOpenDrawer = useCallback(
    async (mode: DrawerMode, credentialSummary?: CredentialSummary) => {
      setDrawerMode(mode);
      if (mode === 'create') {
        setSelectedCredential(null);
        setDrawerOpen(true);
      } else if (credentialSummary && query) {
        // Fetch full credential for view/edit
        try {
          const credential = await query.getById(credentialSummary.id);
          setSelectedCredential(credential);
          setDrawerOpen(true);
        } catch (err) {
          console.error('Failed to load credential:', err);
        }
      }
    },
    [query]
  );

  const handleCloseDrawer = useCallback(() => {
    setDrawerOpen(false);
    setDrawerMode(null);
    setSelectedCredential(null);
  }, []);

  const handleSaveSuccess = useCallback(() => {
    handleLoadCredentials();
    handleCloseDrawer();
  }, [handleLoadCredentials, handleCloseDrawer]);

  const handleConfirmDelete = useCallback((row: CredentialSummary) => {
    setCredentialToDelete(row);
    setDeleteDialogOpen(true);
  }, []);

  const handleDelete = useCallback(() => {
    if (command && credentialToDelete) {
      command
        .delete(credentialToDelete.id)
        .then(() => {
          handleLoadCredentials();
          setDeleteDialogOpen(false);
          setCredentialToDelete(null);
        })
        .catch(() => {
          setDeleteDialogOpen(false);
          setCredentialToDelete(null);
        });
    }
  }, [command, credentialToDelete, handleLoadCredentials]);

  const handleCancelDelete = useCallback(() => {
    setDeleteDialogOpen(false);
    setCredentialToDelete(null);
  }, []);

  const handleCreateClick = useCallback(() => {
    handleOpenDrawer('create');
  }, [handleOpenDrawer]);

  const tableHeaders = useMemo(
    () => ['Name', 'Provider', 'Created At', 'Updated At', 'Actions'],
    []
  );
  const tableDataPath = useMemo<Array<'name' | 'provider' | 'createdAt' | 'updatedAt'>>(
    () => ['name', 'provider', 'createdAt', 'updatedAt'],
    []
  );

  const tableDataTemplate = useMemo(
    () => ({
      provider: (val: string, row: CredentialSummary) => {
        const provider = row.provider as Credential_CredentialProvider;
        return providerConfig[provider]?.label || 'Unknown';
      },
      createdAt: (val: string, row: CredentialSummary) => {
        return row.createdAt ? formatTimestampToDate(row.createdAt, 'DD/MM/YYYY') : '-';
      },
      updatedAt: (val: string, row: CredentialSummary) => {
        return row.updatedAt ? formatTimestampToDate(row.updatedAt, 'DD/MM/YYYY') : '-';
      },
    }),
    []
  );

  const tableActions: ActionMenuProps<CredentialSummary>[] = useMemo(
    () => [
      {
        text: 'View',
        handler: (row: CredentialSummary) => {
          handleOpenDrawer('view', row);
        },
      },
      {
        text: 'Edit',
        handler: (row: CredentialSummary) => {
          handleOpenDrawer('edit', row);
        },
      },
      {
        text: 'Delete',
        handler: (row: CredentialSummary) => {
          handleConfirmDelete(row);
        },
      },
    ],
    [handleOpenDrawer, handleConfirmDelete]
  );

  return (
    <>
      <Box>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h5">Credentials</Typography>
          <Stack flexDirection="row" gap={1}>
            <SimpleSelect
              value={
                providerFilter?.toString() ||
                Credential_CredentialProvider.CREDENTIAL_PROVIDER_UNSPECIFIED.toString()
              }
              onChange={handleProviderFilterChange}
              fullWidth={false}
              options={providerFilterOptions}
              sx={{ minWidth: 250 }}
            />
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
          data={credentials}
          loading={apiLoading}
          options={{
            headers: tableHeaders,
            dataPath: tableDataPath,
            dataTemplate: tableDataTemplate,
            actions: tableActions,
            showPagination: false,
            emptyMessage: 'No credentials found',
            border: true,
          }}
        />
      </Box>

      {/* Drawer for View/Edit/Create */}
      <CredentialDrawer
        open={drawerOpen}
        mode={drawerMode}
        onClose={handleCloseDrawer}
        onSaveSuccess={handleSaveSuccess}
        selectedCredential={selectedCredential}
      />

      {/* Delete Confirmation Dialog */}
      <AlertDialog
        open={deleteDialogOpen}
        onClose={handleCancelDelete}
        onSubmit={handleDelete}
        title="Delete Credential"
        subTitle={`Are you sure you want to delete "${credentialToDelete?.name}"? This action cannot be undone.`}
        submitLabel="Delete"
        submitBtnColor="error"
        cancelLabel="Cancel"
      />
    </>
  );
}
