import { Button } from "antd";
import type { KeyboardEvent } from "react";

type Props = {
  value: string;
  onChange: (value: string) => void;
  onSend: () => void;
};

export default function ChatInput({ value, onChange, onSend }: Props) {
  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      onSend();
    }
  };

  return (
    <footer className="border-t border-slate-200 px-5 py-4 md:px-6">
      <div className="mx-auto max-w-4xl">
        <div className="rounded-3xl border border-slate-200 bg-slate-50 p-3 shadow-inner">
          <label
            htmlFor="chat-input"
            className="mb-2 block text-xs font-semibold uppercase tracking-[0.16em] text-slate-400"
          >
            Describe your requirement
          </label>

          <div className="relative">
            <textarea
              id="chat-input"
              value={value}
              rows={4}
              onChange={(e) => onChange(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="例如：我想找一个适合中小团队、支持 Go 微服务、文档完善的开源脚手架"
              className="w-full resize-none rounded-2xl border border-slate-200 bg-white px-4 py-3 pb-16 text-sm text-slate-700 outline-none transition placeholder:text-slate-400 focus:border-slate-400"
            />

            <div className="absolute bottom-3 right-3 flex items-center gap-2">
              <Button className="rounded-2xl border border-slate-200 px-4 py-2 text-sm font-medium text-slate-700 transition hover:bg-white">
                Attach
              </Button>
              <Button
                className="rounded-2xl bg-slate-900 px-5 py-2 text-sm font-medium text-white transition hover:opacity-90"
                onClick={onSend}
              >
                Send
              </Button>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
}
