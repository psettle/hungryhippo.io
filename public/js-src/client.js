console.log('WebSocket client script will run here.');

const ws = new WebSocket('ws://localhost/ws');

ws.addEventListener('open', () => {
    // Send a message to the WebSocket server
    ws.send(JSON.stringify({ type: "Hello!"}));
  });

ws.addEventListener('message', event => {
  // The `event` object is a typical DOM event object, and the message data sent
  // by the server is stored in the `data` property
  console.log('Received:', event.data);
});