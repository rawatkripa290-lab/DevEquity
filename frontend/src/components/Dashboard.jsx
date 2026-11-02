import StatsCard from "./StatsCard";
import ContributionTable from "./ContributionTable";
import Charts from "./Charts";
import RepositoryOverview from "./RepositoryOverview";
import ContributionSummary from "./ContributionSummary";
function Dashboard({ analysisData }) {
  const contributors = analysisData?.contributors || [];
  const fairScores = analysisData?.fair_scores || [];
  console.log(
    "Repository Details:",
    JSON.stringify(analysisData?.data, null, 2),
  );

  const mergedContributors = contributors
    .map((contributor) => {
      const scoreData = fairScores.find(
        (score) => score.contributor === contributor.contributor,
      );

      return {
        ...contributor,
        commits: contributor.total_commits,
        score: scoreData?.fair_score ?? 0,
      };
    })
    .sort((a, b) => b.score - a.score);

  return (
    <main className="min-h-screen bg-blue-200 px-6 py-10 md:px-12">
      <div className="mx-auto max-w-7xl">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-blue-900 text-center">
            Repository Dashboard
          </h1>
          <p className="mt-2 text-blue-900 text-center font-medium ">
            Explore repository activity and contributor performance.
          </p>
        </div>
        <RepositoryOverview repository={analysisData?.data} />
        <StatsCard
          totalCommits={analysisData?.total_commits_fetched}
          totalContributors={contributors.length}
          totalPRs={analysisData?.pull_request_metrics?.length}
          totalIssues={analysisData?.issue_metrics?.length}
        />
        <div className="mt-8">
          <ContributionSummary contributors={mergedContributors} />
        </div>

        <div className="mt-8">
          <Charts contributors={mergedContributors} />
        </div>

        <div className="mt-8">
          <ContributionTable contributors={mergedContributors} />
        </div>
      </div>
    </main>
  );
}

export default Dashboard;
