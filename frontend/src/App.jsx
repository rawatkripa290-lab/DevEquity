
import { useEffect, useState } from "react";
import Navbar from "./components/NavBar";
import Hero from "./components/Hero";
import RepoInput from "./components/RepoInput";
import Dashboard from "./components/Dashboard";
import {
  analyzeRepository,
  getAnalysis,
  getReport
} from "./services/api";

function App() {
  const [analysisData, setAnalysisData] = useState(null);
  const [analysisId, setAnalysisId] = useState(null);
  const [status, setStatus] = useState("idle");
  const [progress, setProgress] = useState({});
  const [error, setError] = useState("");

  const loading = ["queued", "processing"].includes(status);

  const handleAnalyze = async (repoUrl) => {
    setError("");
    setAnalysisData(null);
    setProgress({});
    setStatus("queued");

    try {
      const result = await analyzeRepository(repoUrl);
      setAnalysisId(result.analysis_id);
      setStatus(result.status || "queued");
    } catch (err) {
      setStatus("failed");
      setError(
        err.response?.data?.error ||
        "Unable to start repository analysis."
      );
    }
  };

  useEffect(() => {
    if (!analysisId || !loading) return;

    let active = true;
    let timer;

    const poll = async () => {
      try {
        const result = await getAnalysis(analysisId);

        if (!active) return;

        setStatus(result.status);
        setProgress(result.progress || {});

        if (result.status === "completed") {
          const report = await getReport(analysisId);
          if (active) setAnalysisData(report);
          return;
        }

        if (result.status === "failed") {
          setError(result.error || "Analysis failed.");
          return;
        }

        timer = setTimeout(poll, 3000);
      } catch (err) {
        if (!active) return;

        setError(
          err.response?.data?.error ||
          "Unable to retrieve analysis status."
        );
        setStatus("failed");
      }
    };

    poll();

    return () => {
      active = false;
      clearTimeout(timer);
    };
  }, [analysisId, loading]);

  return (
    <div className="min-h-screen pixel-desktop">
      <Navbar />

      {!analysisData && (
        <Hero>
          <RepoInput
            onAnalyze={handleAnalyze}
            loading={loading}
          />

          {loading && (
            <div className="pixel-window mt-8 mx-auto max-w-xl p-5">
              <h2 className="pixel-heading">ANALYSIS.EXE</h2>
              <p>
                {progress.stage || status}
              </p>
              <p>
                {progress.completed || 0} / {progress.total || "?"}
              </p>
              <div className="pixel-progress">
                <div
                  style={{
                    width: progress.total
                      ? `${(progress.completed / progress.total) * 100}%`
                      : "0%"
                  }}
                />
              </div>
            </div>
          )}
        </Hero>
      )}

      {error && (
        <p className="text-center mt-6 px-4" role="alert">
          {error}
        </p>
      )}

      {analysisData && (
        <Dashboard analysisData={analysisData} />
      )}
    </div>
  );
}

export default App;
