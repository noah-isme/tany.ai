import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { vi } from "vitest";

import { LoginForm } from "./LoginForm";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ replace: vi.fn(), refresh: vi.fn() }),
}));

const loginActionMock = vi.fn();

vi.mock("./actions", () => ({
  loginAction: (formData: FormData) => loginActionMock(formData),
}));

describe("LoginForm", () => {
  it("validates email field", async () => {
    loginActionMock.mockResolvedValue({ success: false, error: "Invalid" });
    render(<LoginForm />);

    await userEvent.type(screen.getByLabelText(/email/i), "wrong");
    await userEvent.type(screen.getByLabelText(/password/i), "password123");
    await userEvent.click(screen.getByRole("button", { name: /masuk/i }));

    expect(await screen.findByText(/email tidak valid/i)).toBeInTheDocument();
    expect(loginActionMock).not.toHaveBeenCalled();
  });

  it("shows server error", async () => {
    loginActionMock.mockResolvedValue({ success: false, error: "Email atau password salah." });
    render(<LoginForm />);

    await userEvent.type(screen.getByLabelText(/email/i), "admin@example.com");
    await userEvent.type(screen.getByLabelText(/password/i), "Password123!");
    await userEvent.click(screen.getByRole("button", { name: /masuk/i }));

    await waitFor(() => {
      expect(loginActionMock).toHaveBeenCalled();
    });

    expect(await screen.findByText(/email atau password salah/i)).toBeInTheDocument();
  });
});
