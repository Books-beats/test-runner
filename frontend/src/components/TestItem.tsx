import { Play, Activity, BarChart2, Pencil, Trash2 } from "lucide-react";
import { Test, TestRunResult } from "../types";

interface TestItemProps {
  test: Test;
  runData?: { runId: number; status: string; result?: TestRunResult };
  onRunTest: (testId: number) => void;
  onShowResult: (runId: number, testId: number) => void;
  onDelete: (testId: number) => void;
  onEdit: (test: Test) => void;
}

export default function TestItem({
  test,
  runData,
  onRunTest,
  onShowResult,
  onDelete,
  onEdit,
}: TestItemProps) {
  const isRunning = !!runData && runData.status !== "completed";

  return (
    <div
      className="card"
      style={{ display: "flex", flexDirection: "column", gap: "1rem" }}
    >
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
        }}
      >
        <div>
          <h3 style={{ marginBottom: "0.25rem" }}>{test.name}</h3>
          <div
            style={{
              display: "flex",
              gap: "0.5rem",
              alignItems: "center",
              fontSize: "0.875rem",
            }}
            className="text-secondary"
          >
            <span style={{ fontWeight: 600, color: "var(--brand-primary)" }}>
              {test.method}
            </span>
            <span>{test.url}</span>
          </div>
        </div>
        <div style={{ display: "flex", gap: "0.5rem" }}>
          <span
            title={
              isRunning
                ? "Cannot edit while test is running. Wait for it to complete."
                : undefined
            }
            style={{ display: "inline-flex" }}
          >
            <button
              className="btn btn-outline"
              onClick={() => onEdit(test)}
              disabled={isRunning}
              style={isRunning ? { pointerEvents: "none", opacity: 0.45 } : undefined}
            >
              <Pencil size={16} />
            </button>
          </span>

          <button className="btn btn-outline" onClick={() => onDelete(test.id)}>
            <Trash2 size={16} />
          </button>
          {!runData || runData.status === "completed" ? (
            <button
              className="btn btn-primary"
              onClick={() => onRunTest(test.id)}
            >
              <Play size={16} /> {runData ? "Rerun" : "Run"}
            </button>
          ) : (
            <span
              className="badge pending"
              style={{ display: "flex", gap: "0.5rem" }}
            >
              <Activity size={14} className="spin-animation" /> Running...
            </span>
          )}
        </div>
      </div>

      {runData && (
        <div
          style={{
            padding: "1rem",
            backgroundColor: "var(--bg-elevated)",
            borderRadius: "var(--radius-md)",
            border: "1px solid var(--border-color)",
          }}
        >
          <div
            style={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
            }}
          >
            <span style={{ fontWeight: 600, fontSize: "0.875rem" }}>
              Run #{runData.runId}
            </span>

            <div
              style={{
                display: "flex",
                alignItems: "center",
                gap: "1rem",
              }}
            >
              <span className={`badge ${runData.status}`}>
                {runData.status}
              </span>

              {runData.status === "completed" && (
                <button
                  className="btn btn-outline"
                  style={{
                    padding: "0.25rem 0.75rem",
                    fontSize: "0.75rem",
                  }}
                  onClick={() => onShowResult(runData.runId, test.id)}
                >
                  <BarChart2 size={14} /> Show Result
                </button>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
