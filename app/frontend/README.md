# OpenMCF Web App

A Next.js application with Material-UI, ConnectRPC, and Buf integration.

## Features

- **Next.js 14** - Latest stable version
- **Material-UI v6** - With dark/light theme support
- **Styled Components** - Using Emotion
- **ConnectRPC** - For gRPC communication
- **Buf** - Protocol buffer management

## Getting Started

### Prerequisites

- Node.js 18+ 
- Yarn 3.6.4+
- Buf CLI

### Installation

1. Install dependencies:
```bash
make deps
# or
yarn install
```

2. Generate proto code (if you have buf configured):
```bash
make generate
# or
buf generate
```

3. Run development server:
```bash
make dev
# or
yarn dev
```

4. Open [http://localhost:3000](http://localhost:3000) in your browser.

## Project Structure

```
app/frontend/
├── src/
│   ├── app/              # Next.js app directory
│   ├── components/       # React components
│   ├── contexts/         # React contexts
│   ├── themes/          # MUI theme configuration
│   └── gen/             # Generated proto code (from buf)
├── proto/               # Protocol buffer definitions
├── buf.yaml            # Buf configuration
└── buf.gen.yaml        # Buf code generation config
```

## Theme

The application supports both light and dark themes. Toggle between themes using the theme switcher in the header.

## Buf Integration

The project uses Buf for protocol buffer management:

- `proto/service.proto` - Contains CommandService and QueryService definitions
- Run `buf generate` to generate TypeScript code from proto files
- Generated code will be in `src/gen/`

## Development

- `yarn dev` - Start development server
- `yarn build` - Build for production
- `yarn start` - Start production server

## License

Private

