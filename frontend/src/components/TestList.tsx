import { useState, useEffect } from "react";
import { useAuth } from "../context/AuthContext";
import { Test, TestRunResult } from "../types";
import TestItem from "./TestItem";
import TestResultModal from "./TestResultModal";
import RunTestModal from "./RunTestModal";

const API_BASE = import.meta.env.VITE_API_URL || "";

export default function TestList({
  refreshTrigger,
  setEditingTest,
}: {
  refreshTrigger: number;
  setEditingTest: (test: Test | null) => void;
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
  const { token } = useAuth();

  const fetchTests = async () => {
    try {
      const res = await fetch(`${API_BASE}/tests`, {
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

  const openRunDialog = (testId: number) => {
    setTestToRun(testId);
    setConcurrencyDialogOpen(true);
  };

  const handleRunTest = async (concurrencyLevel: number) => {
    if (testToRun === null) return;
    const testId = testToRun;

    try {
      setConcurrencyDialogOpen(false);
      const res = await fetch(`${API_BASE}/tests/${testId}/run`, {
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
      const res = await fetch(`${API_BASE}/tests/${runId}`, {
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

  const handleDeleteTest = async (testId: number) => {
    try {
      const res = await fetch(`${API_BASE}/tests/${testId}/delete`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });
      if (res.ok) {
        setTests((prev) => prev.filter((t) => t.id !== testId));
      } else {
        console.error("Delete failed");
      }
    } catch (e) {
      console.error("Error in deleting", e);
    }
  };

  const handleEditTest = (test: Test) => {
    setEditingTest(test);
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
      <div style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
        {tests.map((test) => (
          <TestItem
            key={test.id}
            test={test}
            runData={runningTests[test.id]}
            onRunTest={openRunDialog}
            onShowResult={handleShowResult}
            onDelete={handleDeleteTest}
            onEdit={handleEditTest}
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
