
import { GitCommit, Users, GitPullRequest, CircleDot } from "lucide-react";

function StatsCard({ totalCommits, totalContributors, totalPRs, totalIssues }) {
  const stats = [
    {
      title: "Total Commits",
      value: totalCommits,
      icon: GitCommit,
    },
    {
      title: "Contributors",
      value: totalContributors,
      icon: Users,
    },
    {
      title: "Pull Requests",
      value: totalPRs,
      icon: GitPullRequest,
    },
    {
      title: "Total Issues",
      value: totalIssues,
      icon: CircleDot,
    },
  ];

  return (
    <div className="grid w-full grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {stats.map((stat, index) => {
        const Icon = stat.icon;

        return (
          <div
            key={index}
            className="rounded-xl border border-slate-200 bg-blue-400 p-5 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md"
          >
            <div className="flex items-center justify-between">
              <p className="text-medium text-blue-900">{stat.title}</p>
            
            </div>

            <h2 className="mt-3 text-2xl font-bold text-blue-900">
              {stat.value ?? 0}
            </h2>
          </div>
        );
      })}
    </div>
  );
}

export default StatsCard;
