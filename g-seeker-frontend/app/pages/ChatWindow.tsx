"use client"
import { useState } from "react";
import ChatSidebar from "../components/ChatSidebar";
import ChatHeader from "../components/ChatHeader";
import MessageBubble from "../components/MessageBubble";
import RecommendationPanel from "../components/RecommendationPanel";
import ChatInput from "../components/ChatInput";
import type { Message, RepoCardData } from "../types/chat";

const initialMessages: Message[] = [
  {
    id: 1,
    role: "assistant",
    content:
      "你好，我是 OpenSource Navigator。告诉我你的开发需求，我会帮你推荐合适的 GitHub 开源项目。",
    time: "09:41",
  },
  {
    id: 2,
    role: "user",
    content: "我想找一个适合中小团队使用的 Go 微服务脚手架。",
    time: "09:42",
  },
  {
    id: 3,
    role: "assistant",
    content:
      "收到，我会优先关注 Go 技术栈、微服务支持、文档质量和社区活跃度。你也可以继续补充是否更偏向轻量方案或企业级方案。",
    time: "09:42",
  },
];

const repoCards: RepoCardData[] = [
  {
    id: 1,
    name: "go-zero",
    stars: "29.8k",
    desc: "内置服务治理、网关与中间件能力，适合快速搭建 Go 微服务系统。",
    reason: "生态成熟，脚手架能力完整，适合中小团队快速起步。",
  },
  {
    id: 2,
    name: "Kratos",
    stars: "24.1k",
    desc: "Bilibili 开源的 Go 微服务框架，分层清晰，工程化能力较强。",
    reason: "架构规范，适合追求可维护性和扩展性的团队。",
  },
];

export default function ChatWindow() {
  const [userMessage, setUserMessage] = useState("");
  const [messages, setMessages] = useState<Message[]>(initialMessages);

  const handleSendMessage = () => {
    const trimmedMessage = userMessage.trim();
    if (!trimmedMessage) return;

    const newMessage: Message = {
      id: Date.now(),
      role: "user",
      content: trimmedMessage,
      time: new Date().toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      }),
    };

    setMessages((prev) => [...prev, newMessage]);
    setUserMessage("");
  };

  return (
    <div className="min-h-screen bg-slate-100 p-6 md:p-10">
      <div className="mx-auto grid max-w-7xl gap-6 lg:grid-cols-[300px_minmax(0,1fr)]">
        <ChatSidebar />

        <main className="flex min-h-[80vh] flex-col rounded-3xl bg-white shadow-sm ring-1 ring-slate-200">
          <ChatHeader />

          <section className="flex-1 overflow-y-auto px-5 py-5 md:px-6">
            <div className="mx-auto max-w-4xl space-y-6">
              {messages.map((message) => (
                <MessageBubble key={message.id} message={message} />
              ))}

              <RecommendationPanel repoCards={repoCards} />
            </div>
          </section>

          <ChatInput
            value={userMessage}
            onChange={setUserMessage}
            onSend={handleSendMessage}
          />
        </main>
      </div>
    </div>
  );
}
