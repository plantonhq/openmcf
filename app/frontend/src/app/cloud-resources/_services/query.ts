import { useContext, useEffect, useMemo, useState } from 'react';
import { create } from '@bufbuild/protobuf';
// Connect RPC clients accept messages directly, no wrapping needed
import { AppContext } from '@/contexts';
import { useConnectRpcClient } from '@/hooks';
import { CloudResourceQueryController } from '@/gen/org/openmcf/app/cloudresource/v1/query_pb';
import {
  ListCloudResourcesRequest,
  ListCloudResourcesResponse,
  GetCloudResourceRequestSchema,
} from '@/gen/org/openmcf/app/cloudresource/v1/io_pb';
import { CloudResource } from '@/gen/org/openmcf/app/cloudresource/v1/api_pb';

interface QueryType {
  listCloudResources: (input: ListCloudResourcesRequest) => Promise<ListCloudResourcesResponse>;
  getById: (id: string) => Promise<CloudResource>;
}

const RESOURCE_NAME = 'Cloud Resources';

export const useCloudResourceQuery = () => {
  const { setPageLoading, openSnackbar } = useContext(AppContext);
  const queryClient = useConnectRpcClient(CloudResourceQueryController);
  const [query, setQuery] = useState<QueryType>(null);

  const cloudResourceQuery: QueryType = useMemo(
    () => ({
      listCloudResources: (
        input: ListCloudResourcesRequest
      ): Promise<ListCloudResourcesResponse> => {
        return new Promise((resolve, reject) => {
          setPageLoading(true);
          queryClient
            .list(input)
            .then((response) => {
              resolve(response);
            })
            .catch((err) => {
              openSnackbar(err.message || `Could not get ${RESOURCE_NAME}!`, 'error');
              reject(err);
            })
            .finally(() => {
              setPageLoading(false);
            });
        });
      },
      getById: (id: string): Promise<CloudResource> => {
        return new Promise((resolve, reject) => {
          setPageLoading(true);
          queryClient
            .get(create(GetCloudResourceRequestSchema, { id }))
            .then((response) => {
              resolve(response.resource!);
            })
            .catch((err) => {
              openSnackbar(err.message || `Could not get ${RESOURCE_NAME}!`, 'error');
              reject(err);
            })
            .finally(() => {
              setPageLoading(false);
            });
        });
      },
    }),
    [queryClient, setPageLoading, openSnackbar]
  );

  useEffect(() => {
    if (queryClient && !query) {
      setQuery(cloudResourceQuery);
    }
  }, [cloudResourceQuery]);

  return { query };
};
