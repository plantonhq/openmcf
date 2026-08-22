'use client';

import { useState, useEffect, useCallback } from 'react';
import { AwsCredential } from '../components/demo/interfaces';

interface UseAwsCredentialDataReturn {
  data: AwsCredential | null;
  loading: boolean;
  error: string | null;
  refetch: () => void;
}

export const useAwsCredentialData = (): UseAwsCredentialDataReturn => {
  const [data, setData] = useState<AwsCredential | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError(null);
    
    try {
      // Add cache-busting parameter to ensure fresh data
      const timestamp = new Date().getTime();
      const response = await fetch(`/demo/aws-credential.json?t=${timestamp}`);
      
      if (response.ok) {
        const jsonData = await response.json();
        
        // Transform the JSON data to match our AwsCredential interface
        const transformedData: AwsCredential = {
          metadata: {
            name: jsonData.metadata.name || '',
            slug: jsonData.metadata.slug || '',
            id: jsonData.metadata.id || '',
            org: jsonData.metadata.org || '',
            env: '',
            labels: jsonData.metadata.labels || {},
            annotations: jsonData.metadata.annotations || {},
            tags: jsonData.metadata.tags || [],
          },
          spec: {
            accountId: jsonData.spec.accountId || '',
            accessKeyId: jsonData.spec.accessKeyId || '',
            secretAccessKey: jsonData.spec.secretAccessKey || '',
            region: jsonData.spec.region || '',
          },
        };
        
        setData(transformedData);
      } else {
        throw new Error(`Failed to fetch data: ${response.status}`);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch AWS credential data';
      setError(errorMessage);
      console.error('Failed to fetch dummy data:', err);
      
      // Fallback to hardcoded data on error
      setData({
        metadata: {
          name: 'main',
          slug: 'main',
          id: 'awscred_01jvpb3y83fftq6e530j71wkr7',
          org: 'planton-cloud',
          env: '',
          labels: {},
          annotations: {},
          tags: ['first', 'second'],
        },
        spec: {
          accountId: 'accountId',
          accessKeyId: 'accessKeyId',
          secretAccessKey: 'secretAccessKey',
          region: 'ap-south-2',
        },
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
