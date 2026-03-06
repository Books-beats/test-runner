import { CheckCircle, XCircle, X } from "lucide-react";
import { TestRunResult } from "../types";

interface TestResultModalProps {
  result: TestRunResult;
  onClose: () => void;
}

export default function TestResultModal({
  result,
  onClose,
}: TestResultModalProps) {
  return (
    <div
      style={{
        position: "fixed",
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: "rgba(0,0,0,0.5)",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        zIndex: 50,
        backdropFilter: "blur(4px)",
      }}
    >
      <div
        className="card"
        style={{ width: "100%", maxWidth: "500px", margin: "1rem" }}
      >
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            marginBottom: "1.5rem",
          }}
        >
          <h2>Run Results</h2>
          <button className="icon-btn" onClick={onClose}>
            <X size={20} />
          </button>
        </div>

        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(2, 1fr)",
            gap: "1.5rem",
          }}
        >
          <div
            style={{
              display: "flex",
              flexDirection: "column",
              gap: "0.25rem",
            }}
          >
            <span
              className="text-secondary"
              style={{
                fontSize: "0.75rem",
                textTransform: "uppercase",
                letterSpacing: "0.05em",
              }}
            >
              Total Requests
            </span>
            <span style={{ fontWeight: 600, fontSize: "2rem" }}>
              {result.total}
            </span>
          </div>
          <div
            style={{
              display: "flex",
              flexDirection: "column",
              gap: "0.25rem",
            }}
          >
            <span
              className="text-secondary"
              style={{
                fontSize: "0.75rem",
                textTransform: "uppercase",
                letterSpacing: "0.05em",
              }}
            >
              Avg Response Time
            </span>
            <span style={{ fontWeight: 600, fontSize: "2rem" }}>
              {result.avg_duration_ms}ms
            </span>
          </div>

          <div
            style={{
              display: "flex",
              flexDirection: "column",
              gap: "0.25rem",
              backgroundColor: "var(--success-bg)",
              padding: "1rem",
              borderRadius: "var(--radius-md)",
            }}
          >
            <span
              style={{
                fontSize: "0.75rem",
                textTransform: "uppercase",
                letterSpacing: "0.05em",
                color: "var(--success)",
                fontWeight: 600,
              }}
            >
              Passed
            </span>
            <span
              style={{
                fontWeight: 600,
                fontSize: "1.5rem",
                color: "var(--success)",
                display: "flex",
                alignItems: "center",
                gap: "0.5rem",
              }}
            >
              <CheckCircle size={24} /> {result.passed}
            </span>
          </div>
          <div
            style={{
              display: "flex",
              flexDirection: "column",
              gap: "0.25rem",
              backgroundColor: "var(--error-bg)",
              padding: "1rem",
              borderRadius: "var(--radius-md)",
            }}
          >
            <span
              style={{
                fontSize: "0.75rem",
                textTransform: "uppercase",
                letterSpacing: "0.05em",
                color: "var(--error)",
                fontWeight: 600,
              }}
            >
              Failed
            </span>
            <span
              style={{
                fontWeight: 600,
                fontSize: "1.5rem",
                color: "var(--error)",
                display: "flex",
                alignItems: "center",
                gap: "0.5rem",
              }}
            >
              <XCircle size={24} /> {result.failed}
            </span>
          </div>
        </div>

        <div
          style={{
            marginTop: "2rem",
            display: "flex",
            justifyContent: "flex-end",
          }}
        >
          <button className="btn btn-primary" onClick={onClose}>
            Close
          </button>
        </div>
      </div>
    </div>
  );
}
