import Mappings from './Mappings';
import ProxiedRequests from './ProxiedRequests';
import Variables from './Variables';

const mappings = new Mappings();
mappings.fetch();

const proxiedRequests = new ProxiedRequests();

const variables = new Variables();
variables.fetch();

export default {
  mappings,
  proxiedRequests,
  variables,
};
