export const ws = new WebSocket(`ws://${location.hostname}:8080/ws`);

ws.onopen = () => {
  console.log("ws opened")
}

ws.onclose = () => {
  console.log("ws closed")
}

ws.onmessage = (event) => {
  const { type, payload } = JSON.parse(event.data);
  if (eventCallbacks.has(type)) {
    eventCallbacks.get(type)?.forEach((callback) => {
      callback(payload);
    });
  }
};

const eventCallbacks = new Map();

export const subscribe = (type, callback) => {
  if (!eventCallbacks.has(type)) {
    eventCallbacks.set(type, []);
  }
  eventCallbacks.get(type)?.push(callback);
};

export const unsubscribe = (type, callback) => {
  const callbacks = eventCallbacks.get(type);
  if (callbacks) {
    eventCallbacks.set(
      type,
      callbacks.filter((cb) => cb !== callback)
    );
  }
};

export const triggerEvent = (type, eventData) => {
  const callbacks = eventCallbacks.get(type);
  if (callbacks) {
    callbacks.forEach((callback) => {
      callback(eventData);
    });
  }
};