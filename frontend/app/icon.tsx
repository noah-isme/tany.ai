import { ImageResponse } from "next/og";

export const size = {
  width: 64,
  height: 64,
};

export const contentType = "image/png";

export default function Icon() {
  return new ImageResponse(
    (
      <div
        style={{
          height: "100%",
          width: "100%",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          background: "radial-gradient(circle at 30% 30%, #38bdf8, #0f172a)",
          color: "#f8fafc",
          fontSize: 32,
          fontWeight: 700,
          letterSpacing: "-0.05em",
          fontFamily: "'Geist', 'Inter', 'Arial', sans-serif",
          textTransform: "lowercase",
        }}
      >
        ta
      </div>
    ),
    {
      width: size.width,
      height: size.height,
    }
  );
}
