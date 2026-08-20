import { Box, Stack } from '@mui/material';
import { FC } from 'react';

interface IImgBox {
  bgImgSrc: string;
  topImgSrc?: string;
  className?: string;
  topImgClassName?: string;
}

interface IPlantoraGridRow {
  className?: string;
  columns: IImgBox[];
}

export const ImgBox: FC<IImgBox> = ({ bgImgSrc, topImgSrc, className, topImgClassName }) => {
  return (
    <Box
      className={`relative w-full h-full rounded-[28px] overflow-hidden bg-cover bg-center ${
        className ?? ''
      }`}
      sx={{ backgroundImage: `url(${bgImgSrc})` }}
    >
      {topImgSrc && (
        <img
          src={topImgSrc}
          className={`absolute top-14 left-1/2 transform -translate-x-1/2 scale-110 ${
            topImgClassName ?? ''
          }`}
        />
      )}
    </Box>
  );
};

const data1: IPlantoraGridRow[] = [
  {
    columns: [
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-3.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        className: 'ml-[300px] min-w-[630px]',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-2.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        className: 'min-w-[630px]',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-3.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        className: 'min-w-[630px]',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
    ],
  },
  {
    columns: [
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-2.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        className: 'ml-[100px] min-w-[630px]',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-3.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        className: 'min-w-[630px]',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-2.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        className: 'min-w-[630px]',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
    ],
  },
  {
    columns: [
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-2.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        className: 'ml-5 min-w-[630px]',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-3.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        className: 'min-w-[630px]',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-2.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        className: 'min-w-[630px]',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
    ],
  },
];

const data2: IPlantoraGridRow[] = [
  {
    className: 'grid-cols-3 gap-8 h-[210px]',
    columns: [
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-3.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-2.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-3.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
    ],
  },
  {
    className: 'grid-cols-12 gap-8 h-[540px]',
    columns: [
      { bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-3.svg', className: 'col-span-2' },
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-2.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-2.svg',
        className: 'col-span-8',
        topImgClassName: 'left-[calc(50%+13px)] scale-[1.5] top-auto -bottom-24',
      },
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-3.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        className: 'col-span-2',
        topImgClassName: 'left-[calc(50%+13px)] top-auto -bottom-8',
      },
    ],
  },
  {
    className: 'grid-cols-4 gap-8 h-[210px]',
    columns: [
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-3.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-2.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-3.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
      {
        bgImgSrc: '/_site/images/features/planton-copilot/gradient-card-2.svg',
        topImgSrc: '/_site/images/features/planton-copilot/chatbox-1.svg',
        topImgClassName: 'left-[calc(50%+13px)]',
      },
    ],
  },
];

export const PlantoraGridBoxes1: FC = () => {
  return (
    <>
      {data1.map((row, index) => (
        <Stack className="w-full flex-row gap-8 h-[320px] overflow-x-hidden" key={index}>
          {row.columns.map((column, index2) => (
            <ImgBox {...column} key={index2} />
          ))}
        </Stack>
      ))}
    </>
  );
};

export const PlantoraGridBoxes2: FC = () => {
  return (
    <>
      {data2.map((row, index) => (
        <Box className={`w-full grid ${row.className}`} key={index}>
          {row.columns.map((column, index2) => (
            <ImgBox {...column} key={index2} />
          ))}
        </Box>
      ))}
    </>
  );
};
