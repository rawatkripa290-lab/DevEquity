import { useState } from "react";
import { Search } from "lucide-react";
import ContributionChart from "./ContributionChart";
import TopContributorsChart from "./TopContributorsChart";
function ContributionTable({ contributors = [] }) {
  const [search, setSearch] = useState("");
  const [sortBy, setSortBy] = useState("highest-score");
  const [showContributionChart, setShowContributionChart] = useState(false);
  const [showTopContributorsChart, setShowTopContributorsChart] =
    useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedContributor, setSelectedContributor] = useState(null);

  const itemsPerPage = 10;

  const getScoreStyle = (score) => {
  if (score >= 75)
    return "bg-emerald-100 text-emerald-800 border border-emerald-200";

  if (score >= 40)
    return "bg-amber-100 text-amber-800 border border-amber-200";

  return "bg-rose-100 text-rose-800 border border-rose-200";
};
  const filteredContributors = contributors.filter((contributor) =>
    contributor.contributor?.toLowerCase().includes(search.toLowerCase()),
  );

  const sortedContributors = [...filteredContributors].sort((a, b) => {
    switch (sortBy) {
      case "lowest-score":
        return a.score - b.score;

      case "most-commits":
        return b.commits - a.commits;

      case "fewest-commits":
        return a.commits - b.commits;

      case "most-prs":
        return (b.total_prs ?? 0) - (a.total_prs ?? 0);

      case "most-changes":
        return (b.total_changes ?? 0) - (a.total_changes ?? 0);

      default:
        return b.score - a.score;
    }
  });

  // Pagination
  const totalPages = Math.ceil(sortedContributors.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;

  const currentContributors = sortedContributors.slice(
    startIndex,
    startIndex + itemsPerPage,
  );
  const totalContributors = contributors.length;

  const totalCommits = contributors.reduce(
    (sum, c) => sum + (c.commits ?? 0),
    0,
  );

  const totalPRs = contributors.reduce((sum, c) => sum + (c.total_prs ?? 0), 0);

  const totalCodeChanges = contributors.reduce(
    (sum, c) => sum + (c.total_changes ?? 0),
    0,
  );

  const averageScore = totalContributors
    ? (
        contributors.reduce((sum, c) => sum + (c.score ?? 0), 0) /
        totalContributors
      ).toFixed(2)
    : 0;
  const exportCSV = () => {
    const headers = [
      "Rank",
      "Contributor",
      "Commits",
      "Fair Score",
      "Additions",
      "Deletions",
      "Total PRs",
      "Merged PRs",
      "Reviews",
      "Issues",
      "Active Days",
    ];

    const rows = sortedContributors.map((contributor, index) => [
      index + 1,
      contributor.contributor,
      contributor.commits ?? 0,
      contributor.score ?? 0,
      contributor.additions ?? 0,
      contributor.deletions ?? 0,
      contributor.total_prs ?? 0,
      contributor.merged_prs ?? 0,
      contributor.total_reviews ?? 0,
      contributor.total_issues ?? 0,
      contributor.active_days ?? 0,
    ]);

    const csvContent = [headers, ...rows]
      .map((row) =>
        row.map((value) => `"${String(value).replace(/"/g, '""')}"`).join(","),
      )
      .join("\n");

    const blob = new Blob([csvContent], {
      type: "text/csv;charset=utf-8;",
    });

    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");

    link.href = url;
    link.download = "contributor-report.csv";
    link.click();

    URL.revokeObjectURL(url);
  };

  const getContributionLevel = (score) => {
    if (score >= 75) return "Excellent";
    if (score >= 40) return "Good";
    return "Needs Improvement";
  };

  return (
    <div className="w-full overflow-x-auto rounded-xl">
      {/* Contribution Analytics Dashboard */}
      <div className="grid grid-cols-1 gap-4 p-6 sm:grid-cols-2 lg:grid-cols-5">
        {[
          ["Total Contributors", totalContributors],
          ["Total Commits", totalCommits],
          ["Total Pull Requests", totalPRs],
          ["Total Code Changes", totalCodeChanges],
          ["Average Fair Score", averageScore],
        ].map(([label, value]) => (
          <div
            key={label}
            className="rounded-xl border border-yellow-300 bg-yellow-200 p-5 shadow-sm transition-all duration-300 hover:-translate-y-1 hover:border-blue-300 hover:shadow-lg hover:shadow-blue-100"
          >
            <p className="text-medium text-blue-900">{label}</p>
            <p className="mt-3 text-2xl font-bold text-blue-900">{value}</p>
          </div>
        ))}
      </div>
      {/* Graph Controls */}
      <div className="flex flex-wrap gap-3 my-6 justify-between p-6">
        <button
          onClick={() => setShowContributionChart(!showContributionChart)}
          className="px-5 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition"
        >
          {showContributionChart
            ? "Hide Contribution Graph"
            : "Show Contribution Graph"}
        </button>

        <button
          onClick={() => setShowTopContributorsChart(!showTopContributorsChart)}
          className="px-5 py-2 rounded-lg bg-yellow-200 text-blue-900 font-medium hover:bg-yellow-200 transition"
        >
          {showTopContributorsChart
            ? "Hide Top Contributors Graph"
            : "Show Top Contributors Graph"}
        </button>
      </div>
      {showContributionChart && (
        <ContributionChart contributors={contributors} />
      )}

      {showTopContributorsChart && (
        <TopContributorsChart contributors={contributors} />
      )}
      {selectedContributor && (
        <div className="m-6 rounded-xl border border-blue-500/30 bg-white p-6">
          <div className="flex items-center justify-between">
            <h3 className="text-xl font-bold text-blue-900">
              Contributor Details
            </h3>

            <button
              onClick={() => setSelectedContributor(null)}
              className="rounded-lg bg-slate-800 px-3 py-2 text-sm text-white hover:bg-slate-700"
            >
              Close
            </button>
          </div>

          <div className="mt-5 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {[
              ["Username", selectedContributor.contributor],
              [
                "Fair Contribution Score",
                `${selectedContributor.score ?? 0}/100`,
              ],
              [
                "Contribution Level",
                getContributionLevel(selectedContributor.score ?? 0),
              ],
              ["Total Commits", selectedContributor.commits],
              ["Total Code Changes", selectedContributor.total_changes],
              ["Total Pull Requests", selectedContributor.total_prs],
              ["Merged Pull Requests", selectedContributor.merged_prs],
              ["Total Reviews", selectedContributor.total_reviews],
              ["Total Issues", selectedContributor.total_issues],
              ["Active Days", selectedContributor.active_days],
            ].map(([label, value]) => (
              <div key={label} className="rounded-lg bg-yellow-200 p-4">
                <p className="text-sm text-blue-900">{label}</p>

                <p className="mt-2 break-words text-xl font-bold text-blue-900">
                  {value ?? 0}
                </p>
              </div>
            ))}
          </div>
        </div>
      )}
      <div className="flex flex-col gap-4 border-b border-slate-800 p-6 md:flex-row md:items-center md:justify-between">
        <div>
          <h2 className="text-2xl font-bold text-blue-900">
            Contributor Rankings
          </h2>

          <p className="mt-1 text-medium text-blue-900">
            Overview of developer contributions
          </p>
        </div>

        <div className="flex flex-col gap-3 sm:flex-row">
          <div className="flex items-center gap-3 rounded-lg border border-slate-700 bg-white px-4 py-3">
            <Search size={18} className="text-blue-900" />

            <input
              type="text"
              placeholder="Search contributor..."
              value={search}
              onChange={(e) => {
                setSearch(e.target.value);
                setCurrentPage(1);
              }}
              className="w-full bg-transparent text-sm text-black outline-none placeholder:text-slate-500"
            />
          </div>

          <button
            onClick={exportCSV}
            className="rounded-lg bg-blue-900 px-5 py-3 text-sm font-semibold text-white transition hover:bg-blue-500"
          >
            Export CSV
          </button>
        </div>
      </div>

      {/* Sorting */}
      <div className="p-5">
        <select
          value={sortBy}
          onChange={(e) => {
            setSortBy(e.target.value);
            setCurrentPage(1);
          }}
          className="rounded-lg border border-slate-700 bg-blue-900 px-4 py-3 text-medium text-white outline-none"
        >
          <option value="highest-score">Highest Score</option>
          <option value="lowest-score">Lowest Score</option>
          <option value="most-commits">Most Commits</option>
          <option value="fewest-commits">Fewest Commits</option>
          <option value="most-prs">Most Pull Requests</option>
          <option value="most-changes">Most Code Changes</option>
        </select>
      </div>

      {/* Contributor Table */}
      
