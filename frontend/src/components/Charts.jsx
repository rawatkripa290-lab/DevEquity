
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

function Charts({ contributors = [] }) {
  const chartData = contributors.map((contributor) => ({
    name: contributor.contributor,
    commits: contributor.commits ?? 0,
    score: contributor.score ?? 0,
  }));

  const chartStyle = {
    backgroundColor: "#0f172a",
    border: "1px solid #334155",
    borderRadius: "8px",
    color: "#fff",
  };

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div className="rounded-xl border border-slate-800 bg-slate-900 p-6">
        <h2 className="text-xl font-bold text-white mb-6">
          Contributor Activity
        </h2>

        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
              <XAxis dataKey="name" stroke="#94a3b8" />
              <YAxis stroke="#94a3b8" />
              <Tooltip contentStyle={chartStyle} />
              <Bar dataKey="commits" fill="#3b82f6" radius={[6, 6, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      <div className="rounded-xl border border-slate-800 bg-slate-900 p-6">
        <h2 className="text-xl font-bold text-white mb-6">
          Fair Contribution Scores
        </h2>

        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
              <XAxis dataKey="name" stroke="#94a3b8" />
              <YAxis domain={[0, 100]} stroke="#94a3b8" />
              <Tooltip contentStyle={chartStyle} />
              <Bar dataKey="score" fill="#22c55e" radius={[6, 6, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
}

export default Charts;
