"use server";

import { type FormState, type FormStateEntry, FormStateEntryStatus } from "./types.ts";

const sendMessage = async (message: string) => {
  await fetch("http://localhost:8080/send", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      channelId: "2",
      message: {
        body: message,
      },
    }),
  });
};

export const sendMessageAction = async (
  prevState: FormState,
  formData: FormData,
): Promise<FormState> => {
  const message = formData.get("message") as string;

  const start = performance.now();

  let newStateEntryStatus: FormStateEntryStatus;

  try {
    const _response = await sendMessage(message);

    newStateEntryStatus = FormStateEntryStatus.SUCCESS;
  } catch (error) {
    newStateEntryStatus = FormStateEntryStatus.ERROR;
  }

  const took = performance.now() - start;

  const newMessage: FormStateEntry = [new Date(), message, newStateEntryStatus, took];

  return [...prevState, newMessage];
};
