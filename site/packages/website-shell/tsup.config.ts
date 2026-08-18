import { defineConfig } from 'tsup';

export default defineConfig({
  entry: ['src/index.ts', 'src/theme/index.ts'],
  format: ['esm'],
  outDir: 'dist',
  dts: false,
  sourcemap: true,
  clean: true,
  splitting: false,
  banner: { js: "'use client';" },
  external: [/^[^.]/],
  esbuildOptions(options) {
    options.jsx = 'automatic';
  },
});
