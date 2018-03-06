import Configurations from './Configurations';
import ProxiedRequests from './ProxiedRequests';

const configurations = new Configurations();
configurations.fetch();

const proxiedRequests = new ProxiedRequests();

export default {
  configurations,
  proxiedRequests,
};
