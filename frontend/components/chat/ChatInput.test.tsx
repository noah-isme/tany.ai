import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

import { ChatInput } from "./ChatInput";

describe("ChatInput", () => {
  it("memanggil onSend ketika form disubmit", async () => {
    const handleSend = vi.fn();
    const user = userEvent.setup();

    render(<ChatInput onSend={handleSend} />);

    const textarea = screen.getByPlaceholderText(
      /tanyakan layanan, harga/i,
    ) as HTMLTextAreaElement;
    await user.type(textarea, "Apa layanan andalanmu?");

    await user.click(screen.getByRole("button", { name: /kirim/i }));

    expect(handleSend).toHaveBeenCalledWith("Apa layanan andalanmu?");
    expect(textarea.value).toBe("");
  });
});
