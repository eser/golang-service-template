import React from "react";
import { Box, Flex, Heading, Link, Text } from "@radix-ui/themes";

export default function Home() {
  return (
    <Flex flexGrow="1" align="center" justify="center" asChild>
      <main>
        <Box className="~max-w-xs/4xl p-4">
          <Flex direction="column" gap="4">
            <Heading as="h1" className="~text-lg/4xl">
              JavaScript 🤝🏻 Golang
            </Heading>
            <Heading as="h2" mb="6" className="~text-base/2xl">
              a tech preview demo crafted by eser.live 👋🏻
            </Heading>

            <Text as="p" className="~text-sm/xl">
              Hello, this page is designed to demonstrate the interoperability between JavaScript and Golang.
            </Text>
            <Text as="p" className="~text-sm/xl">
              Initially, it will be used for a presentation (you may take a look at{" "}
              <Link href="https://speakerdeck.com/eser">speakerdeck.com/eser</Link>{" "}
              for presentation slides), but soon there'll be some content focused on this topic on{" "}
              <Link href="https://eser.live/" target="_blank">eser.live</Link>.
            </Text>
          </Flex>
        </Box>
      </main>
    </Flex>
  );
}
