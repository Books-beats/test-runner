import { useState, useEffect } from "react";
import { Plus, Trash2 } from "lucide-react";
import { useAuth } from "../context/AuthContext";
import { Test } from "../types";

const API_BASE = import.meta.env.VITE_API_URL || "";

interface Header {
  key: string;
  value: string;
}

interface TestFormProps {
  onRefresh: () => void;
  initialTest?: Test | null;
}

export default function TestForm({ onRefresh, initialTest }: TestFormProps) {
  const { token } = useAuth();

  const [name, setName] = useState("");
  const [url, setUrl] = useState("");
  const [method, setMethod] = useState("GET");
  const [headers, setHeaders] = useState<Header[]>([{ key: "", value: "" }]);
  const [body, setBody] = useState("");
  const [expectedResponse, setExpectedResponse] = useState("");
  const [statusCode, setStatusCode] = useState<number | null>(null);

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const hasAtLeastOneExpectation =
    expectedResponse.trim() !== "" || statusCode !== null;

  const isEditing = !!initialTest;

  useEffect(() => {
    if (!initialTest) return;

    setName(initialTest.name);
    setUrl(initialTest.url);
    setMethod(initialTest.method);

    if (initialTest.headers && Object.keys(initialTest.headers).length > 0) {
      const headerArr: Header[] = Object.entries(initialTest.headers).map(
        ([key, value]) => ({
          key,
          value: value === undefined || value === null ? "" : String(value),
        }),
      );
      setHeaders(headerArr);
    } else {
      setHeaders([{ key: "", value: "" }]);
    }

    setBody(initialTest.body || "");
    setExpectedResponse(initialTest.expected_response || "");
    setStatusCode(initialTest.status_code || null);

    document.querySelector(".form-container")?.scrollIntoView({
      behavior: "smooth",
    });
  }, [initialTest]);

  const handleAddHeader = () => {
    setHeaders([...headers, { key: "", value: "" }]);
  };

  const handleUpdateHeader = (
    index: number,
    field: "key" | "value",
    val: string,
  ) => {
    const newHeaders = [...headers];
    newHeaders[index][field] = val;
    setHeaders(newHeaders);
  };

  const handleRemoveHeader = (index: number) => {
    setHeaders(headers.filter((_, i) => i !== index));
  };

  const resetForm = () => {
    setName("");
    setUrl("");
    setMethod("GET");
    setHeaders([{ key: "", value: "" }]);
    setBody("");
    setExpectedResponse("");
    setStatusCode(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!hasAtLeastOneExpectation) {
      setError("At least one of Expected Response or Expected Status Code is required.");
      return;
    }
    setIsSubmitting(true);
    setError(null);

    // Build header JSON
    const headersObj: Record<string, string> = {};
    headers.forEach((h) => {
      if (h.key.trim()) {
        headersObj[h.key.trim()] = h.value.trim();
      }
    });

    const payload = {
      name,
      url,
      method,
      headers: headersObj,
      body: body || null,
      expected_response: expectedResponse,
      status_code: statusCode,
    };

    try {
      const endpoint = isEditing
        ? `${API_BASE}/tests/${initialTest?.id}/edit`
        : `${API_BASE}/tests`;

      const methodType = isEditing ? "PUT" : "POST";

      const res = await fetch(endpoint, {
        method: methodType,
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || "Request failed");
      }

      resetForm();
      onRefresh();
    } catch (err: unknown) {
      if (err instanceof Error) setError(err.message);
      else setError(String(err));
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="card form-container">
      <h2>{isEditing ? "Edit Test" : "Create New Test"}</h2>

      <p
        className="text-secondary"
        style={{ marginTop: "0.5rem", marginBottom: "1.5rem" }}
      >
        Configure the HTTP request settings and expected behavior.
      </p>

      {error && (
        <div
          style={{
            padding: "0.75rem",
            backgroundColor: "var(--error-bg)",
            color: "var(--error)",
            borderRadius: "var(--radius-md)",
            marginBottom: "1rem",
          }}
        >
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label className="form-label">Test Name</label>
          <input
            required
            type="text"
            placeholder="e.g. Health Check API"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>

        <div
          style={{
            display: "grid",
            gridTemplateColumns: "120px 1fr",
            gap: "1rem",
          }}
        >
          <div className="form-group">
            <label className="form-label">Method</label>
            <select value={method} onChange={(e) => setMethod(e.target.value)}>
              <option value="GET">GET</option>
              <option value="POST">POST</option>
              <option value="PUT">PUT</option>
              <option value="DELETE">DELETE</option>
              <option value="PATCH">PATCH</option>
            </select>
          </div>
          <div className="form-group">
            <label className="form-label">URL</label>
            <input
              required
              type="url"
              placeholder="https://api.example.com/v1/health"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
            />
          </div>
        </div>

        <div className="form-group">
          <label className="form-label">Headers</label>
          {headers.map((header, idx) => (
            <div
              key={idx}
              style={{
                display: "grid",
                gridTemplateColumns: "1fr 1fr auto",
                gap: "0.5rem",
                marginBottom: "0.5rem",
              }}
            >
              <input
                type="text"
                placeholder="Header Key (e.g. Authorization)"
                value={header.key}
                onChange={(e) => handleUpdateHeader(idx, "key", e.target.value)}
              />
              <input
                type="text"
                placeholder="Value"
                value={header.value}
                onChange={(e) =>
                  handleUpdateHeader(idx, "value", e.target.value)
                }
              />
              <button
                type="button"
                className="icon-btn"
                onClick={() => handleRemoveHeader(idx)}
                disabled={headers.length === 1 && !header.key}
              >
                <Trash2 size={18} />
              </button>
            </div>
          ))}
          <button
            type="button"
            className="btn btn-outline"
            onClick={handleAddHeader}
            style={{ alignSelf: "flex-start", marginTop: "0.5rem" }}
          >
            <Plus size={16} /> Add Header
          </button>
        </div>

        <div className="form-group">
          <label className="form-label">Request Body</label>
          <textarea
            placeholder='{"foo": "bar"}'
            value={body}
            onChange={(e) => setBody(e.target.value)}
          />
        </div>

        <div className="form-group">
          <label className="form-label">Expected Response String</label>
          <textarea
            placeholder="Exact response body string to match against..."
            value={expectedResponse}
            onChange={(e) => setExpectedResponse(e.target.value)}
          />
        </div>

        <div className="form-group">
          <label className="form-label">Expected Status Code</label>
          <input
            type="number"
            value={statusCode ?? ""}
            placeholder="200"
            onChange={(e) =>
              setStatusCode(e.target.value ? Number(e.target.value) : null)
            }
          />
        </div>

        <button
          type="submit"
          className="btn btn-primary"
          style={{ width: "100%", marginTop: "1rem" }}
          disabled={isSubmitting}
        >
          {isSubmitting
            ? isEditing
              ? "Updating Test..."
              : "Creating Test..."
            : isEditing
              ? "Update Test"
              : "Create Test"}
        </button>
      </form>
    </div>
  );
}
