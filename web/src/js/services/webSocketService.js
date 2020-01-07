import { BehaviorSubject } from 'rxjs';

class WebSocketService extends BehaviorSubject {
  constructor() {
    super();

    this.socket = new WebSocket(`ws://${window.location.host}/api/websocket`);
    this.socket.onmessage = (message) => {
      this.next(JSON.parse(message.data));
    };
  }
}

const webSocketService = new WebSocketService();
export default webSocketService;
