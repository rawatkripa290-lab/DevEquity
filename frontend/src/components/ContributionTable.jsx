import { useState } from "react";
import { Search } from "lucide-react";
import ContributionChart from "./ContributionChart";
import TopContributorsChart from "./TopContributorsChart";
function ContributionTable({ contributors = [] }) {
  const [search, setSearch] = useState("");
  const [sortBy, setSortBy] = useState("highest-score");
  const [showContributionChart, setShowContributionChart] = useState(false);
const [showTopContributorsChart, setShowTopContributorsChart] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedContributor, setSelectedContributor] = useState(null);

  const itemsPerPage = 10;

  const getScoreStyle = (score) => {
    if (score >= 75) return "bg-green-500/15 text-green-400";
    if (score >= 40) return "bg-yellow-500/15 text-yellow-400";
    return "bg-red-500/15 text-red-400";
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
    <div className="w-full overflow-x-auto rounded-xl border border-slate-800 bg-slate-900">
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
            className="rounded-xl border border-slate-700 bg-slate-950 p-5 transition-all duration-300 hover:-translate-y-1 hover:border-blue-500/50 hover:shadow-lg hover:shadow-blue-500/5"
          >
            <p className="text-sm text-slate-400">{label}</p>
            <p className="mt-3 text-2xl font-bold text-blue-400">{value}</p>
          </div>
        ))}
      </div>
      {/* Graph Controls */}
<div className="flex flex-wrap gap-3 my-6 justify-between p-6">

  <button
    onClick={() => setShowContributionChart(!showContributionChart)}
    className="px-5 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition"
  >
    {showContributionChart ? "Hide Contribution Graph" : "Show Contribution Graph"}
  </button>

  <button
    onClick={() => setShowTopContributorsChart(!showTopContributorsChart)}
    className="px-5 py-2 rounded-lg bg-yellow-200 text-blue-900 font-medium hover:bg-yellow-200 transition"
  >
    {showTopContributorsChart ? "Hide Top Contributors Graph" : "Show Top Contributors Graph"}
  </button>

</div>
{showContributionChart && (
  <ContributionChart contributors={contributors} />
)}

{showTopContributorsChart && (
  <TopContributorsChart contributors={contributors} />
)}
      {selectedContributor && (
        <div className="m-6 rounded-xl border border-blue-500/30 bg-slate-950 p-6">
          <div className="flex items-center justify-between">
            <h3 className="text-xl font-bold text-white">
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
              <div key={label} className="rounded-lg bg-slate-900 p-4">
                <p className="text-sm text-slate-400">{label}</p>

                <p className="mt-2 break-words text-xl font-bold text-white">
                  {value ?? 0}
                </p>
              </div>
            ))}
          </div>
        </div>
      )}
      <div className="flex flex-col gap-4 border-b border-slate-800 p-6 md:flex-row md:items-center md:justify-between">
        <div>
          <h2 className="text-xl font-bold text-white">Contributor Rankings</h2>

          <p className="mt-1 text-sm text-slate-400">
            Overview of developer contributions
          </p>
        </div>

        <div className="flex flex-col gap-3 sm:flex-row">
          <div className="flex items-center gap-3 rounded-lg border border-slate-700 bg-slate-950 px-4 py-3">
            <Search size={18} className="text-slate-400" />

            <input
              type="text"
              placeholder="Search contributor..."
              value={search}
              onChange={(e) => {
                setSearch(e.target.value);
                setCurrentPage(1);
              }}
              className="w-full bg-transparent text-sm text-white outline-none placeholder:text-slate-500"
            />
          </div>

          <button
            onClick={exportCSV}
            className="rounded-lg bg-blue-600 px-5 py-3 text-sm font-semibold text-white transition hover:bg-blue-500"
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
          className="rounded-lg border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-white outline-none"
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
      <table className="w-full text-left">
        <thead className="bg-slate-800/60 text-sm text-slate-400">
          <tr>
            <th className="px-6 py-4">Rank</th>
            <th className="px-6 py-4">Contributor</th>
            <th className="px-6 py-4">Commits</th>
            <th className="px-6 py-4">Code Changes</th>
            <th className="px-6 py-4">Pull Requests</th>
            <th className="px-6 py-4">Fair Score</th>
          </tr>
        </thead>

        <tbody>
          {currentContributors.length > 0 ? (
            currentContributors.map((contributor) => {
              const rank = sortedContributors.indexOf(contributor);

              return (
                <tr
                  key={contributor.contributor}
                  className="border-t border-slate-800 text-sm text-slate-300 hover:bg-slate-800/40"
                >
                  <td className="px-6 py-4">#{rank + 1}</td>

                  <td className="px-6 py-4">
                    <button
                      onClick={() => setSelectedContributor(contributor)}
                      className="font-medium text-blue-400 hover:text-blue-300 hover:underline"
                    >
                      {contributor.contributor}
                    </button>
                  </td>

                  <td className="px-6 py-4">{contributor.commits ?? 0}</td>

                  <td className="px-6 py-4">
                    {contributor.total_changes ?? 0}
                  </td>

                  <td className="px-6 py-4">{contributor.total_prs ?? 0}</td>

                  <td className="px-6 py-4">
                    <span
                      className={`rounded-full px-3 py-1 text-xs font-semibold ${getScoreStyle(
                        contributor.score ?? 0,
                      )}`}
                    >
                      {contributor.score ?? 0}/100
                    </span>
                  </td>
                </tr>
              );
            })
          ) : (
            <tr>
              <td colSpan="6" className="px-6 py-10 text-center text-slate-500">
                No contributors found.
              </td>
            </tr>
          )}
        </tbody>
      </table>

      {/* Pagination */}
      <div className="flex items-center justify-between border-t border-slate-800 p-5">
        <p className="text-sm text-slate-400">
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
