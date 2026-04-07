import { Message } from "../types/chat";

type Props = {
  message: Message;
};

export default function MessageBubble({ message }: Props) {
  const isUser = message.role === "user";

  return (
    <div className={`flex ${isUser ? "justify-end" : "justify-start"}`}>
      <div
        className={`max-w-[85%] rounded-3xl px-4 py-3 shadow-sm md:max-w-[70%] ${
          isUser ? "bg-slate-900 text-white" : "bg-slate-100 text-slate-800"
        }`}
      >
        <div className="mb-2 flex items-center gap-2 text-xs opacity-70">
          <span className="font-medium">{isUser ? "You" : "Navigator"}</span>
          <span>{message.time}</span>
        </div>
        <p className="text-sm leading-7">{message.content}</p>
      </div>
    </div>
  );
}
