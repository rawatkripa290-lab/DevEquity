import {
  GitBranch,
  Scale,
  ChartNoAxesCombined,
  Users
} from "lucide-react";

const highlights = [
  {
    icon: Scale,
    title: "Fair Contribution Scoring",
    description:
      "Evaluate developers using multiple contribution metrics instead of relying only on commit counts.",
  },
  {
    icon: ChartNoAxesCombined,
    title: "Meaningful Analytics",
    description:
      "Visualize commits, code changes, pull requests, reviews and other repository activities.",
  },
  {
    icon: Users,
    title: "Team Transparency",
    description:
      "Understand individual contributions and compare participation across your development team.",
  },
  {
    icon: GitBranch,
    title: "GitHub Repository Analysis",
    description:
      "Analyze repository activity using GitHub data and generate an organized contribution report.",
  },
];

function About() {
  return (
    <section id="about" className="scroll-mt-24 bg-slate-50 px-6 py-20 md:px-12">
      <div className="mx-auto max-w-6xl">

        <div className="mb-12 text-center">
          <span className=" px-4 py-2 text-2xl font-extrabold text-blue-700">
            ABOUT DEVEQUITY
          </span>

          <h2 className="mt-5 text-3xl font-extrabold text-slate-900 md:text-xl">
            Beyond Commits. Recognizing Real Contributions.
          </h2>

          <p className="mx-auto mt-4 max-w-2xl text-slate-900">
            DevEquity is a GitHub repository analytics platform designed
            to provide a more balanced and transparent understanding
            of developer contributions.
          </p>
        </div>

        <div className="grid items-center gap-10 md:grid-cols-1">
          <div className="rounded-2xl border border-blue-100 bg-white p-8 shadow-sm">
          

            <h3 className="mb-4 text-2xl font-bold text-slate-900">
              Why DevEquity?
            </h3>

            <p className="leading-7 text-slate-800">
              Software development is a collaborative process.
              A developer's contribution involves more than just
              writing commits. Code changes, pull requests, reviews,
              issue handling and consistency also matter.
            </p>

            <p className="mt-4 leading-7 text-slate-800">
              DevEquity brings these activities together to help
              teams understand contribution patterns and make
              more informed evaluations.
            </p>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            {highlights.map((item) => {

              return (
                <div
                  key={item.title}
                  className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm transition hover:-translate-y-1 hover:shadow-md"
                >
                  

                  <h3 className="mb-2 font-bold text-slate-900">
                    {item.title}
                  </h3>

                  <p className="text-sm leading-6 text-slate-600">
                    {item.description}
                  </p>
                </div>
              );
            })}
          </div>
        </div>

        <div className="mt-12 rounded-2xl bg-blue-700 px-6 py-10 text-center text-white">
          <GitBranch className="mx-auto mb-4" size={34} />

          <h3 className="text-xl font-bold">
            Every Contribution Matters.
          </h3>

          <p className="mx-auto mt-3 max-w-xl text-blue-100">
            Building a clearer picture of teamwork, effort and
            collaboration in software development.
          </p>
        </div>

      </div>
    </section>
  );
}

export default About;
