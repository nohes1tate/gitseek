import { Button } from "antd";

const recentConversations = [
  "Go 微服务项目推荐",
  "适合 AI Agent 的 RAG 框架",
  "前端低代码开源方案",
];

const searchPreferences = ["Go", "Microservice", "High Stars", "Good Docs"];

export default function ChatSidebar() {
  return (
    <aside className="rounded-3xl bg-white p-5 shadow-sm ring-1 ring-slate-200">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-slate-500">Phase 1 MVP</p>
          <h1 className="mt-1 text-xl font-semibold text-slate-900">
            OpenSource Navigator
          </h1>
        </div>
        <span className="rounded-full bg-emerald-50 px-3 py-1 text-xs font-medium text-emerald-700">
          Online
        </span>
      </div>

      <Button className="mt-5 w-full rounded-2xl bg-slate-900 px-4 py-3 text-sm font-medium text-white transition hover:opacity-90">
        + New Chat
      </Button>

      <div className="mt-6">
        <p className="mb-3 text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
          Recent Conversations
        </p>
        <div className="space-y-2">
          {recentConversations.map((item) => (
            <Button
              key={item}
              className="w-full rounded-2xl border border-slate-200 px-4 py-3 text-left text-sm text-slate-700 transition hover:border-slate-300 hover:bg-slate-50"
            >
              <p className="truncate font-medium">{item}</p>
            </Button>
          ))}
        </div>
      </div>

      <div className="mt-6 rounded-2xl bg-slate-50 p-4">
        <p className="text-sm font-semibold text-slate-800">搜索偏好</p>
        <div className="mt-3 flex flex-wrap gap-2 text-xs">
          {searchPreferences.map((tag) => (
            <span
              key={tag}
              className="rounded-full bg-white px-3 py-1 text-slate-600 ring-1 ring-slate-200"
            >
              {tag}
            </span>
          ))}
        </div>
      </div>
    </aside>
  );
}
