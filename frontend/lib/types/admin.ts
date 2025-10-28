export type Profile = {
  id: string;
  name: string;
  title: string;
  bio: string;
  email: string;
  phone: string;
  location: string;
  avatar_url: string;
  updated_at?: string;
};

export type Skill = {
  id: string;
  name: string;
  order: number;
};

export type Service = {
  id: string;
  name: string;
  description: string;
  price_min: number | null;
  price_max: number | null;
  currency: string;
  duration_label: string;
  is_active: boolean;
  order: number;
};

export type Project = {
  id: string;
  title: string;
  description: string;
  tech_stack: string[];
  image_url: string;
  project_url: string;
  category: string;
  duration_label: string;
  price_label: string;
  budget_label: string;
  order: number;
  is_featured: boolean;
};

export type PaginatedResponse<T> = {
  items: T[];
  page: number;
  limit: number;
  total: number;
};

export type ApiListParams = {
  page?: number;
  limit?: number;
  sort?: string;
  dir?: "asc" | "desc";
};
