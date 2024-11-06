import { Link } from "@tanstack/react-router";
import { ExternalLink } from "@components/ExternalLink";
import Logo from "@app/logo.svg";
import { Text, Button, Center, Stack, Container, Image } from "@mantine/core";

export const NotFound = () => {
  return (
    <Container size="md">
      <Center>
        <Image src={Logo} alt="Logo" width={120} height={120} fit="contain" />
      </Center>

      <Stack align="center" mt="md">
        <Text size="xl">404 Page not found</Text>
        <Text size="lg" c="dimmed">
          Oops, looks like there was a little too much plot!
        </Text>
        <Text size="md" c="dimmed">
          In case you think this is a bug rather than too much plot,
        </Text>
        <Text size="md" c="dimmed">
          feel free to report this to our{" "}
          <ExternalLink href="https://github.com/autobrr/autobrr">
            GitHub page
          </ExternalLink>{" "}
          or to{" "}
          <ExternalLink href="https://discord.gg/WQ2eUycxyT">
            our official Discord channel
          </ExternalLink>
          .
        </Text>
        <Text size="md" c="dimmed">
          Otherwise, let us help you to get you back on track for more plot!
        </Text>
      </Stack>

      <Center mt="lg">
        <Link to="/">
          <Button size="md" color="blue">
            Back to Dashboard
          </Button>
        </Link>
      </Center>
    </Container>
  );
};
