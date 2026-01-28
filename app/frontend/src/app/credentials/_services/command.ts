import { useContext, useEffect, useMemo, useState } from 'react';
import { create } from '@bufbuild/protobuf';
// Connect RPC clients accept messages directly, no wrapping needed
import { AppContext } from '@/contexts';
import { useConnectRpcClient } from '@/hooks';
import { CredentialCommandController } from '@/gen/org/openmcf/app/credential/v1/command_pb';
import {
  CreateCredentialRequest,
  CreateCredentialRequestSchema,
  CreateCredentialResponse,
  UpdateCredentialRequest,
  UpdateCredentialRequestSchema,
  UpdateCredentialResponse,
  DeleteCredentialRequestSchema,
  DeleteCredentialResponse,
} from '@/gen/org/openmcf/app/credential/v1/io_pb';
import { Credential, Credential_CredentialProvider } from '@/gen/org/openmcf/app/credential/v1/api_pb';

interface CommandType {
  create: (
    name: string,
    provider: Credential_CredentialProvider,
    providerConfig: CreateCredentialRequest['providerConfig']
  ) => Promise<Credential>;
  update: (
    id: string,
    name: string,
    provider: Credential_CredentialProvider,
    providerConfig: UpdateCredentialRequest['providerConfig']
  ) => Promise<Credential>;
  delete: (id: string) => Promise<void>;
}

const RESOURCE_NAME = 'Credential';

export const useCredentialCommand = () => {
  const { setPageLoading, openSnackbar } = useContext(AppContext);
  const commandClient = useConnectRpcClient(CredentialCommandController);
  const [command, setCommand] = useState<CommandType | null>(null);

  const commandApis: CommandType = useMemo(
    () => ({
      create: (
        name: string,
        provider: Credential_CredentialProvider,
        providerConfig: CreateCredentialRequest['providerConfig']
      ): Promise<Credential> => {
        return new Promise((resolve, reject) => {
          setPageLoading(true);
          commandClient.create(
              create(CreateCredentialRequestSchema, {
                name,
                provider,
                providerConfig,
              })
            )
            .then((response: CreateCredentialResponse) => {
              if (response?.credential) {
                openSnackbar(
                  `${RESOURCE_NAME} ${response.credential.name} created successfully`,
                  'success'
                );
                resolve(response.credential);
              } else {
                reject(new Error('No credential returned from create operation'));
              }
            })
            .catch((err) => {
              openSnackbar(err.message || `Could not create ${RESOURCE_NAME}`, 'error');
              reject(err);
            })
            .finally(() => setPageLoading(false));
        });
      },
      update: (
        id: string,
        name: string,
        provider: Credential_CredentialProvider,
        providerConfig: UpdateCredentialRequest['providerConfig']
      ): Promise<Credential> => {
        return new Promise((resolve, reject) => {
          setPageLoading(true);
          commandClient.update(
              create(UpdateCredentialRequestSchema, {
                id,
                name,
                provider,
                providerConfig,
              })
            )
            .then((response: UpdateCredentialResponse) => {
              if (response?.credential) {
                openSnackbar(
                  `${RESOURCE_NAME} ${response.credential.name} updated successfully`,
                  'success'
                );
                resolve(response.credential);
              } else {
                reject(new Error('No credential returned from update operation'));
              }
            })
            .catch((err) => {
              openSnackbar(err.message || `Could not update ${RESOURCE_NAME}`, 'error');
              reject(err);
            })
            .finally(() => setPageLoading(false));
        });
      },
      delete: (id: string): Promise<void> => {
        return new Promise((resolve, reject) => {
          setPageLoading(true);
          commandClient.delete(create(DeleteCredentialRequestSchema, { id }))
            .then((response: DeleteCredentialResponse) => {
              openSnackbar(
                response.message || `${RESOURCE_NAME} deleted successfully`,
                'success'
              );
              resolve();
            })
            .catch((err) => {
              openSnackbar(err.message || `Could not delete ${RESOURCE_NAME}`, 'error');
              reject(err);
            })
            .finally(() => setPageLoading(false));
        });
      },
    }),
    [commandClient, setPageLoading, openSnackbar]
  );

  useEffect(() => {
    if (commandClient && !command) {
      setCommand(commandApis);
    }
  }, [commandApis]);

  return { command };
};

