
import { Star, GitFork, CircleDot, Code2, GitBranch } from "lucide-react";

function RepositoryOverview({ repository }) {
  const details = [
    {
      title: "Stars",
      value: repository?.stars,
      icon: Star,
    },
    {
      title: "Forks",
      value: repository?.forks,
      icon: GitFork,
    },
    {
      title: "Open Issues",
      value: repository?.open_issues,
      icon: CircleDot,
    },
    {
      title: "Language",
      value: repository?.language,
      icon: Code2,
    },
  ];

  return (
    <div className="mb-8 rounded-xl border border-white bg-white p-6">
      <div className="mb-6">
        <div className="flex items-center gap-3">
          <GitBranch className="text-blue-900" size={24} />
          <h2 className="text-xl font-bold text-blue-900">
            {repository?.owner}/{repository?.repository}
          </h2>
        </div>
        <p className="mt-2 text-medium  text-blue-800">
          Repository Overview
        </p>
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {details.map((detail) => {
          const Icon = detail.icon;

          return (
            <div
              key={detail.title}
              className="rounded-lg border border-slate-400 p-5"
            >
              <div className="flex items-center justify-between">
                <p className="text-medium text-slate-800">{detail.title}</p>
                <Icon size={20} className="text-blue-900" />
              </div>

              <h3 className="mt-3 text-2xl font-bold text-blue-900">
                {detail.value ?? "N/A"}
              </h3>
            </div>
          );
        })}
      </div>
    </div>
  );
}

export default RepositoryOverview;
