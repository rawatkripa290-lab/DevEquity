
import axios from "axios";

const API = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || "http://127.0.0.1:8080/api/v1",
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
export default API;
