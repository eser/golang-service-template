export enum FormStateEntryStatus {
  PENDING = "PENDING",
  SUCCESS = "SUCCESS",
  ERROR = "ERROR",
}

export type FormStateEntry = [Date, string, FormStateEntryStatus, number | null];
export type FormState = Array<FormStateEntry>;
