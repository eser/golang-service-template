"use server";

import { type FormState, type FormStateEntry, FormStateEntryStatus } from "./types.ts";

const sendConcatMessage = async (messages: string[]) => {
  await fetch("http://localhost:8080/send", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      channelId: "2",
      message: {
        body: messages.join(" "), // Concatenate the strings
      },
    }),
  });
};

export const sendConcatMessageAction = async (
  prevState: FormState,
  formData: FormData,
): Promise<FormState> => {
  const message1 = formData.get("message1") as string;
  const message2 = formData.get("message2") as string;
  const messages = [message1, message2];

  const start = performance.now();
  let newStateEntryStatus: FormStateEntryStatus;

  try {
    const _response = await sendConcatMessage(messages);
    newStateEntryStatus = FormStateEntryStatus.SUCCESS;
  } catch (error) {
    newStateEntryStatus = FormStateEntryStatus.ERROR;
  }

  const took = performance.now() - start;
  const newMessage: FormStateEntry = [new Date(), messages, newStateEntryStatus, took];

  return [...prevState, newMessage];
};
