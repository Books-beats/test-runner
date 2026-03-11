interface Header {
  key: string;
  value: string;
}
export interface Test {
  id: number;
  name: string;
  url: string;
  method: string;
  expected_response: string;
  latest_run_id?: number;
  latest_run_status?: string;
  headers?: Header[];
  body?: string;
  status_code?: number;
}

export interface TestRunResult {
  id: number;
  status: string;
  total: number;
  passed: number;
  failed: number;
  avg_duration_ms: number;
}
