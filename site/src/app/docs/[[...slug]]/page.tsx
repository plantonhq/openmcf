import { Metadata } from 'next';
import { notFound } from 'next/navigation';
import {
  getMarkdownContent,
  getDocumentationStructure,
  getNextDocItem,
  generateStaticParamsFromStructure,
  processDocumentationSlug,
} from '@/app/docs/utils/fileSystem';
import { MDXRenderer } from '@/app/docs/components/MDXRenderer';
import { Author, MDXParser } from '@/lib/mdx';
import { DocsLayout } from '@/app/docs/components/DocsLayout';
import matter from 'gray-matter';

type DocsParams = Promise<{ slug?: string[] }>;

export async function generateMetadata({ params }: { params: DocsParams }): Promise<Metadata> {
  const { slug = [] } = await params;
  const { path } = processDocumentationSlug(slug);

  try {
    const content = await getMarkdownContent(path);
    const { data } = matter(content);
    const title = data?.title || slug[slug.length - 1] || 'Documentation';

    return {
      title: `${title} - OpenMCF Documentation`,
      description: data?.description || 'OpenMCF Documentation',
    };
  } catch {
    return {
      title: 'Documentation - OpenMCF',
      description: 'OpenMCF Documentation',
    };
  }
}

export async function generateStaticParams() {
  const structure = await getDocumentationStructure();
  return generateStaticParamsFromStructure(structure);
}

export default async function DocsPage({ params }: { params: DocsParams }) {
  const { slug = [] } = await params;
  const { path } = processDocumentationSlug(slug);

  try {
    const content = await getMarkdownContent(path);
    const { data } = matter(content);
    const mdxContent = MDXParser.reconstructMDX(content);

    // Get the documentation structure to find the next item
    const allDocs = await getDocumentationStructure();
    const nextDocItem = getNextDocItem(path, allDocs);

    // Normal rendering with full layout
    return (
      <DocsLayout author={data?.author as unknown as Author[]} content={content}>
        <MDXRenderer
          mdxContent={mdxContent}
          nextArticle={
            nextDocItem
              ? {
                  title: nextDocItem.title,
                  excerpt: nextDocItem.excerpt,
                  slug: `/docs/${nextDocItem.slug}`,
                }
              : undefined
          }
        />
      </DocsLayout>
    );
  } catch (error) {
    console.error('Error loading documentation:', error);
    notFound();
  }
}

