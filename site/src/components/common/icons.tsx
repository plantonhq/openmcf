'use client';
import React, { FC } from 'react';
import { SvgIcon, SvgIconProps } from '@mui/material';

import videoPlayIcon from '../../../public/_site/images/features/planton-copilot/video-play.svg';

interface IconInterface {
  [name: string]: { icon: FC; color?: string; fill?: string; viewBox?: string };
}

export enum ACTION_ICONS {
  PLAY_VIDEO = 'PLAY_VIDEO',
}

const icons: IconInterface = {
  [ACTION_ICONS.PLAY_VIDEO]: { icon: videoPlayIcon, viewBox: '0 0 16 16' },
};

export interface IActionIcon extends SvgIconProps {
  name: string;
  alt?: string;
}

export const ActionIcon: FC<IActionIcon> = ({ name, sx, ...props }) => {
  const attr = icons[name];

  if (!attr) return <></>;

  const { icon: Icon, color, fill = 'none', viewBox = '0 0 20 20' } = attr;

  return (
    <SvgIcon
      component={Icon}
      sx={{ fontSize: 16, color, fill, cursor: 'pointer', ...sx }}
      viewBox={viewBox}
      {...props}
    />
  );
};
