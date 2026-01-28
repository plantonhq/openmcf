'use client';

import { useEffect, useState, useMemo, useCallback, useContext, useRef } from 'react';
import { useParams } from 'next/navigation';
import { Box, Skeleton, Stack, Typography, Paper, Divider, Tabs, Tab } from '@mui/material';
import { StackUpdateContainer, LogContainer, LogEntry } from '@/app/stack-update/styled';
import { useStackUpdateQuery } from '@/app/stack-update/_services';
import { StackUpdate } from '@/gen/org/openmcf/app/stackupdate/v1/api_pb';
import { Breadcrumb, BreadcrumbStartIcon, IBreadcrumbItem } from '@/components/shared/breadcrumb';
import { StackUpdatesDrawer, StackUpdateHeader } from '@/components/shared/stack-update';
import { ICON_NAMES } from '@/components/shared/icon';
import { formatTimestampToDate } from '@/lib';
import { JsonCode } from '@/components/shared/syntax-highlighter';
import { AppContext, THEME } from '@/contexts';

interface StreamingLog {
  sequenceNum: number;
  content: string;
  streamType: 'stdout' | 'stderr';
  timestamp: Date;
}

export default function StackUpdateDetailPage() {
  const { theme } = useContext(AppContext);
  const params = useParams();
  const { query } = useStackUpdateQuery();
  const [stackUpdate, setStackUpdate] = useState<StackUpdate | null>(null);
  const [stackUpdatesDrawerOpen, setStackUpdatesDrawerOpen] = useState(false);
  const [streamingLogs, setStreamingLogs] = useState<StreamingLog[]>([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [tabIndex, setTabIndex] = useState(0);
  const [intervalCount, setIntervalCount] = useState(0);
  const abortControllerRef = useRef<AbortController | null>(null);
  const logContainerRef = useRef<HTMLDivElement | null>(null);
  const hasStartedStreamingRef = useRef<boolean>(false);
  const streamingJobIdRef = useRef<string | null>(null);
  const isStreamingRef = useRef<boolean>(false);
  const hasLoadedLogsRef = useRef<boolean>(false);

  const stackUpdateId = params?.id as string;

  const handleCloseStackUpdates = useCallback(() => {
    setStackUpdatesDrawerOpen(false);
  }, []);

  const handleStackUpdatesClick = useCallback(() => {
    if (stackUpdate?.cloudResourceId) {
      setStackUpdatesDrawerOpen(true);
    }
  }, [stackUpdate?.cloudResourceId]);

  const breadcrumbs: IBreadcrumbItem[] = useMemo(() => {
    const items: IBreadcrumbItem[] = [];

    // Always show the ID from params, even if stackUpdate is not loaded yet
    if (stackUpdateId) {
      items.push({
        name: stackUpdateId,
        handler: undefined, // Last item is not clickable
      });
    }

    return items;
  }, [stackUpdateId, handleStackUpdatesClick]);

  // Auto-scroll to bottom when new logs arrive
  useEffect(() => {
    if (logContainerRef.current) {
      logContainerRef.current.scrollTop = logContainerRef.current.scrollHeight;
    }
  }, [streamingLogs]);

  // Create interval - updates every 5 seconds
  useEffect(() => {
    const intervalId = setInterval(() => {
      setIntervalCount((prevCount) => prevCount + 1);
    }, 5000); // 5 seconds

    return () => clearInterval(intervalId);
  }, []);

  // Fetch stack-update details when intervalCount changes
  useEffect(() => {
    if (!query || !stackUpdateId) {
      return;
    }

    query
      .getById(stackUpdateId)
      .then((job) => {
        setStackUpdate((prevJob) => {
          // Only update if status or output changed to avoid unnecessary re-renders
          if (!prevJob || prevJob.status !== job.status || prevJob.output !== job.output) {
            return job;
          }
          return prevJob;
        });
      })
      .catch((error) => {
        console.error('Failed to fetch stack-update:', error);
      });
  }, [query, stackUpdateId, intervalCount]);

  // Auto-start streaming when job is in progress or completed
  useEffect(() => {
    if (!stackUpdate || !query || !stackUpdateId) {
      return;
    }

    // Reset streaming state if job ID changed
    if (streamingJobIdRef.current && streamingJobIdRef.current !== stackUpdateId) {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
        abortControllerRef.current = null;
      }
      isStreamingRef.current = false;
      setIsStreaming(false);
      hasStartedStreamingRef.current = false;
      hasLoadedLogsRef.current = false;
      streamingJobIdRef.current = null;
    }

    // Don't start if already streaming this job
    if (streamingJobIdRef.current === stackUpdateId || hasStartedStreamingRef.current) {
      return;
    }

    // Start streaming for in-progress jobs or completed jobs (to load existing logs)
    const shouldStartStreaming =
      stackUpdate.status === 'in_progress' ||
      ((stackUpdate.status === 'success' || stackUpdate.status === 'failed') &&
        !hasLoadedLogsRef.current);

    if (!shouldStartStreaming) {
      return;
    }

    // Mark that we're loading logs for completed jobs
    if (stackUpdate.status === 'success' || stackUpdate.status === 'failed') {
      hasLoadedLogsRef.current = true;
    }

    // Mark as started to prevent re-triggering
    streamingJobIdRef.current = stackUpdateId;
    hasStartedStreamingRef.current = true;
    isStreamingRef.current = true;
    setIsStreaming(true);
    setStreamingLogs([]);

    // Create new AbortController
    const abortController = new AbortController();
    abortControllerRef.current = abortController;

    // Start streaming in async function
    const streamAsync = async () => {
      try {
        const stream = query.streamOutput(stackUpdateId, undefined, abortController.signal);

        for await (const response of stream) {
          if (abortController.signal.aborted) {
            break;
          }

          const log: StreamingLog = {
            sequenceNum: response.sequenceNum,
            content: response.content,
            streamType: response.streamType === 'stderr' ? 'stderr' : 'stdout',
            timestamp: response.timestamp
              ? new Date(Number(response.timestamp.seconds) * 1000)
              : new Date(),
          };

          setStreamingLogs((prev) => [...prev, log]);

          // If stream is completed, stop streaming
          if (response.status === 'completed' || response.status === 'cancelled') {
            isStreamingRef.current = false;
            setIsStreaming(false);
            hasStartedStreamingRef.current = false;
            if (streamingJobIdRef.current === stackUpdateId) {
              streamingJobIdRef.current = null;
            }
            break;
          }
        }
      } catch (error: any) {
        if (error.name !== 'AbortError') {
          console.error('Stream error:', error);
          isStreamingRef.current = false;
          setIsStreaming(false);
          hasStartedStreamingRef.current = false;
          if (streamingJobIdRef.current === stackUpdateId) {
            streamingJobIdRef.current = null;
          }
        }
      }
    };

    streamAsync();

    // Cleanup on unmount or when stackUpdateId changes
    return () => {
      if (streamingJobIdRef.current === stackUpdateId) {
        if (abortControllerRef.current) {
          abortControllerRef.current.abort();
          abortControllerRef.current = null;
        }
        isStreamingRef.current = false;
        setIsStreaming(false);
        hasStartedStreamingRef.current = false;
        streamingJobIdRef.current = null;
      }
    };
  }, [stackUpdate?.status, stackUpdateId, query]);

  const updatedTime = stackUpdate?.updatedAt
    ? formatTimestampToDate(stackUpdate.updatedAt, 'DD/MM/YYYY, HH:mm:ss')
    : '-';

  return (
    <StackUpdateContainer>
      <Stack gap={2}>
        <Breadcrumb
          breadcrumbs={breadcrumbs}
          startBreadcrumb={
            <BreadcrumbStartIcon
              icon={ICON_NAMES.INFRA_HUB}
              iconProps={{ sx: { filter: theme.mode === THEME.DARK ? 'invert(1)' : 'none' } }}
              label="Stack Updates"
              handler={handleStackUpdatesClick}
            />
          }
        />
        <StackUpdateHeader stackUpdate={stackUpdate} updatedTime={updatedTime} />

        <Paper sx={{ p: 2 }}>
          <Tabs value={tabIndex} onChange={(_, newValue) => setTabIndex(newValue)}>
            <Tab label="Streaming Output" />
            <Tab label="Final Output" />
          </Tabs>

          <Divider sx={{ my: 2 }} />

          {tabIndex === 0 ? (
            <Box>
              <Box
                sx={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  mb: 2,
                }}
              >
                <Typography variant="h6">Deployment Logs</Typography>
                <Typography variant="body2" color="text.secondary">
                  {streamingLogs.length} log entries
                  {isStreaming && ' • Streaming...'}
                </Typography>
              </Box>

              <LogContainer ref={logContainerRef}>
                {streamingLogs.length === 0 ? (
                  <Typography
                    variant="body2"
                    color="text.secondary"
                    sx={{ textAlign: 'center', py: 4 }}
                  >
                    {isStreaming
                      ? 'Waiting for logs...'
                      : stackUpdate?.status === 'in_progress'
                        ? 'Click to start streaming'
                        : 'No logs available'}
                  </Typography>
                ) : (
                  streamingLogs.map((log, index) => (
                    <LogEntry key={index} streamType={log.streamType}>
                      <Box sx={{ display: 'flex', gap: 1, alignItems: 'flex-start' }}>
                        <Typography
                          variant="caption"
                          color="text.secondary"
                          sx={{ minWidth: 180, fontFamily: 'monospace' }}
                        >
                          {log.timestamp.toLocaleTimeString()}
                        </Typography>
                        <Typography
                          variant="caption"
                          color="text.secondary"
                          sx={{ minWidth: 80, fontFamily: 'monospace' }}
                        >
                          [{log.streamType}]
                        </Typography>
                        <Typography
                          variant="body2"
                          component="pre"
                          sx={{
                            margin: 0,
                            whiteSpace: 'pre-wrap',
                            wordBreak: 'break-word',
                            flex: 1,
                          }}
                        >
                          {log.content}
                        </Typography>
                      </Box>
                    </LogEntry>
                  ))
                )}
              </LogContainer>
            </Box>
          ) : (
            <Box>
              {stackUpdate ? (
                <JsonCode content={stackUpdate?.output || {}} />
              ) : (
                <Skeleton variant="rounded" width={'100%'} height={200} />
              )}
            </Box>
          )}
        </Paper>

        {/* Stack Updates Drawer */}
        {stackUpdate?.cloudResourceId && (
          <StackUpdatesDrawer
            open={stackUpdatesDrawerOpen}
            cloudResourceId={stackUpdate.cloudResourceId}
            onClose={handleCloseStackUpdates}
          />
        )}
      </Stack>
    </StackUpdateContainer>
  );
}
