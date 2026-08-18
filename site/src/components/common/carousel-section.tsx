'use client';

import { FC, useEffect, useMemo, useState } from 'react';

import { Box, Stack } from '@mui/material';
import { TypoH2 } from './typography';
import { StyledGradientBorderBox2 } from '../features/all/styled';

export interface ICarouselLabel {
  label: string;
  width: string;
}

export interface CarouselSectionProps {
  title: string;
  labels: ICarouselLabel[];
}

interface ICarouselItem {
  labels: ICarouselLabel[];
}

export const CarouselItem: FC<ICarouselItem> = ({ labels }) => {
  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => {
      setCurrentIndex((prevIndex) => (prevIndex + 1) % labels.length);
    }, 2000);

    return () => clearInterval(interval);
  }, [labels.length]);

  const animWidth = useMemo(() => labels[currentIndex]?.width ?? '', [currentIndex, labels]);

  return (
    <StyledGradientBorderBox2
      className={`word-carousel h-[150px] px-5 py-8 rounded-3xl`}
      sx={{ width: animWidth, transition: 'width 0.5s ease-in-out' }}
    >
      <Box
        className="word-carousel-container "
        style={{ marginTop: `-${110 * currentIndex + 5}px` }}
      >
        {labels.map(({ label }) => (
          <TypoH2
            key={label}
            className="!leading-tight text-transparent whitespace-nowrap bg-clip-text bg-[linear-gradient(180deg,_#816FF2_0%,_#104EAE_100%)]"
          >
            {label}
          </TypoH2>
        ))}
      </Box>
    </StyledGradientBorderBox2>
  );
};

export const CarouselSection: FC<CarouselSectionProps> = ({ title, labels }) => {
  return (
    <Stack className="flex-row justify-start items-center gap-5 z-10">
      {labels.length > 0 && <CarouselItem labels={labels} />}
      <TypoH2>{title}</TypoH2>
    </Stack>
  );
};
