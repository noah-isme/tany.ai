import { createServer } from "node:http";
import { createHmac, randomUUID } from "node:crypto";

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
    duration_label: "10 minggu",
    price_label: "Enterprise retainer",
    budget_label: "IDR 220Jt",
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
    duration_label: "6 minggu",
    price_label: "Launch package",
    budget_label: "IDR 120Jt",
    order: 1,
    is_featured: false,
  },
];

function toNumberOrNull(value) {
  if (typeof value === "number") {
    return value;
  }
  if (value === null || value === undefined || value === "") {
    return null;
  }
  const parsed = Number(value);
  return Number.isNaN(parsed) ? null : parsed;
}

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

const JWT_SECRET = process.env.JWT_SECRET ?? "test-secret";

function signToken(payload) {
  const header = Buffer.from(JSON.stringify({ alg: "HS256", typ: "JWT" })).toString("base64url");
  const body = Buffer.from(JSON.stringify(payload)).toString("base64url");
  const signature = createHmac("sha256", JWT_SECRET).update(`${header}.${body}`).digest("base64url");
  return `${header}.${body}.${signature}`;
}

const token = signToken({
  sub: "admin",
  email: "admin@example.com",
  roles: ["admin"],
  exp: Math.floor(Date.now() / 1000) + 3600,
});

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

  if (req.method === "GET" && url === "/api/v1/knowledge-base") {
    const payload = {
      profile: {
        name: profile.name,
        title: profile.title,
        bio: profile.bio,
        email: profile.email,
        phone: profile.phone,
        location: profile.location,
        avatarUrl: profile.avatar_url,
        updatedAt: profile.updated_at,
      },
      skills: skills
        .slice()
        .sort((a, b) => a.order - b.order)
        .map((skill) => ({ name: skill.name })),
      services: services
        .filter((service) => service.is_active)
        .sort((a, b) => a.order - b.order)
        .map((service) => {
          const priceRange = [];
          if (service.price_min !== null && service.price_min !== undefined) {
            priceRange.push(`${service.currency ?? ""} ${service.price_min}`.trim());
          }
          if (
            service.price_max !== null &&
            service.price_max !== undefined &&
            service.price_max !== service.price_min
          ) {
            priceRange.push(`${service.currency ?? ""} ${service.price_max}`.trim());
          }
          return {
            id: service.id,
            name: service.name,
            description: service.description,
            currency: service.currency,
            durationLabel: service.duration_label,
            priceRange,
            order: service.order,
          };
        }),
      projects: projects
        .slice()
        .sort((a, b) => {
          if (a.is_featured === b.is_featured) {
            return a.order - b.order;
          }
          return a.is_featured ? -1 : 1;
        })
        .map((project) => ({
          id: project.id,
          title: project.title,
          description: project.description,
          techStack: project.tech_stack,
          projectUrl: project.project_url,
          category: project.category,
          durationLabel: project.duration_label,
          priceLabel: project.price_label,
          budgetLabel: project.budget_label,
          isFeatured: project.is_featured,
          order: project.order,
        })),
    };

    res.writeHead(200, {
      "Content-Type": "application/json",
      "Cache-Control": "public, max-age=60",
      ETag: 'W/"mock-etag"',
    });
    res.end(JSON.stringify(payload));
    return;
  }

  if (req.method === "POST" && url === "/api/v1/chat") {
    const body = await parseBody(req);
    const id = body.chatId && typeof body.chatId === "string" ? body.chatId : randomUUID();
    const answer = `Pertanyaan diterima: ${body.question}. Layanan aktif kami meliputi ${services
      .filter((service) => service.is_active)
      .map((service) => service.name)
      .join(", ")}.`;
    sendJson(res, 200, {
      chatId: id,
      answer,
      model: "mock-model",
      prompt: "mock-prompt",
    });
    return;
  }

  if (req.method === "POST" && url === "/api/admin/uploads") {
    req.on("data", () => {});
    req.on("end", () => {
      sendJson(res, 201, {
        data: {
          url: `https://mock-storage.example.com/uploads/${randomUUID()}.png`,
          key: `uploads/${randomUUID()}.png`,
          contentType: "image/png",
          size: 2048,
        },
      });
    });
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

  if (req.method === "POST" && url === "/api/admin/services") {
    const body = await parseBody(req);
    const newService = {
      id: randomUUID(),
      name: body.name ?? "",
      description: body.description ?? "",
      price_min: toNumberOrNull(body.price_min),
      price_max: toNumberOrNull(body.price_max),
      currency: body.currency ?? "",
      duration_label: body.duration_label ?? "",
      is_active: typeof body.is_active === "boolean" ? body.is_active : true,
      order: services.length,
    };
    services.push(newService);
    sendJson(res, 201, { data: newService });
    return;
  }

  if (req.method === "PUT" && url.startsWith("/api/admin/services/")) {
    const [, , , , serviceId] = url.split("/");
    const body = await parseBody(req);
    const target = services.find((service) => service.id === serviceId);
    if (target) {
      Object.assign(target, {
        name: body.name ?? target.name,
        description: body.description ?? target.description,
        price_min: body.price_min !== undefined ? toNumberOrNull(body.price_min) : target.price_min,
        price_max: body.price_max !== undefined ? toNumberOrNull(body.price_max) : target.price_max,
        currency: body.currency ?? target.currency,
        duration_label: body.duration_label ?? target.duration_label,
        is_active: typeof body.is_active === "boolean" ? body.is_active : target.is_active,
      });
      sendJson(res, 200, { data: target });
      return;
    }
    sendJson(res, 404, { error: { message: "Service not found" } });
    return;
  }

  if (req.method === "PATCH" && url === "/api/admin/services/reorder") {
    const body = await parseBody(req);
    body.forEach((item) => {
      const target = services.find((service) => service.id === item.id);
      if (target) {
        target.order = item.order;
      }
    });
    res.writeHead(204);
    res.end();
    return;
  }

  if (req.method === "PATCH" && url.match(/^\/api\/admin\/services\/.+\/toggle$/)) {
    const serviceId = url.split("/")[4];
    const body = await parseBody(req);
    const target = services.find((service) => service.id === serviceId);
    if (target) {
      const nextState =
        typeof body.is_active === "boolean" ? body.is_active : !target.is_active;
      target.is_active = nextState;
      sendJson(res, 200, { data: target });
      return;
    }
    sendJson(res, 404, { error: { message: "Service not found" } });
    return;
  }

  if (req.method === "DELETE" && url.startsWith("/api/admin/services/")) {
    const serviceId = url.split("/").pop();
    const index = services.findIndex((service) => service.id === serviceId);
    if (index >= 0) {
      services.splice(index, 1);
    }
    res.writeHead(204);
    res.end();
    return;
  }

  if (req.method === "POST" && url === "/api/admin/projects") {
    const body = await parseBody(req);
    const newProject = {
      id: randomUUID(),
      title: body.title ?? "",
      description: body.description ?? "",
      tech_stack: Array.isArray(body.tech_stack) ? body.tech_stack : [],
      image_url: body.image_url ?? "",
      project_url: body.project_url ?? "",
      category: body.category ?? "",
      duration_label: body.duration_label ?? "",
      price_label: body.price_label ?? "",
      budget_label: body.budget_label ?? "",
      order: projects.length,
      is_featured: typeof body.is_featured === "boolean" ? body.is_featured : false,
    };
    projects.push(newProject);
    sendJson(res, 201, { data: newProject });
    return;
  }

  if (req.method === "PUT" && url.startsWith("/api/admin/projects/")) {
    const [, , , , projectId] = url.split("/");
    const body = await parseBody(req);
    const target = projects.find((project) => project.id === projectId);
    if (target) {
      Object.assign(target, {
        title: body.title ?? target.title,
        description: body.description ?? target.description,
        tech_stack: Array.isArray(body.tech_stack) ? body.tech_stack : target.tech_stack,
        image_url: body.image_url ?? target.image_url,
        project_url: body.project_url ?? target.project_url,
        category: body.category ?? target.category,
        duration_label: body.duration_label ?? target.duration_label,
        price_label: body.price_label ?? target.price_label,
        budget_label: body.budget_label ?? target.budget_label,
        is_featured: typeof body.is_featured === "boolean" ? body.is_featured : target.is_featured,
      });
      sendJson(res, 200, { data: target });
      return;
    }
    sendJson(res, 404, { error: { message: "Project not found" } });
    return;
  }

  if (req.method === "PATCH" && url === "/api/admin/projects/reorder") {
    const body = await parseBody(req);
    body.forEach((item) => {
      const target = projects.find((project) => project.id === item.id);
      if (target) {
        target.order = item.order;
      }
    });
    res.writeHead(204);
    res.end();
    return;
  }

  if (req.method === "PATCH" && url.match(/^\/api\/admin\/projects\/.+\/feature$/)) {
    const projectId = url.split("/")[4];
    const body = await parseBody(req);
    const target = projects.find((project) => project.id === projectId);
    if (target) {
      target.is_featured = typeof body.is_featured === "boolean" ? body.is_featured : true;
      sendJson(res, 200, { data: target });
      return;
    }
    sendJson(res, 404, { error: { message: "Project not found" } });
    return;
  }

  if (req.method === "DELETE" && url.startsWith("/api/admin/projects/")) {
    const projectId = url.split("/").pop();
    const index = projects.findIndex((project) => project.id === projectId);
    if (index >= 0) {
      projects.splice(index, 1);
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
