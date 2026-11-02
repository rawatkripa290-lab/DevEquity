import { useState } from "react";
import Navbar from "./components/NavBar";
import Hero from "./components/Hero";
import RepoInput from "./components/RepoInput";
import Dashboard from "./components/Dashboard";
import { analyzeRepository } from "./services/api";
function App() {
  const [analysisData, setAnalysisData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleAnalyze = async (repoUrl) => {
    setLoading(true);
    setError("");
    setAnalysisData(null);

    try {
      const result = await analyzeRepository(repoUrl);
      setAnalysisData(result);
    } catch (err) {
      setError(
        err.response?.data?.detail ||
        "Unable to analyze repository. Please try again."
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-white">
      <Navbar />

      {!analysisData && (
        <Hero>
          <RepoInput
            onAnalyze={handleAnalyze}
            loading={loading}
          />
        </Hero>
      )}

      {error && (
        <p className="text-center text-red-400 mt-6 px-4">
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
