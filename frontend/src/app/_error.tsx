import Error from 'next/error';
import type { NextPageContext } from 'next';

interface ErrorPageProps {
  statusCode: number;
}

function Page({ statusCode }: ErrorPageProps) {
  return <Error statusCode={statusCode} />;
}

Page.getInitialProps = ({ res, err }: NextPageContext): ErrorPageProps => {
  const statusCode = res?.statusCode || err?.statusCode || 404;
  return { statusCode };
};

export default Page;