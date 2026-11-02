
import { useState } from "react";
import { GitBranch, Search, LoaderCircle } from "lucide-react";

function RepoInput({ onAnalyze, loading }) {
  const [repoUrl, setRepoUrl] = useState("");

  const handleAnalyze = () => {
    if (!repoUrl.trim()) {
      alert("Please enter a GitHub repository URL");
      return;
    }

    onAnalyze(repoUrl);
  };

  return (
    <div className="w-full max-w-2xl mx-auto px-6">
      <div className="bg-blue-400 border border-blue-400 rounded-2xl p-6 md:p-8">
        <label className="block text-white font-semibold mb-3">
          Enter GitHub Repository URL
        </label>

        <div className="flex items-center gap-3 bg-slate-200 border border-white rounded-xl px-4 py-4">
          <GitBranch size={22} className="text-slate-800" />
          <input
            type="text"
            value={repoUrl}
            onChange={(e) => setRepoUrl(e.target.value)}
            placeholder="https://github.com/owner/repository"
            className="w-full bg-transparent text-black outline-none placeholder:text-slate-800"
          />
        </div>
        <button
          onClick={handleAnalyze}
          disabled={loading}
          className="w-full mt-5 bg-blue-900 hover:bg-blue-500 disabled:opacity-50 text-white font-semibold py-4 rounded-xl flex items-center justify-center gap-3 transition"
        >
          {loading ? (
            <>
              <LoaderCircle size={20} className="animate-spin" />
              Analyzing Repository...
            </>
          ) : (
            <>
              <Search size={20} />
              Analyze Repository
            </>
          )}
        </button>
      </div>
    </div>
  );
}

export default RepoInput;
