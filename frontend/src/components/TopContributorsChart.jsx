
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

function TopContributorsChart({ contributors = [] }) {
  const data = [...contributors]
    .sort((a, b) => (b.score ?? 0) - (a.score ?? 0))
    .slice(0, 5)
    .map((contributor) => ({
      name: contributor.contributor,
      score: contributor.score ?? 0,
    }));

  return (
    <div className="mt-6 rounded-xl border  bg-yellow-200 p-6">
      <h2 className="mb-4 text-xl font-bold text-blue-900">
        Top 5 Contributors
      </h2>

      <div className="h-80 w-full">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={data}>
            <CartesianGrid strokeDasharray="3 3" stroke="#334155" />

            <XAxis
              dataKey="name"
              tick={{ fill: "#94a3b8", fontSize: 12 }}
            />

            <YAxis
              domain={[0, 100]}
              tick={{ fill: "#94a3b8" }}
            />

            <Tooltip />

            <Bar
              dataKey="score"
              name="Fair Score"
              fill="#3b82f6"
              radius={[6, 6, 0, 0]}
            />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}

export default TopContributorsChart;
