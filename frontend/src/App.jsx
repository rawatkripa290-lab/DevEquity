
import { useState } from "react";
import Navbar from "./components/NavBar";
import Hero from "./components/Hero";
import RepoInput from "./components/RepoInput";
import Dashboard from "./components/Dashboard";
import About from "./components/About";
import { analyzeRepository } from "./services/api";

function App() {
  const [analysisData, setAnalysisData] = useState(null);
  const [status, setStatus] = useState("idle");
  const [error, setError] = useState("");

  const loading = status === "processing";

  const handleAnalyze = async (repoUrl) => {
    setError("");
    setAnalysisData(null);
    setStatus("processing");

    try {
      const result = await analyzeRepository(repoUrl);

      console.log("Backend Response:", result);

      if (!result || typeof result !== "object" || Object.keys(result).length === 0) {
        throw new Error("The backend returned empty repository data.");
      }

      setAnalysisData(result);
      setStatus("completed");

    } catch (err) {
      console.error("Analysis Error:", err.response?.data || err);

      const errorMessage =
        err.response?.data?.detail ||
        err.response?.data?.error ||
        err.message ||
        "Unable to analyze repository. Please try again.";

      setStatus("failed");
      setError(errorMessage);
    }
  };

  return (
    <div className="min-h-screen pixel-desktop">
      <Navbar />

      {!analysisData && (
        <div id="home">
          <Hero>
            <RepoInput
              onAnalyze={handleAnalyze}
              loading={loading}
            />
          </Hero>
        </div>
      )}

      {loading && (
        <p className="mt-6 text-center">
          Analyzing repository... Please wait.
        </p>
      )}

      {error && (
        <div
          className="fixed top-24 right-6 z-[9999] w-[90%] max-w-sm rounded-xl border border-red-200 bg-white p-4 shadow-2xl"
          role="alert"
        >
          <div className="flex items-start gap-3">
            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-red-100 text-red-600">
              !
            </div>

            <div>
              <h3 className="font-bold text-red-700">
                Analysis Failed
              </h3>
              <p className="mt-1 text-sm text-slate-600">
                {error}
              </p>
            </div>

            <button
              onClick={() => setError("")}
              className="ml-auto text-xl text-slate-400 hover:text-slate-700"
              aria-label="Close error notification"
            >
              &times;
            </button>
          </div>
        </div>
      )}

      {analysisData && (
        <Dashboard analysisData={analysisData} />
      )}

      {!analysisData && <About />}
    </div>
  );
}

export default App;
