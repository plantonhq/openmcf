'use client';

import Link from 'next/link';
import { Stack, Typography } from '@mui/material';
import type { MenuItem } from '../../data/navigation';

export function MegaMenuItem({ label, subLabel, icon, href, onClick }: MenuItem & { onClick?: () => void }) {
  return (
    <Link href={href ?? ''} onClick={onClick} style={{ width: '100%', textDecoration: 'none', color: 'inherit' }}>
      <Stack
        direction="row"
        sx={{
          gap: 2,
          alignItems: { md: 'center' },
          borderRadius: 2,
          px: 1,
          py: 0.75,
          mx: -1,
          transition: 'background-color 150ms ease',
          '&:hover': { bgcolor: 'rgba(255,255,255,0.1)' },
        }}
      >
        {icon}
        <Stack sx={{ justifyContent: 'flex-start' }}>
          <Typography sx={{ color: '#fff', fontWeight: subLabel ? 600 : 400, fontSize: '0.875rem' }}>
            {label}
          </Typography>
          {subLabel && (
            <Typography
              sx={{
                color: '#666',
                fontSize: '0.875rem',
                fontWeight: 400,
                transition: 'color 150ms ease',
                '.MuiStack-root:hover &': { color: '#a0a0a0' },
              }}
            >
              {subLabel}
            </Typography>
          )}
        </Stack>
      </Stack>
    </Link>
  );
}
