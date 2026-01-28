'use client';

import { useEffect, useCallback, useMemo } from 'react';
import { useForm, useWatch } from 'react-hook-form';
import { Drawer } from '@/components/shared/drawer';
import { Stack, Button } from '@mui/material';
import {
  DrawerContainer,
  DrawerContentArea,
  DrawerFooter,
} from '@/app/credentials/_components/styled';
import { SimpleInput } from '@/components/shared/simple-input';
import { SimpleSelect } from '@/components/shared/simple-select';
import {
  CredentialFormData,
  Auth0CredentialForm,
  GcpCredentialForm,
  AwsCredentialForm,
  AzureCredentialForm,
} from '@/app/credentials/_components/forms';
import { useCredentialCommand } from '@/app/credentials/_services';
import {
  Credential_CredentialProvider,
  Credential,
  CredentialProviderConfigSchema,
} from '@/gen/org/openmcf/app/credential/v1/api_pb';
import { CreateCredentialRequest } from '@/gen/org/openmcf/app/credential/v1/io_pb';
import { Auth0ProviderConfig, Auth0ProviderConfigSchema } from '@/gen/org/openmcf/provider/auth0/provider_pb';
import { GcpProviderConfig, GcpProviderConfigSchema } from '@/gen/org/openmcf/provider/gcp/provider_pb';
import { AwsProviderConfig, AwsProviderConfigSchema } from '@/gen/org/openmcf/provider/aws/provider_pb';
import { AzureProviderConfig, AzureProviderConfigSchema } from '@/gen/org/openmcf/provider/azure/provider_pb';
import { create } from '@bufbuild/protobuf';
import { providerConfig } from '@/app/credentials/_components/utils';

export type DrawerMode = 'view' | 'edit' | 'create' | null;

interface CredentialDrawerProps {
  open: boolean;
  mode: DrawerMode;
  onClose: () => void;
  onSaveSuccess: () => void;
  selectedCredential?: Credential | null;
  initialProvider?: Credential_CredentialProvider;
}

