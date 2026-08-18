'use client';

import { Box } from '@mui/material';
import { styled } from '@mui/material/styles';

export const StyledSwiperContainer = styled(Box)(() => ({
  '& .swiper-pagination-bullet': {
    backgroundColor: 'white !important',
  },

  '& .swiper-pagination-bullet-active': {
    backgroundColor: '#03A587 !important',
  },
}));
