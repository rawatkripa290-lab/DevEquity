
import { Trophy, Users, TrendingUp, Award } from "lucide-react";

function ContributionSummary({ contributors = [] }) {
  const totalContributors = contributors.length;

  const totalCommits = contributors.reduce(
    (sum, contributor) => sum + (contributor.commits || 0),
    0
  );

  const averageScore = totalContributors
    ? (
        contributors.reduce(
          (sum, contributor) => sum + (contributor.score || 0),
          0
        ) / totalContributors
      ).toFixed(2)
    : 0;

  const topContributor = [...contributors].sort(
    (a, b) => b.score - a.score
  )[0];

  const summary = [
    {
      title: "Analyzed Contributors",
      value: totalContributors,
      icon: Users,
    },
    {
      title: "Total Contributor Commits",
      value: totalCommits,
      icon: TrendingUp,
    },
    {
      title: "Average Fair Score",
      value: `${averageScore}/100`,
      icon: Award,
    },
    {
      title: "Top Contributor",
      value: topContributor?.contributor || "N/A",
      icon: Trophy,
    },
  ];

  return (
    <div className="mt-8 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
      {summary.map((item) => {
        const Icon = item.icon;

        return (
          <div
            key={item.title}
            className="rounded-xl border border-blue-600 bg-blue-400 p-6"
          >
            <div className="flex items-center justify-between">
              <p className="text-medium font-semibold text-gray-900">{item.title}</p>
              <Icon size={22} className="text-blue-400" />
            </div>

            <h3 className="mt-4 break-words text-2xl font-bold text-blue-900">
              {item.value}
            </h3>
          </div>
        );
      })}
    </div>
  );
}

export default ContributionSummary;
