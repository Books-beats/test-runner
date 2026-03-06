import { useState, useEffect } from "react";
import { useAuth } from "../context/AuthContext";
import { X } from "lucide-react";
import { Test, TestRunResult } from "../types";
import TestItem from "./TestItem";
import TestResultModal from "./TestResultModal";
import RunTestModal from "./RunTestModal";

export default function TestList({
  refreshTrigger,
}: {
  refreshTrigger: number;
}) {
  const [tests, setTests] = useState<Test[]>([]);
  const [runningTests, setRunningTests] = useState<
    Record<number, { runId: number; status: string; result?: TestRunResult }>
  >({});
  const [isLoading, setIsLoading] = useState(true);
  const [selectedResult, setSelectedResult] = useState<TestRunResult | null>(
    null,
  );
  const [concurrencyDialogOpen, setConcurrencyDialogOpen] = useState(false);
  const [testToRun, setTestToRun] = useState<number | null>(null);
  const [pollingError, setPollingError] = useState<string | null>(null);
  const { token } = useAuth();

  const fetchTests = async () => {
    try {
      const res = await fetch("/tests", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (res.ok) {
        const data = await res.json();
        const fetchedTests = data.tests || [];
        setTests(fetchedTests);

        // Pre-populate runningTests state with the latest run ID from the backend to persist history
        const histories: Record<
          number,
          { runId: number; status: string; result?: TestRunResult }
        > = {};

        fetchedTests.forEach((test: Test) => {
          if (test.latest_run_id && test.latest_run_status) {
            histories[test.id] = {
              runId: test.latest_run_id,
              status: test.latest_run_status,
              result: undefined, // Result data will be fetched if the user clicks 'Show Result'
            };
          }
        });

        setRunningTests((prev) => ({ ...prev, ...histories }));
      }
    } catch (err) {
      console.error("Failed to fetch tests", err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchTests();
  }, [refreshTrigger]);

  useEffect(() => {
    const intervals: Record<number, number> = {};
    const errorCounters: Record<number, number> = {};
    const MAX_RETRIES = 5;

    Object.entries(runningTests).forEach(([testIdStr, runData]) => {
      const testId = Number(testIdStr);
      if (runData.status === "pending") {
        errorCounters[testId] = 0;

        const interval = window.setInterval(async () => {
          try {
            const res = await fetch(`/tests/${runData.runId}`, {
              headers: { Authorization: `Bearer ${token}` },
            });

            if (res.ok) {
              // Reset error counter on success
              errorCounters[testId] = 0;
              const data = await res.json();
              const result = data.result;

              if (result && result.status === "completed") {
                setRunningTests((prev) => ({
                  ...prev,
                  [testId]: { ...prev[testId], status: "completed", result },
                }));
                clearInterval(interval);
              }
            } else {
              throw new Error(`HTTP error ${res.status}`);
            }
          } catch (err) {
            console.error(err);
            errorCounters[testId]++;

            if (errorCounters[testId] >= MAX_RETRIES) {
              setPollingError(
                `Failed to poll results for Test Run #${runData.runId}. The server may be down.`,
              );
              clearInterval(interval);
              setRunningTests((prev) => ({
                ...prev,
                [testId]: { ...prev[testId], status: "error" },
              }));
            }
          }
        }, 2000);
        intervals[testId] = interval;
      }
    });

    return () => {
      Object.values(intervals).forEach(clearInterval);
      setPollingError(null);
    };
  }, [runningTests, token]);

  const openRunDialog = (testId: number) => {
    setTestToRun(testId);
    setConcurrencyDialogOpen(true);
  };

  const handleRunTest = async (concurrencyLevel: number) => {
    if (testToRun === null) return;
    const testId = testToRun;

    try {
      setConcurrencyDialogOpen(false);
      const res = await fetch(`/tests/${testId}/run`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ concurrency: concurrencyLevel }),
      });

      if (res.ok) {
        const data = await res.json();
        setRunningTests((prev) => ({
          ...prev,
          [testId]: { runId: data.testRunId, status: data.status },
        }));
      } else {
        alert(
          "Failed to run test. Make sure backend is updated to return testRunId.",
        );
      }
    } catch (err) {
      console.error(err);
      alert("Error running test");
    } finally {
      setTestToRun(null);
    }
  };

  const handleShowResult = async (runId: number, testId: number) => {
    try {
      const res = await fetch(`/tests/${runId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (res.ok) {
        const data = await res.json();
        setSelectedResult(data.result);
      } else {
        const localResult = runningTests[testId]?.result;
        if (localResult) setSelectedResult(localResult);
      }
    } catch (err) {
      console.error("Error fetching result", err);
    }
  };

  if (isLoading) {
    return <div className="text-secondary">Loading tests...</div>;
  }

  if (tests.length === 0) {
    return (
      <div className="card">
        <p className="text-secondary">
          No tests created yet. Use the form to create one.
        </p>
      </div>
    );
  }

  return (
    <>
      {pollingError && (
        <div
          style={{
            backgroundColor: "var(--error-bg)",
            color: "var(--error)",
            padding: "1rem",
            borderRadius: "var(--radius-md)",
            marginBottom: "1rem",
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
          }}
        >
          <span>{pollingError}</span>
          <button
            className="icon-btn"
            style={{ color: "var(--error)" }}
            onClick={() => setPollingError(null)}
          >
            <X size={16} />
          </button>
        </div>
      )}

      <div style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
        {tests.map((test) => (
          <TestItem
            key={test.id}
            test={test}
            runData={runningTests[test.id]}
            onRunTest={openRunDialog}
            onShowResult={handleShowResult}
          />
        ))}
      </div>

      {selectedResult && (
        <TestResultModal
          result={selectedResult}
          onClose={() => setSelectedResult(null)}
        />
      )}

      {concurrencyDialogOpen && (
        <RunTestModal
          onClose={() => {
            setConcurrencyDialogOpen(false);
            setTestToRun(null);
          }}
          onRun={handleRunTest}
        />
      )}
    </>
  );
}
