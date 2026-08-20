import { NextResponse } from 'next/server';
import { getDocumentationStructure } from '../../../(root)/docs/utils/fileSystem';

export const dynamic = 'force-static';

export async function GET() {
  try {
    const structure = await getDocumentationStructure();
    return NextResponse.json(structure);
  } catch (error) {
    console.error('Error loading documentation structure:', error);
    return NextResponse.json({ error: 'Failed to load documentation structure' }, { status: 500 });
  }
}
