'use client';

import { useEffect, useRef, useState } from 'react';
import { Box } from '@mui/material';
import { motion } from 'framer-motion';

interface FlipEffectProps {
  initialRotation: number; // Starting rotation angle
  src: string; // URL of the image
  useIframe?: boolean;
  width?: string | number; // Width of the image container
  height?: string | number; // Height of the image container
  rotationOrigin?: string; // CSS transform-origin value (e.g., "left center", "right bottom")
  extraHeight?: number;
}

export const FlipEffect: React.FC<FlipEffectProps> = ({
  initialRotation,
  src,
  useIframe = false,
  width = 300,
  height = 200,
  extraHeight = 0,
  rotationOrigin = 'center center',
}) => {
  const sectionRef = useRef<HTMLDivElement | null>(null);
  const [rotation, setRotation] = useState(0);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    const section = sectionRef.current;
    if (!section) return;

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true); // Start tracking when section enters viewport
        } else {
          setIsVisible(false); // Stop tracking when section leaves viewport
          setRotation(initialRotation); // Reset rotation
        }
      },
      { threshold: 0 } // Trigger when section starts entering/exiting viewport
    );

    observer.observe(section);

    return () => observer.disconnect();
  }, [initialRotation]);

  useEffect(() => {
    if (!isVisible) return; // Only track scroll when section is visible

    const handleScroll = () => {
      if (!sectionRef.current) return;

      const section = sectionRef.current;
      const sectionHeight = section.clientHeight + extraHeight;
      const sectionTop = section.getBoundingClientRect().top;
      const visiblePercentage = Math.min(
        1,
        Math.max(0, (window.innerHeight - sectionTop) / sectionHeight)
      );
      const newRotation = initialRotation - initialRotation * visiblePercentage;
      setRotation(newRotation);
    };

    const onScroll = () => requestAnimationFrame(handleScroll);

    window.addEventListener('scroll', onScroll);
    handleScroll(); // Run once on mount when visible

    return () => window.removeEventListener('scroll', onScroll);
  }, [isVisible, extraHeight, initialRotation]);

  return (
    <Box
      sx={{
        perspective: '500px',
      }}
      className="flex flex-row justify-center items-center"
    >
      <Box
        component={motion.div}
        ref={sectionRef}
        sx={{
          width,
          height,
          position: 'relative',
          transformStyle: 'preserve-3d',
          transformOrigin: rotationOrigin,
        }}
        style={{ transform: `rotateX(${rotation}deg)` }}
      >
        {useIframe ? (
          <iframe
            width={width}
            height={height}
            src={src}
            title="YouTube video player"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
            allowFullScreen={true}
            style={{ border: 'none' }}
          />
        ) : (
          <video src={src} className="absolute w-full h-full object-cover" autoPlay muted loop />
        )}
      </Box>
    </Box>
  );
};
