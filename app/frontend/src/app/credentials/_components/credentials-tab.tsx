'use client';
import { useState, useCallback, useMemo } from 'react';
import { Grid2, InputAdornment } from '@mui/material';
import { Search } from '@mui/icons-material';
import { Credential_CredentialProvider } from '@/gen/org/openmcf/app/credential/v1/api_pb';
import { Icon, ICON_NAMES } from '@/components/shared/icon';
import { FlexCenterRow } from '@/components/shared/resource-header/styled';
import {
  StyledCard,
  StyledCardContent,
  StyledCardHeader,
  StyledChip,
  StyledButton,
  SearchContainer,
  SearchTextField,
  EmptyStateBox,
} from '@/app/credentials/_components/styled';
import { providerConfig } from '@/app/credentials/_components/utils';

export interface ProviderCardProps {
  provider: Credential_CredentialProvider;
  title: string;
  description: string;
  icon?: ICON_NAMES;
  onClick: (provider: Credential_CredentialProvider) => void;
}

export function ProviderCard({ provider, title, description, onClick }: ProviderCardProps) {
  const config = providerConfig[provider];
  const icon = config?.icon;

  return (
    <StyledCard onClick={() => onClick(provider)}>
      <StyledCardContent>
        <FlexCenterRow justifyContent="space-between">
          {icon && <Icon name={icon} sx={{ fontSize: 32 }} />}
          <StyledChip>Cloud Provider</StyledChip>
        </FlexCenterRow>
        <StyledCardHeader title={title} subheader={description} />
        <StyledButton
          variant="contained"
          color="secondary"
          fullWidth
          onClick={(e) => {
            e.stopPropagation();
            onClick(provider);
          }}
        >
          Connect
        </StyledButton>
      </StyledCardContent>
    </StyledCard>
  );
}

interface IProviders {
  onProviderClick: (provider: Credential_CredentialProvider) => void;
}

export function Providers({ onProviderClick }: IProviders) {
  const [searchValue, setSearchValue] = useState('');

  const allProviders = useMemo(() => {
    return (Object.keys(providerConfig) as unknown as Array<Credential_CredentialProvider>)
      .filter((provider) => {
        // Filter out UNSPECIFIED (value 0) by comparing numeric enum values
        return Number(provider) !== Credential_CredentialProvider.CREDENTIAL_PROVIDER_UNSPECIFIED;
      })
      .map((provider) => {
        const config = providerConfig[provider];
        return {
          provider,
          title: config.label,
          description: config.description,
        };
      });
  }, []);

  const filteredProviders = useMemo(() => {
    if (!searchValue) return allProviders;
    const lowerSearch = searchValue.toLowerCase();
    return allProviders.filter(
      (item) =>
        item.title.toLowerCase().includes(lowerSearch) ||
        item.description.toLowerCase().includes(lowerSearch)
    );
  }, [searchValue, allProviders]);

  const handleSearchChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchValue(e.target.value);
  }, []);

  return (
    <Grid2 container rowSpacing={3} columnSpacing={3}>
      <Grid2 size={12}>
        <SearchContainer>
          <SearchTextField
            fullWidth
            placeholder="Search Connections"
            value={searchValue}
            onChange={handleSearchChange}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <Search />
                </InputAdornment>
              ),
            }}
          />
        </SearchContainer>
      </Grid2>
      {filteredProviders.length > 0 ? (
        filteredProviders.map((item) => (
          <Grid2 size={{ xs: 12, sm: 6, md: 4, lg: 3 }} key={item.provider}>
            <ProviderCard
              provider={item.provider}
              title={item.title}
              description={item.description}
              onClick={onProviderClick}
            />
          </Grid2>
        ))
      ) : (
        <Grid2 size={12}>
          <EmptyStateBox>No connections found</EmptyStateBox>
        </Grid2>
      )}
    </Grid2>
  );
}
