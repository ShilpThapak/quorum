import { Component, type ReactNode } from "react";

type State = { message: string | null };

export class ErrorBoundary extends Component<{ children: ReactNode }, State> {
  state: State = { message: null };

  static getDerivedStateFromError(e: unknown): State {
    return { message: e instanceof Error ? e.message : String(e) };
  }

  componentDidCatch(e: unknown) {
    console.error("Quorum render error:", e);
  }

  render() {
    if (this.state.message !== null) {
      return (
        <div className="banner err" role="alert">
          <b>Something broke in the results panel:</b> {this.state.message}
          <button className="reset" onClick={() => this.setState({ message: null })}>
            Dismiss
          </button>
        </div>
      );
    }
    return this.props.children;
  }
}
