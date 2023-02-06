import Head from 'next/head';
import { slugifyWithCounter } from '@sindresorhus/slugify';

import { Layout } from '@/components/Layout';

import 'focus-visible';
import '@/styles/tailwind.css';

function getNodeText(node) {
  let text = '';
  for (let child of node.children ?? []) {
    if (typeof child === 'string') {
      text += child;
    }
    text += getNodeText(child);
  }
  return text;
}

function collectHeadings(nodes, slugify = slugifyWithCounter()) {
  let sections = [];

  for (let node of nodes) {
    if (node.name === 'h2' || node.name === 'h3') {
      let title = getNodeText(node);
      if (title) {
        let id = slugify(title);
        node.attributes.id = id;
        if (node.name === 'h3') {
          if (!sections[sections.length - 1]) {
            throw new Error(
              'Cannot add `h3` to table of contents without a preceding `h2`'
            );
          }
          sections[sections.length - 1].children.push({
            ...node.attributes,
            title,
          });
        } else {
          sections.push({ ...node.attributes, title, children: [] });
        }
      }
    }

    sections.push(...collectHeadings(node.children ?? [], slugify));
  }

  return sections;
}

export default function App({ Component, pageProps }) {
  let title = pageProps.markdoc?.frontmatter.title;

  let pageTitle =
    pageProps.markdoc?.frontmatter.pageTitle ||
    `${pageProps.markdoc?.frontmatter.title} - Docs`;

  let description = pageProps.markdoc?.frontmatter.description;

  let tableOfContents = pageProps.markdoc?.content
    ? collectHeadings(pageProps.markdoc.content)
    : [];

  return (
    <>
      <Head>
        <title>{pageTitle}</title>
        {description && <meta name="description" content={description} />}
        <link rel="shortcut icon" href="/favicon.png" />
        <meta name="title" content="Unweave: ML, without the FML" />
        <meta
            name="description"
            content="Open source machine learning dev environments. Click the link, try it out 👾"
        />
        <meta property="og:type" content="website" />
        <meta property="og:url" content="https://unweave.io/" />
        <meta property="og:title" content="Unweave: ML, without the FML" />
        <meta
            property="og:description"
            content="SSH into GPU machines across cloud providers. Why are you still reading this? Click the link, try it out 👾"
        />
        <meta property="og:site_name" content="Unweave" />
        <meta property="og:image" content="https://unweave.io/meta-image.png" />
        <meta property="twitter:card" content="summary_large_image" />
        <meta name="twitter:creator" content="@unweaveio" />
        <meta property="twitter:url" content="https://unweave.io/" />
        <meta property="twitter:title" content="Unweave: ML, without the FML" />
        <meta
            property="twitter:image"
            content="https://unweave.io/meta-image.png"
        />
      </Head>
      <Layout title={title} tableOfContents={tableOfContents}>
        <Component {...pageProps} />
      </Layout>
    </>
  );
}
