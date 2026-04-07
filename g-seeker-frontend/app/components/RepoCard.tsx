import { Button } from "antd";

import { RepoCardData } from "../types/chat";

type Props = {
  repo: RepoCardData;
};

export default function RepoCard({ repo }: Props) {
  return (
    <div className="rounded-3xl bg-white p-4 shadow-sm ring-1 ring-slate-200">
      <div className="flex items-start justify-between gap-3">
        <div>
          <h3 className="text-base font-semibold text-slate-900">{repo.name}</h3>
          <p className="mt-1 text-xs text-slate-400">
            GitHub Stars · {repo.stars}
          </p>
        </div>
        <Button className="rounded-full border border-slate-200 px-3 py-1 text-xs text-slate-600 transition hover:bg-slate-50">
          View Repo
        </Button>
      </div>

      <p className="mt-4 text-sm leading-7 text-slate-600">{repo.desc}</p>

      <div className="mt-4 rounded-2xl bg-slate-50 p-3">
        <p className="text-xs font-semibold uppercase tracking-[0.16em] text-slate-400">
          Why recommended
        </p>
        <p className="mt-2 text-sm leading-7 text-slate-700">{repo.reason}</p>
      </div>
    </div>
  );
}
