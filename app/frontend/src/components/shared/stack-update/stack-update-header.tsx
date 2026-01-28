'use client';

import { Typography, Divider, Stack, Box, Skeleton } from '@mui/material';
import { HeaderContainer, TopSection, BottomSection, FlexCenterRow } from '@/app/stack-update/styled';
import { StackUpdate } from '@/gen/org/openmcf/app/stackupdate/v1/api_pb';
import { StatusChip } from '@/components/shared/status-chip';
import { TextCopy } from '@/components/shared/text-copy';

interface StackUpdateHeaderProps {
  stackUpdate: StackUpdate | null;
  updatedTime: string;
}

export function StackUpdateHeader({ stackUpdate, updatedTime }: StackUpdateHeaderProps) {
  return (
    <HeaderContainer>
      <TopSection>
        <Stack gap={1}>
          <FlexCenterRow gap={0.5}>
            {stackUpdate ? (
              <Typography variant="caption" color="text.secondary">
                {stackUpdate?.id}
              </Typography>
            ) : (
              <Skeleton variant="text" width={180} height={15} />
            )}

            {stackUpdate ? (
              <TextCopy text={stackUpdate?.id} />
            ) : (
              <Skeleton variant="rectangular" width={10} height={10} />
            )}
          </FlexCenterRow>
          <Box>
            {stackUpdate ? (
              <StatusChip status={stackUpdate?.status || 'unknown'} />
            ) : (
              <Skeleton variant="rounded" width={80} height={20} />
            )}
          </Box>
        </Stack>
      </TopSection>
      <Divider sx={{ marginY: 1.5 }} />
      <BottomSection>
        <Box />
        {stackUpdate ? (
          <Typography variant="subtitle2" color="text.secondary">
            {updatedTime}
          </Typography>
        ) : (
          <Skeleton variant="text" width={180} height={15} />
        )}
      </BottomSection>
    </HeaderContainer>
  );
}
