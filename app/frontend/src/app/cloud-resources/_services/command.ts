import { useContext, useEffect, useMemo, useState } from 'react';
import { create } from '@bufbuild/protobuf';
// Connect RPC clients accept messages directly, no wrapping needed
import { AppContext } from '@/contexts';
import { useConnectRpcClient } from '@/hooks';
import { CloudResourceCommandController } from '@/gen/org/openmcf/app/cloudresource/v1/command_pb';
import {
  CreateCloudResourceRequestSchema,
  CreateCloudResourceResponse,
  UpdateCloudResourceRequestSchema,
  UpdateCloudResourceResponse,
  DeleteCloudResourceRequestSchema,
  DeleteCloudResourceResponse,
} from '@/gen/org/openmcf/app/cloudresource/v1/io_pb';
import { CloudResource } from '@/gen/org/openmcf/app/cloudresource/v1/api_pb';

interface CommandType {
  create: (manifest: string) => Promise<CloudResource>;
  update: (id: string, manifest: string) => Promise<CloudResource>;
  delete: (id: string) => Promise<void>;
}

const RESOURCE_NAME = 'Cloud Resource';

export const useCloudResourceCommand = () => {
  const { setPageLoading, openSnackbar } = useContext(AppContext);
  const commandClient = useConnectRpcClient(CloudResourceCommandController);
  const [command, setCommand] = useState<CommandType>(null);

  const commandApis: CommandType = useMemo(
    () => ({
      create: (manifest: string): Promise<CloudResource> => {
        setPageLoading(true);
        return new Promise((resolve, reject) => {
          commandClient
            .create(create(CreateCloudResourceRequestSchema, { manifest }))
            .then((response: CreateCloudResourceResponse) => {
              if (response?.resource) {
                openSnackbar(
                  `${RESOURCE_NAME} ${response.resource.name} created successfully`,
                  'success'
                );
                resolve(response.resource);
              } else {
                reject(new Error('No resource returned from create operation'));
              }
            })
            .catch((err) => {
              openSnackbar(err.message || `Could not create ${RESOURCE_NAME}`, 'error');
              reject(err);
            })
            .finally(() => setPageLoading(false));
        });
      },
      update: (id: string, manifest: string): Promise<CloudResource> => {
        return new Promise((resolve, reject) => {
          setPageLoading(true);
          commandClient
            .update(create(UpdateCloudResourceRequestSchema, { id, manifest }))
            .then((response: UpdateCloudResourceResponse) => {
              if (response?.resource) {
                openSnackbar(
                  `${RESOURCE_NAME} ${response.resource.name} updated successfully`,
                  'success'
                );
                resolve(response.resource);
              } else {
                reject(new Error('No resource returned from update operation'));
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
          commandClient
            .delete(create(DeleteCloudResourceRequestSchema, { id }))
            .then((response: DeleteCloudResourceResponse) => {
              openSnackbar(response.message || `${RESOURCE_NAME} deleted successfully`, 'success');
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
    [commandClient, openSnackbar, setPageLoading]
  );

  useEffect(() => {
    if (commandClient && !command) {
      setCommand(commandApis);
    }
  }, [commandApis, commandClient, command]);

  return { command };
};
