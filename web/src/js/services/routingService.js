import { BehaviorSubject } from 'rxjs';

class RoutingService extends BehaviorSubject {
  constructor() {
    super();

    this.currentPath = window.location.pathname;
    this.next(this.currentPath);

    window.addEventListener('popstate', this.handlePathChanged.bind(this));
  }

  goTo(path) {
    history.pushState(null, null, path);
    this.handlePathChanged();
  }

  handlePathChanged() {
    let newPath = window.location.pathname;
    if (this.currentPath != newPath) {
      this.currentPath = newPath;
      this.next(this.currentPath);
    }
  }
}

const routingService = new RoutingService();

export default routingService;
