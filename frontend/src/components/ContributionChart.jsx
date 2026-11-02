
import {
  PieChart,
  Pie,
  Cell,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";

function ContributionChart({ contributors = [] }) {
  const excellent = contributors.filter(
    (c) => c.score >= 75
  ).length;

  const good = contributors.filter(
    (c) => c.score >= 40 && c.score < 75
  ).length;

  const needsImprovement = contributors.filter(
    (c) => c.score < 40
  ).length;

  const data = [
    { name: "Excellent", value: excellent },
    { name: "Good", value: good },
    { name: "Needs Improvement", value: needsImprovement },
  ];

  const COLORS = ["#22c55e", "#eab308", "#ef4444"];

  return (
    <div className="mt-6 rounded-xl border border-slate-800 bg-blue-900 p-6">
      <h2 className="mb-4 text-xl font-bold text-white">
        Contribution Distribution
      </h2>

      <div className="h-80 w-full">
        <ResponsiveContainer width="100%" height="100%">
          <PieChart>
            <Pie
              data={data}
              cx="50%"
              cy="50%"
              outerRadius={100}
              dataKey="value"
              nameKey="name"
              label
            >
              {data.map((entry, index) => (
                <Cell
                  key={entry.name}
                  fill={COLORS[index]}
                />
              ))}
            </Pie>

            <Tooltip />
            <Legend />
          </PieChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}

export default ContributionChart;