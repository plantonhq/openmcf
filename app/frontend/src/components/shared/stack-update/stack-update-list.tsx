'use client';

import { useState, useMemo, useCallback, useEffect } from 'react';
import { Typography, Box } from '@mui/material';
import { useRouter } from 'next/navigation';
import { create } from '@bufbuild/protobuf';
import { TableComp } from '@/components/shared/table';
import { PAGINATION_MODE } from '@/models/table';
import { StatusChip } from '@/components/shared/status-chip';
import { useStackUpdateQuery } from '@/app/stack-update/_services';
import { ListStackUpdatesRequestSchema } from '@/gen/org/openmcf/app/stackupdate/v1/io_pb';
import { StackUpdate } from '@/gen/org/openmcf/app/stackupdate/v1/api_pb';
import { PageInfoSchema } from '@/gen/org/openmcf/app/commons/page_info_pb';
import { formatTimestampToDate } from '@/lib';

export interface StackUpdatesListProps {
  cloudResourceId: string;
}

export function StackUpdatesList({ cloudResourceId }: StackUpdatesListProps) {
  const router = useRouter();
  const { query } = useStackUpdateQuery();
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [stackUpdates, setStackUpdates] = useState<StackUpdate[]>([]);
  const [totalPages, setTotalPages] = useState(0);
  const [apiLoading, setApiLoading] = useState(true);

  // Function to call the API
  const handleLoadStackUpdates = useCallback(() => {
    setApiLoading(true);
    if (query) {
      query
        .listStackUpdates(
          create(ListStackUpdatesRequestSchema, {
            cloudResourceId,
            pageInfo: create(PageInfoSchema, {
              num: page,
              size: rowsPerPage,
            }),
          })
        )
        .then((result) => {
          setStackUpdates(result.stackUpdates);
          setTotalPages(result.totalPages || 0);
          setApiLoading(false);
        })
        .catch(() => {
          setApiLoading(false);
        });
    }
  }, [query, cloudResourceId, page, rowsPerPage]);

  // Reset to first page when cloudResourceId changes
  useEffect(() => {
    setPage(0);
  }, [cloudResourceId]);

  // Auto-load on mount and when dependencies change
  useEffect(() => {
    if (query && cloudResourceId) {
      handleLoadStackUpdates();
    }
  }, [query, cloudResourceId, handleLoadStackUpdates]);

  // Handle page change
  const handlePageChange = useCallback((newPage: number, newRowsPerPage: number) => {
    setPage(newPage);
    setRowsPerPage(newRowsPerPage);
  }, []);

  const tableDataTemplate = useMemo(
    () => ({
      id: (val: string, row: StackUpdate) => {
        return <Typography variant="subtitle2">{row.id.substring(0, 12)}...</Typography>;
      },
      status: (val: string, row: StackUpdate) => {
        return <StatusChip status={row.status || 'unknown'} />;
      },
      createdAt: (val: string, row: StackUpdate) => {
        return row.createdAt ? formatTimestampToDate(row.createdAt, 'DD/MM/YYYY HH:mm') : '-';
      },
      updatedAt: (val: string, row: StackUpdate) => {
        return row.updatedAt ? formatTimestampToDate(row.updatedAt, 'DD/MM/YYYY HH:mm') : '-';
      },
    }),
    []
  );

  const clickableColumns = useMemo(
    () => ({
      id: (row: StackUpdate) => {
        router.push(`/stack-update/${row.id}`);
      },
    }),
    [router]
  );

  return (
    <Box>
      <TableComp
        data={stackUpdates}
        loading={apiLoading}
        options={{
          headers: ['ID', 'Status', 'Created At', 'Updated At'],
          dataPath: ['id', 'status', 'createdAt', 'updatedAt'],
          dataTemplate: tableDataTemplate,
          clickableColumns,
          actions: [],
          showPagination: true,
          paginationMode: PAGINATION_MODE.SERVER,
          currentPage: page,
          rowsPerPage: rowsPerPage,
          totalPages: totalPages,
          onPageChange: handlePageChange,
          emptyMessage: 'No stack-updates found',
          border: true,
        }}
      />
    </Box>
  );
}
