'use client';
import { useState, useEffect, useCallback } from 'react';
import { create } from '@bufbuild/protobuf';
import { Box, Grid2, Stack } from '@mui/material';
import { Cloud as CloudIcon, Key as KeyIcon, Storage as StorageIcon } from '@mui/icons-material';
import { StyledGrid2 } from '@/app/dashboard/styled';
import { StatCard } from '@/components/shared/stat-card';
import { CloudResourcesList } from '@/components/shared/cloud-resources-list';
import { useCloudResourceQuery } from '@/app/cloud-resources/_services';
import { ListCloudResourcesRequestSchema } from '@/gen/org/openmcf/app/cloudresource/v1/io_pb';
import { SectionHeader } from '@/components/shared/section-header';

export default function DashboardPage() {
  const { query } = useCloudResourceQuery();
  const [cloudResourceCount, setCloudResourceCount] = useState<number | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const getCloudResourceCount = useCallback(() => {
    if (query) {
      setIsLoading(true);
      query
        .listCloudResources(create(ListCloudResourcesRequestSchema, {}))
        .then((response) => {
          setCloudResourceCount(response.resources.length);
          setIsLoading(false);
        })
        .catch(() => {
          setIsLoading(false);
        });
    }
  }, [query]);

  useEffect(() => {
    getCloudResourceCount();
  }, [getCloudResourceCount]);

  return (
    <Stack height={'100%'}>
      <SectionHeader
        title="Dashboard"
        borderBottom
        containerProps={{ paddingX: 4, paddingY: 3 }}
      />
      <Box bgcolor="grey.20" height={'100%'} px={4}>
        <StyledGrid2 container spacing={2.5}>
          <Grid2 size={{ xs: 12, sm: 6, md: 4 }}>
            <StatCard
              title="Cloud Resources"
              value={cloudResourceCount}
              icon={<CloudIcon />}
              loading={isLoading}
              href="/cloud-resources"
              accent
            />
          </Grid2>
          <Grid2 size={{ xs: 12, sm: 6, md: 4 }}>
            <StatCard
              title="Credentials"
              value="—"
              icon={<KeyIcon />}
              href="/credentials"
            />
          </Grid2>
          <Grid2 size={{ xs: 12, sm: 6, md: 4 }}>
            <StatCard
              title="Stack Updates"
              value="—"
              icon={<StorageIcon />}
              href="/stack-updates"
            />
          </Grid2>
        </StyledGrid2>
        <Box mt={3}>
          <CloudResourcesList
            showErrorAlerts={true}
            onChange={getCloudResourceCount}
          />
        </Box>
      </Box>
    </Stack>
  );
}
