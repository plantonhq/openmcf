import { CSSProperties, FC } from 'react';
import { Box, Stack, Typography } from '@mui/material';
import { TypoP1 } from './typography';

interface IFlipBox {
  title: string;
  subTitle: string;
  bgImgSrc: string;
  topImgSrc: string;
  reverse?: boolean;
  showPlantoraFace?: boolean;
  gap?: CSSProperties['gap'];
  maxWidth?: CSSProperties['maxWidth'];
}

export const FlipBox: FC<IFlipBox> = ({
  title,
  subTitle,
  bgImgSrc,
  topImgSrc,
  reverse,
  showPlantoraFace,
  gap = '96px',
  maxWidth = '1225px',
}) => {
  return (
    <Stack
      className={`${reverse ? 'flex-row-reverse' : 'flex-row'} items-center`}
      style={{ gap, maxWidth }}
    >
      <Stack className="gap-4 max-w-[400px]">
        <Typography className="text-[40px] font-semibold leading-[48px]">{title}</Typography>
        <TypoP1 className="text-text-secondary">{subTitle}</TypoP1>
      </Stack>
      <Box className="relative">
        <img src={bgImgSrc} />
        <img src={topImgSrc} className="absolute bottom-0 left-1/2 transform -translate-x-1/2" />
        {showPlantoraFace && (
          <img
            className="absolute right-[calc(0%-20px)] bottom-[calc(0%-60px)]  w-[185px] h-[160px]"
            src={'/_site/images/features/planton-copilot/face-1.gif'}
          />
        )}
      </Box>
    </Stack>
  );
};
