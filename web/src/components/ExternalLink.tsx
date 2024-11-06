import { Anchor } from "@mantine/core";

type ExternalLinkProps = {
  href: string;
  children?: React.ReactNode;
};

export const ExternalLink = ({ href, children }: ExternalLinkProps) => (
  <Anchor href={href} target="_blank" rel="noopener noreferrer">
    {children}
  </Anchor>
);

export const DocsLink = ({ href }: { href: string }) => (
  <ExternalLink href={href}>{href}</ExternalLink>
);
