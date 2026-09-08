import axios from "axios";

// 接后端时改为 false；页面组件不需要改动。
const USE_MOCK = false;
const client = axios.create({
  baseURL: "http://localhost:8080",
  timeout: 40000,
  headers: { "Content-Type": "application/json" },
});

export const api = {
  async listHosts() {
    const { data } = await client.get("/hosts");
    return Array.isArray(data) ? data : [];
  },
  async getHosts() {
    const { data } = await client.get("/hostlist");
    return Array.isArray(data) ? data : [];
  },
  async getItems(host) {
    const { data } = await client.get("/items", {
      params: { host },
    });
    return data;
  },
  async createHost(payload) {
    if (USE_MOCK) return { ...payload };
    const response = await client.post("/post/data", payload);
    return response;
  },
  async checkHosts() {
    const response = await client.post("/check");
    return response;
  },
  async getMetrics({ host, days = 15, metric = "cpu,memory" }) {
    const path = host
      ? `/hosts/${encodeURIComponent(host)}/metrics`
      : "/metrics";
    const { data } = await client.get(path, { params: { days, metric } });
    return data;
  },
};
