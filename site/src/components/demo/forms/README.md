# Demo Forms

This directory contains form components for the demo application, based on the original Planton forms but adapted for the demo environment.

## Components

### AwsCredentialForm

A complete AWS credential creation form that matches the original Planton form structure.

**Features:**

- Simple React state management (no external form libraries)
- No validation - simple form submission
- **Dynamic default values** from `aws-credential.json`
- **Full-screen modal layout** with proper scrolling
- **Fixed header and footer** with scrollable content
- **Summary sidebar** with progress indicators
- **Accordion-style sections** (Metadata, AWS Account Details)
- **Grid layout** for form fields
- **Help text icons** for each field
- **Refresh data button** to reload from JSON
- Type-safe with TypeScript interfaces

**Usage:**

```typescript
import { AwsCredentialForm } from '../forms';

<AwsCredentialForm
  onSubmit={(data) => console.log('Form submitted:', data)}
  onCancel={() => console.log('Form cancelled')}
  initialData={optionalInitialData}
/>;
```

### AwsCredentialFormElements

Reusable form elements for AWS credential fields.

**Fields:**

- Account ID (text input)
- Default Region (text input - not dropdown)
- Access Key ID (text input)
- Secret Access Key (password input)

**Layout:**

- **2-column grid** layout for better space utilization
- **Help text icons** for each field
- **Consistent styling** with rounded inputs and focus states

**Usage:**

```typescript
import { AwsCredentialFormElements } from '../forms';

<AwsCredentialFormElements
  formData={formData}
  onInputChange={handleInputChange}
  helpTextMapping={helpTextMapping}
  section="AWS Account Details"
/>;
```

## Integration

The forms are integrated into the CloudConnections modal:

- Click on AWS card → Opens full-screen modal with AWS credential form
- **Default values loaded** → Form is populated with data from `aws-credential.json`
- **Full-screen layout** → Modal takes up most of the screen with proper scrolling
- **Fixed header/footer** → Only content area scrolls, header and footer stay fixed
- **Summary sidebar** → Shows progress and form sections
- Form submission → Calls the modal's submit handler
- Form cancellation → Closes the modal

## Form Layout Features

### Modal Structure

- **Full-screen modal** (max-w-7xl, max-h-[95vh])
- **Fixed header** with title, refresh button, and action buttons
- **Scrollable content** area with form sections
- **Fixed summary sidebar** with progress indicators

### Form Sections

- **Metadata Section**: Name and Slug fields in accordion
- **AWS Account Details**: Account ID, Region, Access Key ID, Secret Access Key
- **Accordion headers** with expand/collapse icons
- **Grid layout** for better field organization

### Dynamic Data Features

- **Default values** from JSON file on load
- **Manual refresh** button to reload data
- **Live updates** every 5 seconds automatically
- **Cache busting** to ensure fresh data
- **Error handling** with fallback data

## Form Structure

The forms follow the same structure as the original Planton forms:

### Data Flow

1. **JSON Data Fetching**: `useAwsCredentialData` hook fetches from `/demo/aws-credential.json`
2. **Default Value Setting**: Form is populated with JSON data when available
3. **User Interaction**: User can modify any field
4. **Form Submission**: Data is passed to parent component

### Key Differences from Original

- **No validation**: Simplified for demo purposes
- **No react-hook-form**: Uses simple React state
- **Dynamic data**: Loads from JSON file instead of API
- **Full-screen modal**: Better UX for demo
- **Text input for region**: Matches the image specification

## File Structure

```
forms/
├── AwsCredentialForm.tsx          # Main form component
├── AwsCredentialFormElements.tsx  # Form field components
├── index.ts                       # Export index
└── README.md                      # This documentation
```

## Dependencies

- **React**: For component state and lifecycle
- **TypeScript**: For type safety
- **Tailwind CSS**: For styling
- **Custom hooks**: `useAwsCredentialData` for data fetching
