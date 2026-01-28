import { useContext, useEffect, useMemo, useState } from 'react';
import { create } from '@bufbuild/protobuf';
// Connect RPC clients accept messages directly, no wrapping needed
import { AppContext } from '@/contexts';
import { useConnectRpcClient } from '@/hooks';
import { CredentialQueryController } from '@/gen/org/openmcf/app/credential/v1/query_pb';
import {
  ListCredentialsRequest,
  ListCredentialsResponse,
  GetCredentialRequestSchema,
} from '@/gen/org/openmcf/app/credential/v1/io_pb';
import { Credential } from '@/gen/org/openmcf/app/credential/v1/api_pb';

interface QueryType {
  listCredentials: (input: ListCredentialsRequest) => Promise<ListCredentialsResponse>;
  getById: (id: string) => Promise<Credential>;
}

const RESOURCE_NAME = 'Credentials';

export const useCredentialQuery = () => {
  const { setPageLoading, openSnackbar } = useContext(AppContext);
  const queryClient = useConnectRpcClient(CredentialQueryController);
  const [query, setQuery] = useState<QueryType | null>(null);

  const credentialQuery: QueryType = useMemo(
    () => ({
      listCredentials: (input: ListCredentialsRequest): Promise<ListCredentialsResponse> => {
        return new Promise((resolve, reject) => {
          setPageLoading(true);
          queryClient.list(input)
            .then(resolve)
            .catch((err) => {
              openSnackbar(err.message || `Could not get ${RESOURCE_NAME}!`, 'error');
              reject(err);
            })
            .finally(() => {
              setPageLoading(false);
            });
        });
      },
      getById: (id: string): Promise<Credential> => {
        return new Promise((resolve, reject) => {
          setPageLoading(true);
          queryClient.get(create(GetCredentialRequestSchema, { id }))
            .then((response) => {
              resolve(response.credential!);
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
      setQuery(credentialQuery);
    }
  }, [credentialQuery]);

  return { query };
};

