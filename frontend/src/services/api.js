
import axios from "axios";

const API = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api/v1",
  headers: {
    "Content-Type": "application/json"
  }
});

export async function analyzeRepository(repoUrl) {
  const { data } = await API.post("/analyses", {
    repository_url: repoUrl
  });

  return data;
}

export async function getAnalysis(id) {
  const { data } = await API.get(`/analyses/${id}`);

  // The Go handler currently returns progress as JSON text.
  if (typeof data.progress === "string") {
    try {
      data.progress = JSON.parse(data.progress);
    } catch {
      data.progress = {};
    }
  }

  return data;
}

export async function getReport(id) {
  const { data } = await API.get(`/analyses/${id}/report`);
  return data;
}

export default API;
