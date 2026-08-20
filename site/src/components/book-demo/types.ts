export interface DemoFormData {
  firstName: string;
  lastName: string;
  workEmail: string;
  company: string;
  jobTitle: string;
  companySize: string;
}

export const JOB_TITLE_OPTIONS = [
  { value: 'CTO', label: 'CTO' },
  { value: 'VP Engineering', label: 'VP Engineering' },
  { value: 'Director of Engineering', label: 'Director of Engineering' },
  { value: 'Engineering Manager', label: 'Engineering Manager' },
  { value: 'Platform Engineer', label: 'Platform Engineer' },
  { value: 'DevOps Engineer', label: 'DevOps Engineer' },
  { value: 'SRE', label: 'SRE' },
  { value: 'Software Engineer', label: 'Software Engineer' },
  { value: 'Solutions Architect', label: 'Solutions Architect' },
  { value: 'Founder / CEO', label: 'Founder / CEO' },
  { value: 'Other', label: 'Other' },
] as const;

export const COMPANY_SIZE_OPTIONS = [
  { value: '1-10', label: '1-10' },
  { value: '11-50', label: '11-50' },
  { value: '51-200', label: '51-200' },
  { value: '200+', label: '200+' },
] as const;

export type Phase = 'form' | 'scheduler';

export type SubmissionStatus = 'idle' | 'submitting' | 'success' | 'error';

export interface FieldErrors {
  firstName?: string;
  lastName?: string;
  workEmail?: string;
  company?: string;
  jobTitle?: string;
  companySize?: string;
}

export const SUBMISSION_ENDPOINT = 'https://webhooks.planton.ai/demo';
