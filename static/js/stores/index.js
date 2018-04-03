import Mappings from './Mappings';
import ProxiedRequests from './ProxiedRequests';

const mappings = new Mappings();
mappings.fetch();

const proxiedRequests = new ProxiedRequests();

export default {
  mappings,
  proxiedRequests,
};
