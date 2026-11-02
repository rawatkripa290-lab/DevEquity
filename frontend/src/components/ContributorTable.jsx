
function ContributorTable({ contributors = [] }) {
  return (
    <div className="w-full overflow-x-auto rounded-xl border border-white bg-white">
      <div className="border-b border-blue-800 p-6">
        <h2 className="text-xl font-bold text-white">
          Contributor Rankings
        </h2>
        <p className="mt-1 text-sm text-slate-400">
          Overview of developer contributions
        </p>
      </div>

      <table className="w-full text-left">
        <thead className="bg-white text-sm text-slate-400">
          <tr>
            <th className="px-6 py-4">Rank</th>
            <th className="px-6 py-4">Contributor</th>
            <th className="px-6 py-4">Commits</th>
            <th className="px-6 py-4">Contribution Score</th>
          </tr>
        </thead>

        <tbody>
          {contributors.length > 0 ? (
            contributors.map((contributor, index) => (
              <tr
                key={contributor.contributor || index}
                className="border-t border-slate-800 text-sm text-slate-300 hover:bg-slate-800/40"
              >
                <td className="px-6 py-4">{index + 1}</td>
                <td className="px-6 py-4 font-medium text-white">
                  {contributor.contributor}
                </td>
                <td className="px-6 py-4">
                  {contributor.commits ?? 0}
                </td>
                <td className="px-6 py-4">
                  {contributor.score ?? 0}
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td
                colSpan="4"
                className="px-6 py-10 text-center text-slate-500"
              >
                No contributor data available yet.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

export default ContributorTable;
