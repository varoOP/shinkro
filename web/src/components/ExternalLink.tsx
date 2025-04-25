import {Anchor} from "@mantine/core";
import React from "react";

type ExternalLinkProps = {
    href: string;
    children?: React.ReactNode;
};

export const ExternalLink = ({href, children}: ExternalLinkProps) => (
    <Anchor href={href} target="_blank" rel="noopener noreferrer">
        {children}
    </Anchor>
);
