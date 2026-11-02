
import axios from "axios";

const API = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || "http://127.0.0.1:8000/api",
  headers: {
    "Content-Type": "application/json"
  }
});

export async function analyzeRepository(repoUrl) {
  const { data } = await API.post("/analyze", {
    repo_url: repoUrl
  });

  return data;
}

export default API;
