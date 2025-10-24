"use client";

import { Component, ErrorInfo, ReactNode } from "react";

type GlobalErrorBoundaryProps = {
  children: ReactNode;
};

type GlobalErrorBoundaryState = {
  hasError: boolean;
  error?: Error;
};

export class GlobalErrorBoundary extends Component<
  GlobalErrorBoundaryProps,
  GlobalErrorBoundaryState
> {
  public state: GlobalErrorBoundaryState = {
    hasError: false,
  };

  static getDerivedStateFromError(error: Error): GlobalErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, info: ErrorInfo): void {
    if (process.env.NODE_ENV !== "production") {
      console.error("Global error captured", error, info);
    }
  }

  render(): ReactNode {
    if (this.state.hasError && this.state.error) {
      return (
        <div className="p-4 text-red-600">
          Unexpected error: {this.state.error.message}
        </div>
      );
    }

    return this.props.children;
  }
}
