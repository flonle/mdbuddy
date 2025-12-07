const eventSource = new EventSource('/sse-refresh');

eventSource.onmessage = (event) => {
  eventSource.close();
  // console.log("I would reload now!");
  // console.log(Date.now());
  location.reload();
};

eventSource.onerror = (error) => {
  console.error('SSE connection error:', error);
  eventSource.close();
};
