const eventSource = new EventSource('/sse-refresh');

eventSource.onmessage = (event) => {
  eventSource.close();
  location.reload();
};

eventSource.onerror = (error) => {
  console.error('SSE connection error:', error);
  eventSource.close();
};
