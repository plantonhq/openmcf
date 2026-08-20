'use client';
import { ReactNode } from 'react';
import { Swiper, SwiperSlide } from 'swiper/react';
import { Pagination } from 'swiper/modules';
import 'swiper/css';
import 'swiper/css/pagination';
import { StyledSwiperContainer } from './styled';

interface ISlider {
  items: ReactNode[];
}

export const Slider: React.FC<ISlider> = ({ items }) => {
  return (
    <StyledSwiperContainer>
      <Swiper slidesPerView={1} modules={[Pagination]} pagination={{ clickable: true }}>
        {items.map((item, index) => (
          <SwiperSlide key={index}>{item}</SwiperSlide>
        ))}
      </Swiper>
    </StyledSwiperContainer>
  );
};