{/* Contributor Table */}
<div className="overflow-x-auto rounded-xl border border-slate-200 bg-white shadow-lg">
  <table className="w-full min-w-[750px] text-left">
    <thead className="bg-slate-900 text-white">
      <tr>
        <th className="px-6 py-5 text-sm font-semibold uppercase tracking-wider">
          Rank
        </th>
        <th className="px-6 py-5 text-sm font-semibold uppercase tracking-wider">
          Contributor
        </th>
        <th className="px-6 py-5 text-sm font-semibold uppercase tracking-wider">
          Commits
        </th>
        <th className="px-6 py-5 text-sm font-semibold uppercase tracking-wider">
          Code Changes
        </th>
        <th className="px-6 py-5 text-sm font-semibold uppercase tracking-wider">
          Pull Requests
        </th>
        <th className="px-6 py-5 text-sm font-semibold uppercase tracking-wider">
          Fair Score
        </th>
      </tr>
    </thead>

    <tbody className="divide-y divide-slate-200">
      {currentContributors.length > 0 ? (
        currentContributors.map((contributor) => {
          const rank = sortedContributors.indexOf(contributor);

          return (
            <tr
              key={contributor.contributor}
              className="transition-colors duration-200 hover:bg-blue-50"
            >
              <td className="px-6 py-5">
                <span className="inline-flex h-8 w-8 items-center justify-center rounded-full bg-slate-100 text-sm font-bold text-slate-700">
                  {rank + 1}
                </span>
              </td>

              <td className="px-6 py-5">
                <button
                  onClick={() => setSelectedContributor(contributor)}
                  className="font-semibold text-blue-700 hover:text-blue-900 hover:underline"
                >
                  {contributor.contributor}
                </button>
              </td>

              <td className="px-6 py-5 text-sm font-medium text-slate-700">
                {contributor.commits ?? 0}
              </td>

              <td className="px-6 py-5 text-sm font-medium text-slate-700">
                {(contributor.total_changes ?? 0).toLocaleString()}
              </td>

              <td className="px-6 py-5 text-sm font-medium text-slate-700">
                {contributor.total_prs ?? 0}
              </td>

              <td className="px-6 py-5">
                <span
                  className={`inline-flex min-w-[85px] justify-center rounded-full px-3 py-1.5 text-xs font-bold ${getScoreStyle(
                    contributor.score ?? 0
                  )}`}
                >
                  {(contributor.score ?? 0).toFixed(1)}/100
                </span>
              </td>
            </tr>
          );
        })
      ) : (
        <tr>
          <td
            colSpan="6"
            className="px-6 py-12 text-center text-sm font-medium text-slate-500"
          >
            No contributors found.
          </td>
        </tr>
      )}
    </tbody>
  </table>
</div>


      {/* Pagination */}
      <div className="flex items-center justify-between border-t border-slate-800 p-5">
        <p className="text-sm text-blue-900">
          Page {currentPage} of {Math.max(totalPages, 1)}
        </p>

        <div className="flex gap-3">
          <button
            onClick={() => setCurrentPage((page) => Math.max(page - 1, 1))}
            disabled={currentPage === 1}
            className="rounded-lg bg-slate-800 px-4 py-2 text-sm text-white disabled:opacity-40"
          >
            Previous
          </button>

          <button
            onClick={() =>
              setCurrentPage((page) => Math.min(page + 1, totalPages))
            }
            disabled={currentPage >= totalPages || totalPages === 0}
            className="rounded-lg bg-blue-600 px-4 py-2 text-sm text-white disabled:opacity-40"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  );
}

export default ContributionTable;
