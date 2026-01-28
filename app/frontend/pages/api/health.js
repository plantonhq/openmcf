// Health check endpoint for Docker health checks
export default function handler(req, res) {
  res.status(200).json({
    status: 'ok',
    timestamp: new Date().toISOString(),
    service: 'openmcf-frontend'
  });
}
