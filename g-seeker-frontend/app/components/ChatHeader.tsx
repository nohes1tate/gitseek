import { Button } from "antd";

export default function ChatHeader() {
  return (
    <header className="flex items-center justify-between border-b border-slate-200 px-5 py-4 md:px-6">
      <div>
        <h2 className="text-lg font-semibold text-slate-900">需求对话窗口</h2>
        <p className="text-sm text-slate-500">
          自然语言输入需求，返回候选仓库与推荐原因
        </p>
      </div>

      <div className="hidden items-center gap-2 md:flex">
        <Button className="rounded-2xl border border-slate-200 px-4 py-2 text-sm text-slate-700 transition hover:bg-slate-50">
          Clear
        </Button>
        <Button className="rounded-2xl border border-slate-900 bg-slate-900 px-4 py-2 text-sm text-white transition hover:opacity-90">
          Export
        </Button>
      </div>
    </header>
  );
}
