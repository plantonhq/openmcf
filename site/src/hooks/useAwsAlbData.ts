import { useState, useEffect, useCallback } from 'react';
import { AwsAlb } from '../components/demo/interfaces';

interface UseAwsAlbDataReturn {
  data: AwsAlb | null;
  loading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

export const useAwsAlbData = (): UseAwsAlbDataReturn => {
  const [data, setData] = useState<AwsAlb | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Add cache-busting parameter to ensure fresh data
      const timestamp = new Date().getTime();
      const response = await fetch(`/demo/aws-alb.json?t=${timestamp}`);
      
      if (!response.ok) {
        throw new Error(`Failed to fetch AWS ALB data: ${response.statusText}`);
      }
      
      const jsonData = await response.json();
      setData(jsonData);
    } catch (err) {
      console.error('Error fetching AWS ALB data:', err);
      setError(err instanceof Error ? err.message : 'Failed to fetch data');
      
      // Set fallback data if fetch fails
      setData({
        apiVersion: 'aws.planton.dev/v1',
        kind: 'AwsAlb',
        metadata: {
          name: 'demo-alb',
          slug: 'demo-alb',
          id: '',
          org: '',
          env: '',
          labels: {},
          annotations: {},
          tags: [],
        },
        spec: {
          subnets: [
            { value: 'subnet-12345678' },
            { value: 'subnet-87654321' }
          ],
          securityGroups: [
            { value: 'sg-12345678' }
          ],
          internal: false,
          deleteProtectionEnabled: true,
          idleTimeoutSeconds: 60,
          dns: {
            enabled: true,
            route53ZoneId: { value: 'Z1234567890ABC' },
            hostnames: ['app.example.com']
          },
          ssl: {
            enabled: true,
            certificateArn: { value: 'arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012' }
          }
        }
      });
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return {
    data,
    loading,
    error,
    refetch: fetchData,
  };
};
