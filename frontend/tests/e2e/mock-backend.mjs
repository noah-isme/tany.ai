import { createServer } from "node:http";
import { randomUUID } from "node:crypto";

const PORT = Number(process.env.MOCK_API_PORT ?? 4000);

const skills = [
  { id: randomUUID(), name: "Next.js", order: 0 },
  { id: randomUUID(), name: "Golang", order: 1 },
];

const profile = {
  id: "profile-1",
  name: "Tanya Admin",
  title: "Lead AI Consultant",
  bio: "Membantu tim membangun knowledge base tany.ai.",
  email: "admin@example.com",
  phone: "+62 812-0000-0000",
  location: "Jakarta, Indonesia",
  avatar_url: "https://avatars.dicebear.com/api/initials/TA.svg",
  updated_at: new Date().toISOString(),
};

const services = [
  {
    id: randomUUID(),
    name: "AI Discovery Workshop",
    description: "Eksplorasi kebutuhan dan peluang otomasi AI untuk bisnis.",
    price_min: 15000000,
    price_max: 25000000,
    currency: "IDR",
    duration_label: "2 minggu",
    is_active: true,
    order: 0,
  },
  {
    id: randomUUID(),
    name: "Chatbot Optimization",
    description: "Audit dan optimasi alur percakapan untuk peningkatan konversi.",
    price_min: 10000000,
    price_max: 18000000,
    currency: "IDR",
    duration_label: "10 hari",
    is_active: true,
    order: 1,
  },
];

const projects = [
  {
    id: randomUUID(),
    title: "Enterprise Knowledge Assistant",
    description: "Implementasi chatbot internal untuk knowledge base perusahaan.",
    tech_stack: ["Next.js", "LangChain", "Postgres"],
    image_url: "https://images.unsplash.com/photo-1523475472560-d2df97ec485c?auto=format&fit=crop&w=640&q=80",
    project_url: "https://example.com/enterprise-assistant",
    category: "AI Platform",
    order: 0,
    is_featured: true,
  },
  {
    id: randomUUID(),
    title: "Sales Enablement Chatbot",
    description: "Automasi Q&A produk untuk tim sales global.",
    tech_stack: ["Next.js", "OpenAI", "Supabase"],
    image_url: "https://images.unsplash.com/photo-1525182008055-f88b95ff7980?auto=format&fit=crop&w=640&q=80",
    project_url: "https://example.com/sales-bot",
    category: "Sales Enablement",
    order: 1,
    is_featured: false,
  },
];

function sendJson(res, status, body) {
  const payload = JSON.stringify(body);
  res.writeHead(status, {
    "Content-Type": "application/json",
    "Access-Control-Allow-Origin": "*",
  });
  res.end(payload);
}

function parseBody(req) {
  return new Promise((resolve, reject) => {
    let data = "";
    req.on("data", (chunk) => {
      data += chunk;
    });
    req.on("end", () => {
      try {
        resolve(data ? JSON.parse(data) : {});
      } catch (error) {
        reject(error);
      }
    });
    req.on("error", reject);
  });
}

const payload = Buffer.from(
  JSON.stringify({
    sub: "admin",
    email: "admin@example.com",
    roles: ["admin"],
    exp: Math.floor(Date.now() / 1000) + 3600,
  }),
).toString("base64url");

const token = `header.${payload}.signature`;

const server = createServer(async (req, res) => {
  const url = req.url ?? "";
  if (req.method === "OPTIONS") {
    res.writeHead(204, {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET,POST,PATCH,PUT,DELETE,OPTIONS",
      "Access-Control-Allow-Headers": "Content-Type,Authorization",
    });
    res.end();
    return;
  }

  if (req.method === "POST" && url === "/api/auth/login") {
    sendJson(res, 200, {
      accessToken: token,
      user: { id: "admin", email: "admin@example.com", roles: ["admin"] },
    });
    return;
  }

  if (req.method === "POST" && url === "/api/auth/logout") {
    res.writeHead(204);
    res.end();
    return;
  }

  if (req.method === "GET" && url === "/api/admin/profile") {
    sendJson(res, 200, { data: profile });
    return;
  }

  if (req.method === "PUT" && url === "/api/admin/profile") {
    const body = await parseBody(req);
    Object.assign(profile, body);
    profile.updated_at = new Date().toISOString();
    sendJson(res, 200, { data: profile });
    return;
  }

  if (req.method === "GET" && url.startsWith("/api/admin/services")) {
    const sorted = services.slice().sort((a, b) => a.order - b.order);
    sendJson(res, 200, { items: sorted, page: 1, limit: sorted.length || 1, total: sorted.length });
    return;
  }

  if (req.method === "GET" && url.startsWith("/api/admin/projects")) {
    const sorted = projects.slice().sort((a, b) => a.order - b.order);
    sendJson(res, 200, { items: sorted, page: 1, limit: sorted.length || 1, total: sorted.length });
    return;
  }

  if (req.method === "GET" && url.startsWith("/api/admin/skills")) {
    const sorted = skills.slice().sort((a, b) => a.order - b.order);
    sendJson(res, 200, { items: sorted, page: 1, limit: sorted.length, total: sorted.length });
    return;
  }

  if (req.method === "POST" && url === "/api/admin/skills") {
    const body = await parseBody(req);
    const newSkill = { id: randomUUID(), name: body.name ?? "", order: skills.length };
    skills.push(newSkill);
    sendJson(res, 201, { data: newSkill });
    return;
  }

  if (req.method === "PATCH" && url === "/api/admin/skills/reorder") {
    const body = await parseBody(req);
    body.forEach((item) => {
      const target = skills.find((skill) => skill.id === item.id);
      if (target) {
        target.order = item.order;
      }
    });
    res.writeHead(204);
    res.end();
    return;
  }

  if (req.method === "DELETE" && url.startsWith("/api/admin/skills/")) {
    const id = url.split("/").pop();
    const index = skills.findIndex((skill) => skill.id === id);
    if (index >= 0) {
      skills.splice(index, 1);
    }
    res.writeHead(204);
    res.end();
    return;
  }

  res.writeHead(404);
  res.end();
});

server.listen(PORT, () => {
  console.log(`[mock-backend] listening on http://localhost:${PORT}`);
});
