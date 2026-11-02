
import axios from "axios";

const API = axios.create({
  baseURL: "http://127.0.0.1:8000/api",
  headers: {
    "Content-Type": "application/json",
  },
});

export const analyzeRepository = async (repoUrl) => {
  try {
    const response = await API.post("/analyze", {
      repo_url: repoUrl,
    });

    return response.data;
  } catch (error) {
    console.error("Repository analysis failed:", error);

    throw error;
  }
};

export default API;
