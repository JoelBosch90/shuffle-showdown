import type { Answer } from '$lib/types/Answer';

export interface AnswerSelectionUpdate {
  sentAt: Date;
  playerId: string;
  answer: Answer;
}