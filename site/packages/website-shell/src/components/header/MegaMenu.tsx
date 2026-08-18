'use client';

import { useState } from 'react';
import Link from 'next/link';
import { Menu, Paper, Stack, Typography } from '@mui/material';
import { NavigateNext, KeyboardArrowDown } from '@mui/icons-material';
import { MegaMenuItem } from './MegaMenuItem';
import type { MenuSection, MenuItem } from '../../data/navigation';

interface MegaMenuProps {
  title: string;
  leftMenu: MenuSection[];
  rightMenu?: MenuSection[];
  footerMenu?: MenuItem;
  leftWidth?: number;
}

export function MegaMenu({
  title,
  leftMenu,
  rightMenu,
  footerMenu,
  leftWidth = 270,
}: MegaMenuProps) {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);

  const handleClick = (event: React.MouseEvent<HTMLDivElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  return (
    <>
      <Stack
        aria-controls={open ? 'mega-menu' : undefined}
        aria-expanded={open ? 'true' : undefined}
        aria-haspopup="true"
        direction="row"
        onClick={handleClick}
        sx={{
          cursor: 'pointer',
          alignItems: 'center',
          '&:hover': { color: '#fff' },
        }}
      >
        <Typography
          sx={{
            fontSize: '0.875rem',
            fontWeight: 500,
            color: 'grey.100',
            transition: 'color 150ms ease',
          }}
        >
          {title}
        </Typography>
        <KeyboardArrowDown
          sx={{
            fontSize: '20px !important',
            fontVariationSettings: "'FILL' 1, 'wght' 500, 'GRAD' 200, 'opsz' 48",
          }}
        />
      </Stack>

      <Menu
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
        MenuListProps={{ sx: { padding: 0 } }}
        slotProps={{
          paper: {
            sx: {
              mt: 1.5,
              backgroundColor: 'transparent',
              backgroundImage: 'none',
              boxShadow: '0 0 0 1px rgba(255,255,255,0.06), 0 8px 30px rgba(0,0,0,0.5)',
            },
          },
        }}
      >
        <Stack
          component={Paper}
          onClick={handleClose}
          sx={{
            gap: 2,
            justifyContent: 'space-between',
            borderRadius: 3,
            bgcolor: '#1a1a1a',
            border: '1px solid #3a3a3a',
          }}
        >
          <Stack direction="row" sx={{ p: 2.5 }}>
            <Stack sx={{ width: leftWidth, mr: 2 }}>
              {leftMenu.map((section, index) => (
                <Stack key={index} sx={{ mb: 3, '&:last-child': { mb: 0 } }}>
                  {section.title && (
                    <Typography
                      sx={{
                        fontWeight: 600,
                        mb: 1.5,
                        color: '#666',
                        fontSize: '0.75rem',
                        textTransform: 'uppercase',
                        letterSpacing: '0.05em',
                      }}
                    >
                      {section.title}
                    </Typography>
                  )}
                  <Stack gap={section.title ? 0.5 : 1}>
                    {section.items.map((item) => (
                      <MegaMenuItem key={item.label} {...item} />
                    ))}
                  </Stack>
                </Stack>
              ))}
            </Stack>

            {rightMenu && (
              <Stack sx={{ width: 170, pl: 2, borderLeft: '1px solid #3a3a3a' }}>
                {rightMenu.map((section, index) => (
                  <Stack key={index} sx={{ mb: 3, '&:last-child': { mb: 0 } }}>
                    {section.title && (
                      <Typography
                        sx={{
                          fontWeight: 600,
                          mb: 1.5,
                          color: '#666',
                          fontSize: '0.75rem',
                          textTransform: 'uppercase',
                          letterSpacing: '0.05em',
                        }}
                      >
                        {section.title}
                      </Typography>
                    )}
                    <Stack gap={section.title ? 0.5 : 1}>
                      {section.items.map((item) => (
                        <MegaMenuItem key={item.label} {...item} />
                      ))}
                    </Stack>
                  </Stack>
                ))}
              </Stack>
            )}
          </Stack>

          {footerMenu && (
            <Link href={footerMenu.href} style={{ textDecoration: 'none', color: 'inherit' }}>
              <Stack
                direction="row"
                sx={{
                  py: 1.25,
                  px: 2.5,
                  gap: 2,
                  justifyContent: 'space-between',
                  cursor: 'pointer',
                  bgcolor: 'rgba(255,255,255,0.05)',
                  '&:hover': { bgcolor: 'rgba(255,255,255,0.1)' },
                  transition: 'background-color 150ms ease',
                  borderRadius: '0 0 12px 12px',
                  borderTop: '1px solid #3a3a3a',
                }}
              >
                <Typography sx={{ fontSize: '0.875rem', fontWeight: 500, color: '#a0a0a0' }}>
                  {footerMenu.label}
                </Typography>
                <NavigateNext sx={{ color: '#666' }} />
              </Stack>
            </Link>
          )}
        </Stack>
      </Menu>
    </>
  );
}
