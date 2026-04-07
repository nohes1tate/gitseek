export type MessageRole = "assistant" | "user";

export interface Message {
  id: number;
  role: MessageRole;
  content: string;
  time: string;
}

export interface RepoCardData {
  id: number;
  name: string;
  stars: string;
  desc: string;
  reason: string;
}
