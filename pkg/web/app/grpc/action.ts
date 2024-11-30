"use server";

import * as grpc from "@grpc/grpc-js";
import * as broadcast from "../../proto-ts/broadcast/broadcast.ts";
import { type FormState, type FormStateEntry, FormStateEntryStatus } from "./types.ts";

const sendMessage = (message: string, retryCount = 0) => {
  return new Promise((resolve, reject) => {
    const client = new broadcast.MessageServiceClient(
      "localhost:9090",
      grpc.credentials.createInsecure(),
      {
        "grpc.enable_retries": 1,
        "grpc.keepalive_time_ms": 10000,
        "grpc.keepalive_timeout_ms": 5000,
        "grpc.keepalive_permit_without_calls": 1,
      },
    );

    // Wait for client to be ready before sending
    client.waitForReady(Date.now() + 5000, (error) => {
      if (error) {
        reject(error);
        return;
      }

      client.send(
        {
          channelId: "1",
          message: {
            body: message,
          },
        },
        (error, response) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
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
