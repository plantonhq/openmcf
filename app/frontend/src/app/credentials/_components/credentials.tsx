'use client';
import { useState } from 'react';
import { Box, Stack, Tab } from '@mui/material';
import { StyledTabs } from '@/app/credentials/_components/styled';
import { SectionHeader } from '@/components/shared/section-header';
import { TabPanel } from '@/components/shared/tabpanel';
import { Providers } from '@/app/credentials/_components/credentials-tab';
import { Credential_CredentialProvider } from '@/gen/org/openmcf/app/credential/v1/api_pb';
import { CredentialDrawer } from '@/app/credentials/_components/forms';
import { CredentialsList } from '@/components/shared/credentials-list';

export default function Credentials() {
  const [tabIndex, setTabIndex] = useState<'providers' | 'credentials'>('providers');
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [selectedProvider, setSelectedProvider] = useState<Credential_CredentialProvider | null>(null);

  const handleTabChange = (_: React.SyntheticEvent, newValue: 'providers' | 'credentials') => {
    setTabIndex(newValue);
  };

  const handleProviderClick = (provider: Credential_CredentialProvider) => {
    setSelectedProvider(provider);
    setDrawerOpen(true);
  };

  const handleDrawerClose = () => {
    setDrawerOpen(false);
    setSelectedProvider(null);
  };

  const handleCredentialCreated = () => {
    setDrawerOpen(false);
    setSelectedProvider(null);
    setTabIndex('credentials');
  };

  return (
    <Stack height={'100%'}>
      <SectionHeader
        title="Credentials"
        borderBottom
        containerProps={{ paddingX: 4, paddingY: 3 }}
      />
      <Stack gap={3} bgcolor="grey.20" height={'100%'}>
        <Box px={2.5} borderBottom={'1px solid'} borderColor={'grey.60'}>
          <StyledTabs value={tabIndex} onChange={handleTabChange}>
            <Tab label="Providers" id="credentials-tabs-0" value="providers" />
            <Tab label="Credentials" id="credentials-tabs-1" value="credentials" />
          </StyledTabs>
        </Box>
        <Box px={4}>
          <TabPanel value={tabIndex} index="providers">
            <Providers onProviderClick={handleProviderClick} />
          </TabPanel>
          <TabPanel value={tabIndex} index="credentials">
            <CredentialsList />
          </TabPanel>
        </Box>
      </Stack>

      {selectedProvider && (
        <CredentialDrawer
          open={drawerOpen}
          mode="create"
          onClose={handleDrawerClose}
          onSaveSuccess={handleCredentialCreated}
          initialProvider={selectedProvider}
        />
      )}
    </Stack>
  );
}
