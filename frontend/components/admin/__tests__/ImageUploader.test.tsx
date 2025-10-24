import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { act, useState } from "react";
import { vi, beforeEach, afterEach, describe, it, expect } from "vitest";

import { ImageUploader } from "../ImageUploader";
import { ToastProvider } from "../ToastProvider";

class MockXMLHttpRequest {
  static instances: MockXMLHttpRequest[] = [];
  method = "";
  url = "";
  headers: Record<string, string> = {};
  responseType = "";
  withCredentials = false;
  status = 200;
  response: unknown = null;
  responseText = "";
  upload = { onprogress: null as ((event: ProgressEvent<EventTarget>) => void) | null };
  onload: (() => void) | null = null;
  onerror: (() => void) | null = null;

  open(method: string, url: string) {
    this.method = method;
    this.url = url;
  }

  setRequestHeader(name: string, value: string) {
    this.headers[name] = value;
  }

  send(body?: Document | BodyInit | null) {
    void body;
    MockXMLHttpRequest.instances.push(this);
  }

  abort() {}
}

const createFile = () => new File(["fake"], "avatar.png", { type: "image/png" });

function ControlledUploader({ onValueChange }: { onValueChange?: (value: string) => void }) {
  const [value, setValue] = useState("");
  const handleChange = (next: string) => {
    setValue(next);
    onValueChange?.(next);
  };
  return (
    <ToastProvider>
      <ImageUploader value={value} onChange={handleChange} />
    </ToastProvider>
  );
}

describe("ImageUploader", () => {
  const createObjectURL = vi.fn(() => "blob:preview");
  const revokeObjectURL = vi.fn();
  const hadCreateObjectURL = typeof URL.createObjectURL === "function";
  const hadRevokeObjectURL = typeof URL.revokeObjectURL === "function";

  beforeEach(() => {
    MockXMLHttpRequest.instances = [];
    vi.stubGlobal("XMLHttpRequest", MockXMLHttpRequest as unknown as typeof XMLHttpRequest);
    if (hadCreateObjectURL) {
      vi.spyOn(URL, "createObjectURL").mockImplementation(createObjectURL as unknown as typeof URL.createObjectURL);
    } else {
      Object.defineProperty(URL, "createObjectURL", { configurable: true, value: createObjectURL });
    }
    if (hadRevokeObjectURL) {
      vi.spyOn(URL, "revokeObjectURL").mockImplementation(revokeObjectURL as unknown as typeof URL.revokeObjectURL);
    } else {
      Object.defineProperty(URL, "revokeObjectURL", { configurable: true, value: revokeObjectURL });
    }
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
    if (!hadCreateObjectURL) {
      delete (URL as unknown as Record<string, unknown>).createObjectURL;
    }
    if (!hadRevokeObjectURL) {
      delete (URL as unknown as Record<string, unknown>).revokeObjectURL;
    }
    createObjectURL.mockClear();
    revokeObjectURL.mockClear();
  });

  it("uploads file and updates value", async () => {
    const handleChange = vi.fn();
    render(<ControlledUploader onValueChange={handleChange} />);

    const input = screen.getByTestId("image-uploader-input") as HTMLInputElement;
    const file = createFile();
    fireEvent.change(input, { target: { files: [file] } });

    const request = MockXMLHttpRequest.instances.at(-1);
    expect(request).toBeTruthy();
    if (!request) return;

    request.status = 201;
    request.response = { data: { url: "https://cdn.example.com/avatar.png", key: "x", contentType: "image/png", size: 10 } };
    await act(async () => {
      request.onload?.();
    });

    await waitFor(() => expect(handleChange).toHaveBeenCalledWith("https://cdn.example.com/avatar.png"));
    expect(screen.getByAltText("Pratinjau gambar")).toBeInTheDocument();
  });

  it("shows error when upload fails", async () => {
    render(<ControlledUploader onValueChange={vi.fn()} />);

    const input = screen.getByTestId("image-uploader-input") as HTMLInputElement;
    fireEvent.change(input, { target: { files: [createFile()] } });

    const request = MockXMLHttpRequest.instances.at(-1);
    expect(request).toBeTruthy();
    if (!request) return;

    request.status = 415;
    request.response = { error: { message: "unsupported" } };
    await act(async () => {
      request.onload?.();
    });

    await waitFor(() => {
      const messages = screen.getAllByText("unsupported");
      expect(messages.length).toBeGreaterThan(0);
    });
  });
});
