'use client';
import { FC, useEffect, useState } from 'react';
import { Box, Grid2 } from '@mui/material';

interface IBgDots {
  count?: number;
  width?: number;
}

interface ISquareBoxes {
  rows?: number;
  columns?: number;
  className?: string;
}

const generateRandomDots = (count: number, screenWidth: number, height: number, width: number) =>
  Array.from({ length: count }).map(() => ({
    top: `${Math.random() * height}px`,
    left: `${Math.random() * width + (screenWidth / 2 - width / 2)}px`,
  }));

export const BgDots: FC<IBgDots> = ({ count = 100, width = 1800 }) => {
  const [dots, setDots] = useState<{ top: string; left: string }[]>([]);

  useEffect(() => {
    const updateDimensions = () => {
      const height = document.body.scrollHeight;
      setDots(generateRandomDots(count, document.body.scrollWidth, height, width));
    };

    updateDimensions();
    window.addEventListener('resize', updateDimensions);

    return () => window.removeEventListener('resize', updateDimensions);
  }, [count, width]);

  return (
    <Box>
      {dots.map((dot, index) => (
        <img
          key={index}
          src="/_site/images/solutions/by-size/growing-teams/dot.svg"
          style={{
            top: dot.top,
            left: dot.left,
          }}
          className="absolute text-blue-500 z-50"
        />
      ))}
    </Box>
  );
};

export const SquareBox: FC = () => {
  return (
    <Grid2 size={1} className="border-r border-b border-[#42424247]">
      <Box className="pt-[100%]" />
    </Grid2>
  );
};

export const SquareBoxes: FC<ISquareBoxes> = ({ columns = 12, rows = 1, className }) => {
  return (
    <Grid2 container columns={columns} className={className}>
      {[...Array(columns * rows)].map((_, index) => (
        <SquareBox key={index} />
      ))}
    </Grid2>
  );
};