export function CredentialDrawer({
  open,
  mode,
  onClose,
  onSaveSuccess,
  selectedCredential,
  initialProvider,
}: CredentialDrawerProps) {
  const { command } = useCredentialCommand();
  const isView = mode === 'view';
  const submitLabel = mode === 'edit' ? 'Update' : 'Create';

  const { register, handleSubmit, reset, setValue, control, watch } = useForm<CredentialFormData>({
    defaultValues: {
      name: '',
      provider: Credential_CredentialProvider.CREDENTIAL_PROVIDER_UNSPECIFIED,
      auth0: {},
      gcp: {},
      aws: {},
      azure: {},
    },
  });

  useEffect(() => {
    if (initialProvider) {
      setValue('provider', initialProvider);
    }
  }, [initialProvider, setValue]);

  const formProvider = useWatch({ control, name: 'provider' });

  const providerOptions = useMemo(() => {
    return (Object.keys(providerConfig) as unknown as Array<Credential_CredentialProvider>)
      .filter((provider) => {
        // Filter out UNSPECIFIED (value 0) by comparing numeric enum values
        return Number(provider) !== Credential_CredentialProvider.CREDENTIAL_PROVIDER_UNSPECIFIED;
      })
      .map((provider) => ({
        label: providerConfig[provider].label,
        value: provider,
      }));
  }, []);

  // Populate form when selectedCredential changes
  useEffect(() => {
    if (selectedCredential && (mode === 'view' || mode === 'edit')) {
      const providerConfigData = selectedCredential.providerConfig;
      const formData: CredentialFormData = {
        name: selectedCredential.name,
        provider: selectedCredential.provider,
        auth0: {},
        gcp: {},
        aws: {},
        azure: {},
      };
      if (providerConfigData?.data?.case === 'auth0') {
        formData.auth0 = {
          domain: providerConfigData.data.value.domain,
          clientId: providerConfigData.data.value.clientId,
          clientSecret: providerConfigData.data.value.clientSecret,
        };
      } else if (providerConfigData?.data?.case === 'gcp') {
        formData.gcp = {
          serviceAccountKeyBase64: providerConfigData.data.value.serviceAccountKeyBase64,
        };
      } else if (providerConfigData?.data?.case === 'aws') {
        formData.aws = {
          accountId: providerConfigData.data.value.accountId,
          accessKeyId: providerConfigData.data.value.accessKeyId,
          secretAccessKey: providerConfigData.data.value.secretAccessKey,
          region: providerConfigData.data.value.region,
          sessionToken: providerConfigData.data.value.sessionToken,
        };
      } else if (providerConfigData?.data?.case === 'azure') {
        formData.azure = {
          clientId: providerConfigData.data.value.clientId,
          clientSecret: providerConfigData.data.value.clientSecret,
          tenantId: providerConfigData.data.value.tenantId,
          subscriptionId: providerConfigData.data.value.subscriptionId,
        };
      }
      reset(formData);
    } else if (mode === 'create') {
      reset({
        name: '',
        provider: initialProvider || Credential_CredentialProvider.CREDENTIAL_PROVIDER_UNSPECIFIED,
        auth0: {},
        gcp: {},
        aws: {},
        azure: {},
      });
    }
  }, [selectedCredential, mode, initialProvider, reset]);

  const handleSave = useCallback(
    (formData: CredentialFormData) => {
      if (!command) return;

      let providerConfig: CreateCredentialRequest['providerConfig'];

      if (
        formData.provider == Credential_CredentialProvider.AUTH0 &&
        formData.auth0?.domain &&
        formData.auth0?.clientId &&
        formData.auth0?.clientSecret
      ) {
        providerConfig = create(CredentialProviderConfigSchema, {
          data: {
            case: 'auth0',
            value: create(Auth0ProviderConfigSchema, formData.auth0 as Auth0ProviderConfig),
          },
        });
      } else if (
        formData.provider == Credential_CredentialProvider.GCP &&
        formData.gcp?.serviceAccountKeyBase64
      ) {
        providerConfig = create(CredentialProviderConfigSchema, {
          data: {
            case: 'gcp',
            value: create(GcpProviderConfigSchema, formData.gcp as GcpProviderConfig),
          },
        });
      } else if (
        formData.provider == Credential_CredentialProvider.AWS &&
        formData.aws?.accountId &&
        formData.aws?.accessKeyId &&
        formData.aws?.secretAccessKey
      ) {
        providerConfig = create(CredentialProviderConfigSchema, {
          data: {
            case: 'aws',
            value: create(AwsProviderConfigSchema, formData.aws as AwsProviderConfig),
          },
        });
      } else if (
        formData.provider == Credential_CredentialProvider.AZURE &&
        formData.azure?.clientId &&
        formData.azure?.clientSecret &&
        formData.azure?.tenantId &&
        formData.azure?.subscriptionId
      ) {
        providerConfig = create(CredentialProviderConfigSchema, {
          data: {
            case: 'azure',
            value: create(AzureProviderConfigSchema, formData.azure as AzureProviderConfig),
          },
        });
      } else {
        return;
      }

      if (mode === 'create') {
        command.create(formData.name, formData.provider, providerConfig).then(() => {
          onSaveSuccess();
        });
      } else if (mode === 'edit' && selectedCredential) {
        command
          .update(selectedCredential.id, formData.name, formData.provider, providerConfig)
          .then(() => {
            onSaveSuccess();
          });
      }
    },
    [command, mode, selectedCredential, onSaveSuccess]
  );

  const handleClose = () => {
    reset({
      name: '',
      provider: initialProvider || Credential_CredentialProvider.CREDENTIAL_PROVIDER_UNSPECIFIED,
      auth0: {},
      gcp: {},
      aws: {},
      azure: {},
    });
    onClose();
  };

  const onProviderChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      if (isView || initialProvider) return;
      const newProvider = parseInt(e.target.value, 10) as Credential_CredentialProvider;
      setValue('provider', newProvider);
      setValue('auth0', {});
      setValue('gcp', {});
      setValue('aws', {});
      setValue('azure', {});
    },
    [setValue, isView, initialProvider]
  );

  const title =
    mode === 'view'
      ? 'View Credential'
      : mode === 'edit'
        ? 'Edit Credential'
        : initialProvider
          ? `Create ${providerConfig[initialProvider].label} Credential`
          : 'Create Credential';

  return (
    <Drawer open={open} onClose={handleClose} title={title} width={600}>
      <DrawerContainer>
        <DrawerContentArea $hasFooter={!isView}>
          <Stack spacing={3}>
            <Stack>
              <SimpleSelect
                name="Provider"
                value={formProvider}
                required
                disabled={isView || !!initialProvider}
                onChange={onProviderChange}
                options={providerOptions}
                sx={{ minWidth: 250 }}
              />
              {!!formProvider && (
                <SimpleInput
                  register={register}
                  path="name"
                  name="Name"
                  registerOptions={{ required: true }}
                  disabled={isView}
                />
              )}
              {formProvider == Credential_CredentialProvider.AUTH0 && (
                <Auth0CredentialForm register={register} disabled={isView} />
              )}
              {formProvider == Credential_CredentialProvider.GCP && (
                <GcpCredentialForm
                  setValue={setValue}
                  watch={watch}
                  disabled={isView}
                  credentialName={selectedCredential?.name}
                />
              )}
              {formProvider == Credential_CredentialProvider.AWS && (
                <AwsCredentialForm register={register} disabled={isView} />
              )}
              {formProvider == Credential_CredentialProvider.AZURE && (
                <AzureCredentialForm register={register} disabled={isView} />
              )}
            </Stack>
          </Stack>
        </DrawerContentArea>
        {!isView && (
          <DrawerFooter>
            <Button variant="contained" color="secondary" onClick={handleClose}>
              Cancel
            </Button>
            <Button variant="contained" color="primary" onClick={handleSubmit(handleSave)}>
              {submitLabel}
            </Button>
          </DrawerFooter>
        )}
      </DrawerContainer>
    </Drawer>
  );
}
