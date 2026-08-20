export * from './icons';
export * from './typography';
export * from './flip-box';
export * from './bg-dots';
export * from './carousel-section';
export * from './plantora-grid';
export * from './flip-effect';
export * from './slider';
export * from './images';
export * from './content-layout';
export * from './content-sidebar';
export * from './content-details-sidebar';
export * from './CodeBlock';
export { default as MermaidDiagram } from './MermaidDiagram';
export { PageActions } from './PageActions';

export const randomBetween = (min: number, max: number): number => {
  return Math.floor(Math.random() * (max - min + 1)) + min;
};
