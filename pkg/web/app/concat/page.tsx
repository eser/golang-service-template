"use client";

import React, { useActionState, useOptimistic } from "react";
import { Badge, Box, Button, Flex, Heading, Kbd, Spinner, Text, TextField } from "@radix-ui/themes";
import { sendConcatMessageAction } from "./action.ts";
import { type FormState, type FormStateEntry, FormStateEntryStatus } from "./types.ts";

function formatDate(date: Date): string {
  const formatted = date.toISOString();
  return formatted.substring(11, formatted.length - 1);
}

export default function Page() {
  const [state, action, isPending] = useActionState(
    async (currentState: FormState, payload: FormData) => {
      const message1 = (await payload.get("message1")) as string;
      const message2 = (await payload.get("message2")) as string;

      addMessagesOptimistic([message1, message2]);
      const newState = await sendConcatMessageAction(currentState, payload);

      return newState;
    },
    [],
  );

  const [optimisticState, addMessagesOptimistic] = useOptimistic<
    FormState,
    string[]
  >(state, (prevState, newMessages) => [
    ...prevState,
    [new Date(), newMessages, FormStateEntryStatus.PENDING, null],
  ]);

  return (
    <Flex flexGrow="1" align="center" justify="center" direction="column">
      <main className="flex-grow">
        <Box className="~max-w-xs/4xl p-4">
          <Flex direction="column" gap="4">
            <Heading as="h1" className="~text-lg/4xl">
              String Concatenation Demo
            </Heading>
            <Heading as="h2" mb="6" className="~text-base/2xl">
              Combine Two Messages
            </Heading>

            <Flex direction="column" gap="2" asChild>
              <form action={action}>
                <div>
                  {optimisticState.map(
                    (entry: FormStateEntry, index: number) => (
                      <Text as="div" key={index}>
                        <Kbd>{formatDate(entry[0])}</Kbd>
                        <Text ml="2">{entry[1].join(" ")}</Text>
                        {entry[2] === FormStateEntryStatus.PENDING ? (
                          <Badge ml="2" variant="soft">
                            sending...
                          </Badge>
                        ) : entry[2] === FormStateEntryStatus.ERROR ? (
                          <Badge ml="2" variant="soft" color="red">
                            error
                          </Badge>
                        ) : (
                          <Text ml="2">{entry[3]!.toFixed(2)}ms</Text>
                        )}
                      </Text>
                    ),
                  )}
                </div>

                <Flex direction="column" gap="2">
                  <TextField.Root
                    type="text"
                    name="message1"
                    placeholder="First Message"
                    required={true}
                    size="3"
                  />
                  <TextField.Root
                    type="text"
                    name="message2"
                    placeholder="Second Message"
                    required={true}
                    size="3"
                  />
                </Flex>

                <Button
                  type="submit"
                  disabled={isPending}
                  variant="soft"
                  size="3"
                >
                  {isPending && <Spinner />}
                  Concatenate & Send
                </Button>
              </form>
            </Flex>
          </Flex>
        </Box>
      </main>
      <footer className="w-full p-4 bg-gray-100">
        <Text size="2" align="center">
          Enter two messages above to concatenate them together
        </Text>
      </footer>
    </Flex>
  );
}
