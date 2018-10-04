import Configurations from './Configurations';
import Mappings from './Mappings';
import PossibleValues from './PossibleValues';
import ProxiedRequests from './ProxiedRequests';
import SelectedValues from './SelectedValues';
import Variables from './Variables';

const configurations = new Configurations();
configurations.fetch();

const mappings = new Mappings();
mappings.fetch();

const possibleValues = new PossibleValues();
possibleValues.fetch();

const proxiedRequests = new ProxiedRequests();

const selectedValues = new SelectedValues();
selectedValues.fetch();

const variables = new Variables();
variables.fetch();

export default {
  configurations,
  mappings,
  possibleValues,
  proxiedRequests,
  selectedValues,
  variables,
};
