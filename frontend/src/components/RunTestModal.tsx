import { useState } from "react";
import { Play, X } from "lucide-react";

interface RunTestModalProps {
  onClose: () => void;
  onRun: (concurrency: number) => void;
}

export default function RunTestModal({ onClose, onRun }: RunTestModalProps) {
  const [concurrency, setConcurrency] = useState<number>(1);

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
        style={{ width: "100%", maxWidth: "400px", margin: "1rem" }}
      >
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            marginBottom: "1.5rem",
          }}
        >
          <h2>Run Test</h2>
          <button className="icon-btn" onClick={onClose}>
            <X size={20} />
          </button>
        </div>

        <div style={{ marginBottom: "1.5rem" }}>
          <label
            htmlFor="concurrency-input"
            style={{
              display: "block",
              marginBottom: "0.5rem",
              fontSize: "0.875rem",
              fontWeight: 500,
            }}
          >
            Concurrency
          </label>
          <input
            id="concurrency-input"
            type="number"
            min="1"
            className="input"
            value={concurrency}
            onChange={(e) => setConcurrency(parseInt(e.target.value) || 1)}
            style={{
              width: "100%",
              padding: "0.75rem",
              borderRadius: "var(--radius-md)",
              border: "1px solid var(--border-color)",
              backgroundColor: "var(--bg-elevated)",
              color: "var(--text-primary)",
            }}
          />
        </div>

        <div
          style={{
            display: "flex",
            justifyContent: "flex-end",
            gap: "1rem",
          }}
        >
          <button className="btn btn-outline" onClick={onClose}>
            Cancel
          </button>
          <button
            className="btn btn-primary"
            onClick={() => onRun(concurrency)}
          >
            <Play size={16} /> Start
          </button>
        </div>
      </div>
    </div>
  );
}
